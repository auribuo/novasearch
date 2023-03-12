package types

type Repo interface {
	Init() error
	Galaxies() GalaxyCollection
	Galaxy(id int) Galaxy
	GalaxiesFiltered(data BaseData) GalaxyCollection
	ViewportsFiltered(collection GalaxyCollection, data BaseData) []Viewport
	AddAzimuthalCoordinates(galaxies []Galaxy, data BaseData) []Galaxy
}
