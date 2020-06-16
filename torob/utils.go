package torob

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var cookieJar, _ = cookiejar.New(nil)
var myClient = &http.Client{
	Timeout: 10 * time.Second,
	Jar: cookieJar,
}

func getJson(url string, target interface{}) error {
	CurrentRuntimeInfo.WorkerPool <- 1
	time.Sleep(time.Duration(rand.Intn(500) + 1000)* time.Millisecond)
	defer func() {
		<-CurrentRuntimeInfo.WorkerPool
	}()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0")
	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func StripQueryParam(inURL string, stripKey string) string {
	u, err := url.Parse(inURL)
	if err != nil {
		return inURL
	}
	q := u.Query()
	q.Del(stripKey)
	u.RawQuery = q.Encode()
	return u.String()
}

func CleanProductUrl(url string) string {
	url = StripQueryParam(url, "source")
	url = StripQueryParam(url, "discover_method")
	url = StripQueryParam(url, "experiment")
	return url
}

func GetHostName(address string) string {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}
	return u.Hostname()
}


func GetQueryParam(address string, param string) []string {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	return m[param]
}

func getText(url string) (string, error) {
	CurrentRuntimeInfo.WorkerPool <- 1
	time.Sleep(time.Duration(rand.Intn(2000) + 1000)* time.Millisecond)
	defer func() {
		<-CurrentRuntimeInfo.WorkerPool
	}()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0")
	r, err := myClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}


func Base64Decode(data string) string {
	sDec, _ := b64.StdEncoding.DecodeString(data)
	return string(sDec)
}
