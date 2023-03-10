package situational

import (
	"fmt"
	"github.com/auribuo/novasearch/types"
	"github.com/auribuo/novasearch/types/coordinates"
)

func Filter(galaxies []types.Galaxy, hemisphere types.Hemisphere, minHeight float64, maxMajorAxis float64, maxMinorAxis float64) ([]types.Galaxy, error) {
	filtered := make([]types.Galaxy, 0)
	for _, galaxy := range galaxies {
		if (galaxy.AzimuthalCoordinates == coordinates.AzimuthalCoordinates{}) {
			return nil, fmt.Errorf("galaxy UGC%d has no azimuthal coordinates", galaxy.UgcNumber)
		}
		if galaxy.Position().Elevation < minHeight {
			continue
		}

		if (galaxy.SemiMajorAxis > maxMajorAxis) || (galaxy.SemiMinorAxis > maxMinorAxis) {
			continue
		}

		if hemisphere == types.East {
			if galaxy.Position().Azimuth > 0 && galaxy.Position().Azimuth < 180 {
				filtered = append(filtered, galaxy)
			}
		} else if hemisphere == types.West {
			if galaxy.Position().Azimuth > 180 && galaxy.Position().Azimuth < 360 {
				filtered = append(filtered, galaxy)
			}
		} else {
			return nil, fmt.Errorf("invalid hemisphere %s", hemisphere)
		}
	}
	return filtered, nil
}
