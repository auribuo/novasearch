package types

import "sort"

type GalaxyCollection []Galaxy

func (c GalaxyCollection) Len() int {
	return len(c)
}

func (c GalaxyCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c GalaxyCollection) Less(i, j int) bool {
	return c[i].Id < c[j].Id
}

func (c GalaxyCollection) Sort() {
	sort.Sort(c)
}

func (c GalaxyCollection) Ratable() []Ratable {
	ratable := make([]Ratable, len(c))
	for i := range c {
		ratable[i] = &c[i]
	}
	return ratable
}
