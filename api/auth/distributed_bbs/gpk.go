package distributed_bbs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/pj-aias/matching-app-server/util"
)

const (
	MIME_PLAINTEXT = "text/plain"
)

type GpkRegistry interface {
	GetGpk(gm Gm, gms Gms) (Gpk, error)
}

type Gms []Gm

func (gms Gms) Len() int           { return len(gms) }
func (gms Gms) Less(i, j int) bool { return gms[i] < gms[j] }
func (gms Gms) Swap(i, j int)      { gms[i], gms[j] = gms[j], gms[i] }

func CombineGpk(gpks []Gpk) (Gpk, error) {
	// TODO implement CpmbineGpk process
	return gpks[0], nil
}

func GetGpkFromGms(registry GpkRegistry, gms Gms) (Gpk, error) {
	gpks := []Gpk{}
	sort.Sort(gms)

	for _, gm := range gms {
		gpk, err := registry.GetGpk(gm, gms)
		if err != nil {
			return "", fmt.Errorf("could not get gpk from %v: %v", gm, err)
		}

		gpks = append(gpks, gpk)
	}

	combinedGpk, err := CombineGpk(gpks)
	if err != nil {
		return "", fmt.Errorf("could not combine gpks: %v", err)
	}

	return combinedGpk, nil
}

type CachedRegistry struct {
	Cache map[Gm]Gpk
}

func (r CachedRegistry) GetGpk(gm Gm, gms Gms) (Gpk, error) {
	cacheKey := gms.toKey()
	// if cache found, return it
	if cache := r.Cache[cacheKey]; cache != "" {
		return cache, nil
	}

	gpks := [3]gpkResp{}

	for i, gm := range gms {
		resp, err := sendRequest(gm, gms)
		if err != nil {
			return "", err
		}

		gpks[i] = resp
	}

	return combine(gpks)
}

func (gms Gms) toKey() string {
	return fmt.Sprintf("%v,%v,%v", gms[0], gms[1], gms[2])
}

type fp [6]uint64

type gpkResp struct {
	H       fp         `json:"h"`
	U       fp         `json:"u"`
	V       fp         `json:"v"`
	Partial partialGpk `json:"partial"`
}

func sendRequest(gm Gm, gms Gms) (gpkResp, error) {
	tor, err := util.NewTorClient()
	if err != nil {
		return gpkResp{}, err
	}

	body := bufio.NewReader(nil)

	resp, err := tor.Post("http://"+gm+"/pubkey", MIME_PLAINTEXT, body)
	if err != nil {
		return gpkResp{}, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return gpkResp{}, err
	}

	gpk := gpkResp{}

	err = json.Unmarshal(respBody, &gpk)
	return gpk, err
}

type partialGpk map[string]interface{}

type combinedGpk struct {
	H           fp           `json:"h"`
	U           fp           `json:"u"`
	V           fp           `json:"v"`
	W           fp           `json:"w"`
	PartialGpks []partialGpk `json:"partial_gpks"`
}

func combine(gpks [3]gpkResp) (Gpk, error) {
	combined := combinedGpk{
		H: gpks[0].H,
		U: gpks[0].U,
		V: gpks[0].V,
		W: gpks[1].U,
		PartialGpks: []partialGpk{
			gpks[0].Partial,
			gpks[1].Partial,
			gpks[2].Partial,
		},
	}

	gpk, err := combined.toGpk()
	if err != nil {
		return "", fmt.Errorf("failed to combine gpks: %v", err)
	}

	return gpk, nil
}

func (combined combinedGpk) toGpk() (Gpk, error) {
	bytes, err := json.Marshal(combined)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
