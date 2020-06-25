package rotator

import (
	"net/url"
	"time"
)

type ProxyProvider struct {
	Provider func() []url.URL
	RefreshRate time.Duration // Every duration
	Meta map[string]interface{}
}

func NewProviderInstance(provider func() []url.URL, refreshRate time.Duration) *ProxyProvider {
	return &ProxyProvider{
		Provider: provider,
		RefreshRate: refreshRate,
		Meta: make(map[string]interface{}),
	}
}
