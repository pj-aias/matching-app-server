package distributed_bbs

import (
	"fmt"
	"sort"
)

type GpkRegistry interface {
	GetGpk(gm Gm) (Gpk, error)
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
		gpk, err := registry.GetGpk(gm)
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

func (r CachedRegistry) Get(gm Gm) (Gpk, error) {
	// if cache found, return it
	if cache := r.Cache[gm]; cache != "" {
		return cache, nil
	}

	// TODO: get from the internet
	gpk := "someGpk"

	return gpk, nil
}

type fp [6]uint64

type GpkResp struct {
	H fp `json:"h"`
	U fp `json:"u"`
	V fp `json:"v"`
}

func sendRequest(gm Gm)
