package filter

import (
	"math"
	"reflect"
	"sort"
	"time"

	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/types"
)

type path struct {
	objects []types.Ratable
	time    time.Duration
}

func (p path) deconstruct() ([]types.Ratable, time.Duration) {
	return p.objects, p.time
}

func Generic(objects []types.Ratable, data types.CalculationData, repo types.Repo) ([]types.Ratable, time.Duration) {
	startPoints := determineStartPoints(objects, data)
	paths := make([]path, len(startPoints))
	waiters := make([]chan path, len(startPoints))

	for i := range startPoints {
		waiters[i] = make(chan path)
		go func(i int) {
			calculatePath(startPoints[i], objects, data, &waiters, i, repo)
		}(i)
	}

	for i := range waiters {
		paths[i] = <-waiters[i]
	}

	return bestPath(paths).deconstruct()
}

func bestPath(paths []path) path {
	bestPath := paths[0]
	bestScore := pathScore(bestPath.objects)
	for _, path := range paths {
		score := pathScore(path.objects)
		if score > bestScore {
			bestScore = score
			bestPath = path
		}
	}
	return bestPath
}

func pathScore(path []types.Ratable) float64 {
	score := 0.0
	for _, object := range path {
		score += object.Quality()
	}
	return score
}

func calculatePath(start types.Ratable, objects []types.Ratable, data types.CalculationData, waiters *[]chan path, pathIndex int, repo types.Repo) {
	calculatedPath := make([]types.Ratable, 0)
	usedTime := time.Duration(start.WaitTime(0)) * time.Second
	start.Visit()
	pickNext(start, objects, &data, &calculatedPath, &usedTime, time.Duration(data.MaxSearchTime)*time.Second, repo)
	(*waiters)[pathIndex] <- path{
		objects: calculatedPath,
		time:    usedTime,
	}
}

func pickNext(start types.Ratable, objects []types.Ratable, data *types.CalculationData, path *[]types.Ratable, usedTime *time.Duration, maxSearchTime time.Duration, repo types.Repo) {
	if *usedTime > maxSearchTime || (len(objects) == 0) {
		return
	}

	var next types.Ratable

	if data.SearchRadius < 0 {
		log.Default.Warnf("Search radius is negative, using infinite search radius")
		nextIndex := pickBestIndex(objects, start)
		next = objects[nextIndex]
		next.Visit()
		*path = append(*path, next)
	} else {
		nextIndex := pickBestInRangeIndex(start, objects, (*data).SearchRadius)
		next = objects[nextIndex]
		next.Visit()
		*path = append(*path, next)
	}

	distance := next.DistanceTo(start)
	waitTime := time.Duration(next.WaitTime(distance)) * time.Second
	*usedTime += waitTime
	nextTime := data.StartTime.Convert().Add(*usedTime)
	nextSkyState := repo.GalaxiesFiltered(data.BaseData)
	data.StartTime = types.IsoTime(nextTime.Format(time.RFC3339))
	next.Mark(nextTime)
	*path = append(*path, next)
	pickNext(next, nextSkyState.Ratable(), data, path, usedTime, maxSearchTime, repo)
}

func pickBestIndex(objects []types.Ratable, start types.Ratable) int {
	best := objects[0]
	bestIndex := 0
	for i, object := range objects {
		if object.Quality() > best.Quality() && !reflect.DeepEqual(object, start) {
			best = object
			bestIndex = i
		}
	}
	return bestIndex
}

func pickBestInRangeIndex(start types.Ratable, objects []types.Ratable, radius float64) int {
	best := objects[0]
	bestIndex := 0
	foundAny := false
	for i, object := range objects {
		if object.DistanceTo(start) <= radius && object.Quality() > best.Quality() && !reflect.DeepEqual(object, start) {
			best = object
			bestIndex = i
			foundAny = true
		}
	}
	if !foundAny {
		log.Default.Warnf("No objects in range, using next best object")
		return pickBestIndex(objects, start)
	}
	return bestIndex
}

func determineStartPoints(objects []types.Ratable, data types.CalculationData) []types.Ratable {
	maxPoints := int(math.Min(float64(data.StartPointCount), float64(len(objects))))

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Quality() > objects[j].Quality()
	})

	startPoints := make([]types.Ratable, maxPoints)

	for i := 0; i < maxPoints; i++ {
		startPoints[i] = objects[i]
	}

	return startPoints
}
