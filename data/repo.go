package data

import (
	"strings"

	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/types"
)

var repo types.Repo

func (repo *repoImpl) cleanAndMap(data []types.Galaxy) {
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
		repo.galaxies[galaxy.Id] = galaxy
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

func R() types.Repo {
	return repo
}

func NewDefault() types.Repo {
	if repo == nil {
		repo = New(DefaultFetcher)
	}
	return repo
}

func New(fetcher Fetcher) types.Repo {
	return &repoImpl{
		fetcher:       fetcher,
		galaxies:      map[int]types.Galaxy{},
		isInitialized: false,
	}
}
