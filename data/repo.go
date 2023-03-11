package data

import (
	"sort"
	"strings"

	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/sql"

	"github.com/auribuo/novasearch/types"
)

type Repo struct{}

var isInitialized bool

var galaxies map[int]types.Galaxy

func Init(fetcher Fetcher) error {
	data, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	galaxies = make(map[int]types.Galaxy, len(data))
	cleanAndMap(data)
	isInitialized = true
	log.Default.Debugf("initialized repo with %d galaxies", len(galaxies))
	return nil
}

func verifyInitialized() {
	if !isInitialized {
		log.Default.Error("trying to get galaxies before repo is initialized")
		panic("accessing repo before it is initialized")
	}
}

func Galaxies() []types.Galaxy {
	verifyInitialized()
	values := sql.Values(galaxies)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Id < values[j].Id
	})
	return values
}

func GalaxiesFiltered(data types.BaseData) []types.Galaxy {
	verifyInitialized()
	values := sql.Values(galaxies)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Id < values[j].Id
	})
	// TODO: filter galaxies
	return values
}

func Galaxy(id int) types.Galaxy {
	verifyInitialized()
	return galaxies[id]
}

func cleanAndMap(data []types.Galaxy) {
	for _, galaxy := range data {
		if galaxy.Magnitude < 0 {
			log.Default.Debugf("skipping galaxy UGC%d because it has a negative magnitude", galaxy.Id)
			continue
		}
		if galaxy.SemiMajorAxis < 0 {
			log.Default.Debugf("skipping galaxy UGC%d because it has a negative semi-major axis", galaxy.Id)
			continue
		}
		if galaxy.SemiMinorAxis < 0 {
			log.Default.Debugf("skipping galaxy UGC%d because it has a negative semi-minor axis", galaxy.Id)
			continue
		}
		if galaxy.Redshift < 0 {
			log.Default.Debugf("skipping galaxy UGC%d because it has a negative redshift", galaxy.Id)
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
