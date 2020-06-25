package rotator

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxyChecker struct {
	validator func(proxy url.URL) bool
	WorkersSize int
}

func NewCheckerInstance(validator func(proxy url.URL) bool) *ProxyChecker {
	return &ProxyChecker{
		validator: validator,
	}
}


func (proxyChecker *ProxyChecker) IsProxyWorking(url url.URL) bool {
	return proxyChecker.validator(url)
}


func RecaptchaChecker(proxy url.URL) bool {
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&proxy),DisableKeepAlives: true}, Timeout: time.Second * 5}
	resp, err := myClient.Get("https://api.torob.com/v4/product-page/redirect/?prk=5ced2835-f5a2-4fe5-9236-993a5cb04fd5&source=next&uid=&discover_method=search&_bt__experiment=members_clicks_04")
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return !strings.Contains(string(body), "captcha") && !strings.Contains(string(body), "Captcha")
}