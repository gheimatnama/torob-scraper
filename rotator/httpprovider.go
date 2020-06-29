package rotator

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func ParseScyllaProxies() []url.URL {
	var proxyUrls []url.URL
	myClient := &http.Client{Timeout: time.Second * 10}
	resp, err := myClient.Get("http://localhost:8899/api/v1/proxies?limit=500")
	if err != nil {
		logrus.Error(err)
		return proxyUrls
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return proxyUrls
		logrus.Error(err)
	}
	resp.Body.Close()
	var responseJson map[string]interface{}
	json.Unmarshal(body, &responseJson)
	proxies := responseJson["proxies"].([]interface{})
	for _, proxyRaw := range proxies {
		proxy := proxyRaw.(map[string]interface{})
		ip := proxy["ip"].(string)
		port := fmt.Sprintf("%d", int(proxy["port"].(float64)))
		u, _ := url.Parse("http://" + ip + ":" + port)
		proxyUrls = append(proxyUrls, *u)
	}
	logrus.Info("Parsed urls from scylla found : ", len(proxyUrls))
	return proxyUrls
}