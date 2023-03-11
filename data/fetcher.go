package data

import (
	"github.com/auribuo/novasearch/data/ned"
	"github.com/auribuo/novasearch/data/ugc"
	"github.com/auribuo/novasearch/sql"
	"github.com/auribuo/novasearch/types"
	"github.com/auribuo/novasearch/types/coordinates"
)

type Fetcher interface {
	Fetch() ([]types.Galaxy, error)
}

var DefaultFetcher Fetcher = defaultDataFetcher{}

type defaultDataFetcher struct{}

func (f defaultDataFetcher) Fetch() ([]types.Galaxy, error) {
	ugcData, err := ugc.Fetch()
	if err != nil {
		return nil, err
	}
	nedData, err := ned.Fetch(ugcData)
	if err != nil {
		return nil, err
	}

	data := make([]types.Galaxy, 0)

	for _, ugcResponse := range ugcData {
		nedResponse := sql.FirstF(nedData, func(t ned.Response) bool {
			return t.UgcNumber == ugcResponse.UgcNumber
		})
		if nedResponse == nil {
			continue
		}

		magnitude := ugcResponse.Magnitude
		if nedResponse.Magnitude > 0 {
			magnitude = nedResponse.Magnitude
		}

		data = append(data, types.Galaxy{
			Morphology:    nedResponse.HubbleType,
			Id:            ugcResponse.UgcNumber,
			PreferredName: nedResponse.PreferredName,
			Magnitude:     magnitude,
			EquatorialCoordinates: coordinates.EquatorialCoordinates{
				RightAscension: ugcResponse.RightAscension,
				Declination:    ugcResponse.Declination,
			},
			AzimuthalCoordinates: coordinates.AzimuthalCoordinates{},
			SemiMajorAxis:        ugcResponse.SemiMajorAxis,
			SemiMinorAxis:        ugcResponse.SemiMinorAxis,
			Redshift:             nedResponse.Redshift,
		})
	}

	return data, nil
}
