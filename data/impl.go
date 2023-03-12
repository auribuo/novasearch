package data

import (
	"sort"

	"github.com/auribuo/novasearch/filter"
	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/sql"
	"github.com/auribuo/novasearch/types"
)

type repoImpl struct {
	fetcher       Fetcher
	isInitialized bool
	galaxies      map[int]types.Galaxy
}

func (repo *repoImpl) Init() error {
	data, err := repo.fetcher.Fetch()
	if err != nil {
		return err
	}

	repo.galaxies = make(map[int]types.Galaxy, len(data))
	repo.cleanAndMap(data)
	repo.isInitialized = true
	log.Default.Debugf("initialized repo with %d galaxies", len(repo.galaxies))
	return nil
}

func (repo *repoImpl) verifyInitialized() {
	if !repo.isInitialized {
		log.Default.Errorf("trying to get galaxies before repo is initialized")
		panic("accessing repo before it is initialized")
	}
}

func (repo *repoImpl) Galaxies() types.GalaxyCollection {
	repo.verifyInitialized()
	values := sql.Values(repo.galaxies)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Id < values[j].Id
	})
	return values
}

func (repo *repoImpl) GalaxiesFiltered(data types.BaseData) types.GalaxyCollection {
	repo.verifyInitialized()
	values := sql.Values(repo.galaxies)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Id < values[j].Id
	})

	for i := range values {
		values[i].AzimuthalCoordinates = values[i].EquatorialCoordinates.ToAzimuthalCoordinates(data.StartTime.Convert(), data.Location)
	}

	values, err := filter.Situational(values, data.Hemisphere, data.MinHeight, data.MaxSemiMajorAxis, data.MaxSemiMinorAxis)

	if err != nil {
		log.Default.Error(err)
		return []types.Galaxy{}
	}

	return values
}

func (repo *repoImpl) ViewportsFiltered(galaxies types.GalaxyCollection, data types.BaseData) []types.Viewport {
	repo.verifyInitialized()

	galaxies, err := filter.Situational(galaxies, data.Hemisphere, data.MinHeight, data.MaxSemiMajorAxis, data.MaxSemiMinorAxis)

	if err != nil {
		log.Default.Error(err)
		return []types.Viewport{}
	}

	viewports := types.CalculateViewports(galaxies, data.Fov, data.Location, data.StartTime.Convert(), 0.25)

	return viewports
}

func (repo *repoImpl) Galaxy(id int) types.Galaxy {
	repo.verifyInitialized()
	return repo.galaxies[id]
}

func (repo *repoImpl) AddAzimuthalCoordinates(galaxies []types.Galaxy, data types.BaseData) []types.Galaxy {
	for i := range galaxies {
		galaxies[i].AzimuthalCoordinates = galaxies[i].EquatorialCoordinates.ToAzimuthalCoordinates(data.StartTime.Convert(), data.Location)
	}
	return galaxies
}
