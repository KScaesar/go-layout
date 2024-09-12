package adapters

import (
	"net/http"
)

func NewHttpClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConnsPerHost: 5,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}
