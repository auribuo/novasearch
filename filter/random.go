package filter

import (
	random "math/rand"
	"time"

	"github.com/auribuo/novasearch/types"
)

func Random(objects []types.Ratable, maxSearchTime time.Duration) ([]types.Ratable, time.Duration) {
	randSource := random.NewSource(time.Now().UnixNano())
	rand := random.New(randSource)
	path := make([]types.Ratable, 0)

	startIndex := rand.Intn(len(objects))
	start := objects[startIndex]
	path = append(path, start)
	start.Visit()

	usedTime := time.Duration(0)

	pickRandom(rand, start, objects, &path, maxSearchTime, &usedTime)
	return path, usedTime
}

func pickRandom(rand *random.Rand, start types.Ratable, objects []types.Ratable, path *[]types.Ratable, maxSearchTime time.Duration, usedTime *time.Duration) {
	if (*usedTime > maxSearchTime) || (len(objects) == 0) {
		return
	}

	nextIndex := rand.Intn(len(objects))
	next := objects[nextIndex]

	waitTime := next.WaitTime(next.DistanceTo(start))
	*usedTime += time.Duration(waitTime) * time.Second
	if *usedTime > maxSearchTime {
		return
	}

	*path = append(*path, next)
	next.Visit()
	pickRandom(rand, next, objects, path, maxSearchTime, usedTime)
}
