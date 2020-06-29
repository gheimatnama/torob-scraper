package rotator

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"sort"
	"time"
)


type (
	ProxyRotator struct {
		proxies []*Proxy
		proxyChecker *ProxyChecker
		providers []*ProxyProvider
		ProxyConnectionDelay time.Duration // Connections will not be concurrent in a proxy
		ParallelProxyConnection bool // Using a proxy multiple time in request
		TotalParallelProxyConnection uint
		CheckProxyInterval time.Duration
		CheckProxyBeforeConnection bool
		checkerPool chan int
		ProxyQueueRetryTimeout time.Duration // Client will sleep in this time, then rotator tries to find a proxy, this cycle goes on
		ProxyQueueTimeout time.Duration // Total time client waits for a proxy
	}

	Proxy struct {
		Url url.URL
		CurrentConnections uint
		TotalConnections uint
		IsDead bool
		LastConnectionTime time.Time
		LastTimeChecked time.Time
	}
)

func (proxy *Proxy) IsWorking(checker *ProxyChecker) bool {
	return checker.IsProxyWorking(proxy.Url)
}

func (proxy *Proxy) MarkDead() {
	proxy.IsDead = true
}


func (proxy *Proxy) BeenCheckedSince(duration time.Duration) bool {
	return time.Now().Sub(proxy.LastTimeChecked) < duration
}

func NewInstance(proxyChecker *ProxyChecker, providers []*ProxyProvider) *ProxyRotator {
	return &ProxyRotator{providers: providers, proxyChecker : proxyChecker, checkerPool: make(chan int, 30)}
}

func NewProxy(url url.URL) *Proxy {
	return &Proxy{Url: url, IsDead: true}
}


func (proxyRotator *ProxyRotator) hasProxy(url url.URL) bool {
	for _, proxy := range proxyRotator.proxies {
		if url == proxy.Url {
			return true
		}
	}
	return false
}

func (proxyRotator *ProxyRotator) addProxy(url url.URL)  {
	if !proxyRotator.hasProxy(url) {
		proxyRotator.proxies = append(proxyRotator.proxies, NewProxy(url))
	}
}

func (proxyRotator *ProxyRotator) addProxies(urls []url.URL)  {
	for _, url := range urls {
		proxyRotator.addProxy(url)
	}
}


func (proxyRotator *ProxyRotator) Init(minimumAvailableProxies int) {
	proxyRotator.initializeProviders()
	proxyRotator.initializeProxyChecker()
	for minimumAvailableProxies >= proxyRotator.TotalAliveProxies() {
		fmt.Print("Current available proxies : ", proxyRotator.TotalAliveProxies(), "\r")
		time.Sleep(200 * time.Millisecond)
	}
	logrus.Info("Proxy pool reached required conditions")
}

func (proxyRotator *ProxyRotator) initializeProviders()  {
	go func() {
		for _, prv := range proxyRotator.providers {
			go func(proxyRotator *ProxyRotator, provider *ProxyProvider) {
				for {
					urls := provider.Provider()
					proxyRotator.addProxies(urls)
					time.Sleep(provider.RefreshRate)
				}
			} (proxyRotator, prv)
		}
	}()
}


func (proxyRotator *ProxyRotator) initializeProxyChecker() {
	go func(proxyRotator *ProxyRotator) {
		for {
			if proxyRotator.proxies == nil {
				continue
			}
			for _, proxy := range proxyRotator.proxies {
				if !proxy.BeenCheckedSince(proxyRotator.CheckProxyInterval) {
					proxy.LastTimeChecked = time.Now()
					proxyRotator.checkerPool <- 1
					go func(proxy *Proxy, proxyRotator *ProxyRotator) {
						proxy.IsDead = !proxy.IsWorking(proxyRotator.proxyChecker)
						<- proxyRotator.checkerPool
					}(proxy, proxyRotator)
				}
			}
		}
	}(proxyRotator)
}

func (proxyRotator *ProxyRotator) TotalAliveProxies() int {
	count := 0
	for _, proxy := range proxyRotator.proxies {
		if !proxy.IsDead {
			count++
		}
	}
	return count
}


func (proxyRotator *ProxyRotator) GetProxySync() *Proxy {
	logrus.Info("Total alive proxies : ", proxyRotator.TotalAliveProxies())
	startedAt := time.Now()
	var workingProxy *Proxy
	for {
		proxyRotator.sortProxiesByTime()
		if proxyRotator.proxies == nil {
			time.Sleep(proxyRotator.ProxyQueueRetryTimeout)
			continue
		}
		for _, proxy := range proxyRotator.proxies {
			if proxy.IsDead {
				continue
			}
			if proxyRotator.ParallelProxyConnection && proxy.CurrentConnections < proxyRotator.TotalParallelProxyConnection {
				workingProxy = proxy
			} else {
				if time.Now().Sub(proxy.LastConnectionTime) > proxyRotator.ProxyConnectionDelay {
					workingProxy = proxy
				}
			}
			if workingProxy != nil && proxyRotator.CheckProxyBeforeConnection {
				if !workingProxy.IsWorking(proxyRotator.proxyChecker) {
					workingProxy = nil
				}
			}
		}
		if workingProxy == nil {
			if time.Now().Sub(startedAt) > proxyRotator.ProxyQueueTimeout {
				return nil
			}
			time.Sleep(proxyRotator.ProxyQueueRetryTimeout)
		} else {
			break
		}
	}
	workingProxy.TotalConnections++
	workingProxy.LastConnectionTime = time.Now()
	return workingProxy
}


func (proxyRotator *ProxyRotator) sortProxiesByTime() {
	sort.Slice(proxyRotator.proxies, func(i, j int) bool {
		return proxyRotator.proxies[i].LastConnectionTime.Before(proxyRotator.proxies[j].LastConnectionTime)
	})
}