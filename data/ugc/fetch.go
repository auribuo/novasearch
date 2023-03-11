package ugc

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/auribuo/novasearch/fs"
	"github.com/auribuo/novasearch/log"
	"github.com/go-resty/resty/v2"
)

func Fetch() ([]Response, error) {
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

	cacheFile := filepath.FromSlash(cacheFolder + "/ugc.json")

	log.Default.Debug("trying to fetch UGC data from cache")
	local, err := fetchLocal(cacheFile)
	if err != nil {
		return nil, err
	}
	if len(local) > 0 {
		log.Default.Debug("successfully fetched NED data from cache")
		return local, nil
	}
	log.Default.Warn("no cache found. trying to fetch UGC data from remote")
	data, err := fetchRemote()
	if err != nil {
		return nil, err
	}
	log.Default.Debug("successfully fetched UGC data from remote. saving to cache")
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
		if err != nil {
			return nil, err
		}
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
	if time.Since(cache.LastUpdated) < 365*time.Hour*24 {
		return cache.Items, nil
	}
	return nil, nil
}

func fetchRemote() ([]Response, error) {
	client := resty.New()

	bodyValues := url.Values{}
	bodyValues.Add("-ref", "VIZ64060aaa23fbb0")
	bodyValues.Add("-to", "4")
	bodyValues.Add("-from", "-3")
	bodyValues.Add("-this", "-3")
	bodyValues.Add("//source", "VII/26D")
	bodyValues.Add("//tables", "VII/26D/catalog")
	bodyValues.Add("//tables", "VII/26D/errors")
	bodyValues.Add("-out.max", "unlimited")
	bodyValues.Add("//CDSportal", "http://cdsportal.u-strasbg.fr/StoreVizierData.html")
	bodyValues.Add("-out.form", "ascii+text/plain")
	bodyValues.Add("-out.add", "_RAJ,_DEJ")
	bodyValues.Add("//outaddvalue", "default")
	bodyValues.Add("-order", "I")
	bodyValues.Add("-oc.form", "sexa")
	bodyValues.Add("-out.src", "VII/26D/catalog,VII/26D/errors")
	bodyValues.Add("-nav", "cat:VII/26D&tab:{VII/26D/catalog}&tab:{VII/26D/errors}&key:source=VII/26D&HTTPPRM:&")
	bodyValues.Add("-c", "")
	bodyValues.Add("-c.eq", "J2000")
	bodyValues.Add("-c.r", "++2")
	bodyValues.Add("-c.u", "arcmin")
	bodyValues.Add("-c.geom", "r")
	bodyValues.Add("-source", "")
	bodyValues.Add("-source", "VII/26D/catalog VII/26D/errors")
	bodyValues.Add("-out", "UGC")
	bodyValues.Add("-out", "MajAxis")
	bodyValues.Add("-out", "MinAxis")
	bodyValues.Add("-out", "PA")
	bodyValues.Add("-out", "Hubble")
	bodyValues.Add("-out", "Pmag")
	bodyValues.Add("-out", "i")
	bodyValues.Add("-meta.ucd", "1")
	bodyValues.Add("-meta", "1")
	bodyValues.Add("-meta.foot", "1")
	bodyValues.Add("-usenav", "1")
	bodyValues.Add("-bmark", "POST")
	body := bodyValues.Encode()
	resp, err := client.R().SetBody(body).Post("https://vizier.unistra.fr/viz-bin/asu-txt")
	if err != nil {
		return nil, err
	}
	responseString := resp.String()
	return parseResponse(responseString)
}
