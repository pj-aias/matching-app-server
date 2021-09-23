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
