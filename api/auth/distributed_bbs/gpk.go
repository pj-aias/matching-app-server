package distributed_bbs

import "fmt"

type GpkRegistry interface {
	GetGpk(gm Gm) (Gpk, error)
}

func CombineGpk(gpks []Gpk) (Gpk, error) {
	// TODO implement CpmbineGpk process
	return gpks[0], nil
}

func GetGpkFromGms(registry GpkRegistry, gms []Gm) (Gpk, error) {
	gpks := []Gpk{}
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
