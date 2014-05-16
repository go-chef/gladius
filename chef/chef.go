package chef

import (
	"crypto/rsa"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	Uri  *url.URL
	Key  *rsa.PrivateKey
	Name string
}

type Client struct {
	Request *http.Request
	ClientConfig
}
