package util

import (
	"fmt"
	"net/http"

	"golang.org/x/net/proxy"
)

type TorClient struct {
	*http.Client
}

func NewTorClient() (*TorClient, error) {
	p, err := proxy.SOCKS5("tcp", "tor:9050", nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to open Tor proxy: %v", err)
	}

	var client TorClient = TorClient{http.DefaultClient}
	client.Transport = &http.Transport{
		Dial: p.Dial,
	}

	return &client, nil
}
