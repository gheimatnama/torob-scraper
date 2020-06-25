package rotator

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"strings"
)

func ParseProxyFile() []url.URL {
	var proxyUrls []url.URL
	file, err := os.Open("proxies.txt")
	if err != nil {
		return proxyUrls
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		u, _ := url.Parse("http://" +  text)
		proxyUrls = append(proxyUrls, *u)
	}
	logrus.Info("Parsed file found ", len(proxyUrls))
	return proxyUrls
}