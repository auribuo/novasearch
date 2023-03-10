package types

type ExtendedGalaxy struct {
	Galaxy
	Distance float64 `json:"distance"`
}
