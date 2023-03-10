package data

import (
	"github.com/auribuo/novasearch/sql"
	"github.com/auribuo/novasearch/types"
	"strings"
)

type Repo struct{}

var isInitialized bool

var galaxies map[int]types.Galaxy

func Init(fetcher Fetcher) error {
	isInitialized = true
	data, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	cleanAndMap(data)
	return nil
}

func Galaxies() []types.Galaxy {
	if !isInitialized {
		panic("Repo is not initialized")
	}
	return sql.Values(galaxies)
}

func Galaxy(id int) types.Galaxy {
	if !isInitialized {
		panic("Repo is not initialized")
	}
	return galaxies[id]
}

func cleanAndMap(galaxies []types.Galaxy) {
	for _, galaxy := range galaxies {
		if galaxy.Magnitude < 0 {
			continue
		}
		if galaxy.SemiMajorAxis < 0 {
			continue
		}
		if galaxy.SemiMinorAxis < 0 {
			continue
		}
		if galaxy.Redshift < 0 {
			continue
		}
		galaxy.Morphology = normalizeMorphology(galaxy.Morphology)
		galaxies[galaxy.Id] = galaxy
	}
}

func normalizeMorphology(morphology string) string {
	if len(morphology) == 0 {
		return "Unknown"
	}

	morphology = strings.Replace(morphology, "?", "", -1)
	morphology = strings.Replace(morphology, ":", " ", -1)

	if strings.HasPrefix(morphology, "E") {
		return "E"
	}

	relevantTypes := []string{
		"S0",
		"Sa",
		"Sb",
		"Sbc",
		"Sc",
		"Scd",
	}

	for _, relevantType := range relevantTypes {
		if strings.HasPrefix(morphology, relevantType) {
			return relevantType
		}
	}

	return "Unknown"
}
