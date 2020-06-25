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
	"torobSpider/rotator"
)

var UAs = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36","Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0","Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36","Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"}
var cookieJar, _ = cookiejar.New(nil)

func getClient(useProxy bool) (*http.Client, *rotator.Proxy) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar: cookieJar,
	}
	if useProxy {
		proxy := CurrentRuntimeInfo.ProxyRotator.GetProxySync()
		client.Transport = &http.Transport{Proxy: http.ProxyURL(&proxy.Url)}
		return client, proxy
	}
	return client, nil
}

func DeleteCookie() {
	cookieJar, _ = cookiejar.New(nil)
}

func getJson(url string, target interface{}) error {
	CurrentRuntimeInfo.WorkerPool <- 1
	defer func() {
		<-CurrentRuntimeInfo.WorkerPool
	}()
	time.Sleep((time.Duration(rand.Intn(5)) + 1) * time.Second)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UAs[rand.Intn(3)])
	client, proxy := getClient(false)
	r, err := client.Do(req)
	if err != nil {
		if proxy != nil {
			proxy.MarkDead()
		}
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
	defer func() {
		<-CurrentRuntimeInfo.WorkerPool
	}()
	time.Sleep((time.Duration(rand.Intn(5)) + 2) * time.Second)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UAs[rand.Intn(3)])
	client, proxy := getClient(true)
	r, err := client.Do(req)
	if err != nil {
		if proxy != nil {
			proxy.MarkDead()
		}
		return "", err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if proxy != nil {
			proxy.MarkDead()
		}
		return "", err
	}
	return string(data), nil
}

func getFakeText(url string) (string, error) {
	time.Sleep((time.Duration(rand.Intn(7)) + 2) * time.Second)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UAs[rand.Intn(3)])
	client, proxy := getClient(true)
	r, err := client.Do(req)
	if err != nil {
		if proxy != nil {
			proxy.MarkDead()
		}
		return "", err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if proxy != nil {
			proxy.MarkDead()
		}
		return "", err
	}
	return string(data), nil
}


func Base64Decode(data string) string {
	sDec, _ := b64.StdEncoding.DecodeString(data)
	return string(sDec)
}
