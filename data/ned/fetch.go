package ned

import (
	"encoding/json"
	"github.com/auribuo/novasearch/data/ugc"
	"github.com/auribuo/novasearch/fs"
	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/sql"
	"github.com/go-resty/resty/v2"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Fetch(responses []ugc.Response) ([]Response, error) {
	cacheFolder, err := fs.CacheDir()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cacheFolder); os.IsNotExist(err) {
		err := os.MkdirAll(cacheFolder, 0755)
		if err != nil {
			return nil, err
		}
	}

	cacheFile := filepath.FromSlash(cacheFolder + "/ned.json")

	log.Logger.Info("trying to fetch NED data from cache")
	local, err := fetchLocal(cacheFile)
	if err != nil {
		return nil, err
	}
	if local != nil && len(local) > 0 {
		log.Logger.Info("successfully fetched NED data from cache")
		return local, nil
	}
	log.Logger.Info("no hit. trying to fetch NED data from remote")
	data, err := fetchRemote(responses)
	if err != nil {
		return nil, err
	}
	log.Logger.Info("successfully fetched NED data from remote. saving to cache")
	cache := Cache{
		LastUpdated: time.Now(),
		Items:       data,
	}
	cacheData, err := json.Marshal(cache)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(cacheFile, cacheData, 0644)
	if os.IsNotExist(err) {
		f, err := os.Create(cacheFile)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(cacheData)
	}
	return data, nil
}

func fetchLocal(file string) ([]Response, error) {
	cacheContent, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var cache Cache
	err = json.Unmarshal(cacheContent, &cache)
	if err != nil {
		return nil, err
	}
	if time.Now().Sub(cache.LastUpdated) < 365*time.Hour*24 {
		return cache.Items, nil
	}
	return nil, nil
}

func fetchRemote(responses []ugc.Response) ([]Response, error) {
	client := resty.New()

	numbers := sql.Map(responses, func(t ugc.Response) int {
		return t.UgcNumber
	})

	numberChunks := sql.Chunk(numbers, 500)

	waiters := make([]chan []Response, len(numberChunks))
	var data []Response

	for i, chunk := range numberChunks {
		waiters[i] = make(chan []Response)
		go func(chunk []int, waiter chan []Response) {
			data, err := fetchChunk(client, chunk)
			if err != nil {
				waiter <- nil
				return
			}
			waiter <- data
		}(chunk, waiters[i])
	}

	for _, waiter := range waiters {
		data = append(data, <-waiter...)
	}

	return data, nil
}

func fetchChunk(client *resty.Client, numbers []int) ([]Response, error) {
	ugcNumbers := sql.Map(numbers, func(t int) string {
		return "UGC" + strconv.Itoa(t)
	})

	queryValues := url.Values{}
	queryValues.Add("uplist", strings.Join(ugcNumbers, "\r"))
	queryValues.Add("delimiter", "bar")
	queryValues.Add("NO_LINKS", "1")
	queryValues.Add("nondb", "user_objname")
	queryValues.Add("crosid", "objname")
	queryValues.Add("position", "z")
	queryValues.Add("gadata", "magnit")
	queryValues.Add("attdat_CON", "M")
	queryValues.Add("attdat", "attned")
	queryValues.Add("gphotoms", "q_value")
	queryValues.Add("gphotoms", "q_unc")
	queryValues.Add("gphotoms", "ned_value")
	queryValues.Add("gphotoms", "ned_unc")
	queryValues.Add("diamdat", "ned_maj_dia")
	queryValues.Add("distance", "avg")
	queryValues.Add("distance", "stddev_samp")

	resp, err := client.R().SetQueryParamsFromValues(queryValues).Get("https://ned.ipac.caltech.edu/cgi-bin/gmd")

	if err != nil {
		return nil, err
	}

	responseString := resp.String()

	return parseResponse(responseString)
}
