package chef

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// ClientConfig contains the types used for authentication with the Chef Server
// Uri: The URL of the chef server
// Key: The rsa Private Key we authenticate with
// Name: The username that matches the private key we use
type ClientConfig struct {
	Uri  *url.URL
	Key  *rsa.PrivateKey
	Name string
}

// Client combines a standard http Request with the ClientConfig structure
type Client struct {
	Request *http.Request
	ClientConfig
}

// Object is the basic chef object type (node/role/client/env....)
type Object struct{}

// Object Reader is the RO interface to Objects
type ObjectReader interface {
	Read(p []byte) (size int, err error)
}

// ObjectWriter is the WO interface to Objects
type ObjectWriter interface {
	Write()
}

// ObjectReadWriter is an interface type combining the ObjectReader and Writer used for modification
type ObjectReadWriter struct {
	ObjectReader
	ObjectWriter
}

type ArbitraryJson map[string]interface{}

// ApiClient is presently arbitrary json data
type ApiClient ArbitraryJson

// Cookbook is presently arbitrary json data
type Cookbook ArbitraryJson

// Databag is presently arbitrary json data
type Databag ArbitraryJson

// Environment is presently arbitrary json data
type Environment ArbitraryJson

// Node is presently arbitrary json data
type Node ArbitraryJson

// Role is presently arbitrary json data
type Role ArbitraryJson

//This will probably end up being a private read() and delegating to Object Read
func (c *ApiClient) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&c); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

//This will probably end up being a private read() and delegating to Object Read
func (c *Cookbook) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&c); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

//This will probably end up being a private read() and delegating to Object Read
func (d *Databag) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&d); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

//This will probably end up being a private read() and delegating to Object Read
func (e *Environment) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&e); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

//This will probably end up being a private read() and delegating to Object Read
func (r *Role) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&r); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

//This will probably end up being a private read() and delegating to Object Read
func (n *Node) Read(p []byte) (size int, err error) {
	if b, err := json.Marshal(&n); err == nil {
		copy(p, b)
		return len(p), io.EOF
	}
	return len(p), nil
}

var chefTypeMap = map[string]interface{}{
	`Chef::ApiClient`:   ApiClient{},
	`Chef::Cookbook`:    Cookbook{},
	`Chef::Databag`:     Databag{},
	`Chef::Environment`: Environment{},
	`Chef::Node`:        Node{},
	`Chef::Role`:        Role{},
}

// ErrUninferableType allows an implementer to gracefully handle maybeChefType
// When the type is not inferable
var ErrUninferableType = errors.New("could not infer type from chef_type or json_class keys")

// This thing might figure out what chef type some json amalgamet is
// example:
//   json, err :=  jsonDecoder(buff)
//   var t = maybeChefType(json)
//   switch t :=  t.(type) {
//   default:
//     //unexpected type
//   Node:
//     NodeDoSOmething()
//   }
func maybeChefType(obj map[string]interface{}) (interface{}, error) {
	for key := range obj {
		switch key {
		case "chef_type", "json_class":
			var ok bool
			var maybeType interface{}
			if maybeType, ok = chefTypeMap[key]; ok {
				return maybeType, nil
			}
			return nil, ErrUninferableType
		}
	}
	return nil, nil
}
