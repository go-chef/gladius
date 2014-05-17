package chef

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	Uri  *url.URL
	Key  *rsa.PrivateKey
	Name string
}

type Client struct {
	Request http.Request
	ClientConfig
}

// Object is the basic chef object type (node/role/client/env....)
type Object struct{}

// Object Reader is the RO interface to Objects
type ObjectReader interface {
	Read()
}

// ObjectWriter is the WO interface to Objects
type ObjectWriter interface {
	Write()
}

// chef.Object is any chef object (
type ObjectReadWriter interface {
	ObjectReader
	ObjectWriter
}

type Node map[string]interface{}

func (n *Node) Read() error {
	// Append into nodes
	// Use makerequest to pull
	// MAGIC
	return nil
}

// where does it downlad too ?
//  memory returns the buff then you can do what you will
func Download(orw *ObjectReadWriter) *Object {
	//MAGIC
	// Example from gladius
	//  n  := chef.NewNode("NodeName")
	//  chef.Upload(n)
	return &Object{}
}

// NodeReader reads nodes from io buffers and creates a Node
func jsonDecoder(buff io.Reader) (data map[string]interface{}, err error) {
	err = json.NewDecoder(buff).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

var chefTypeMap = map[string]interface{}{
	`Chef::Node`: Node{},
}

var UninferableType = errors.New("Could not infer type from chef_type or json_class keys")

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
	for key, _ := range obj {
		switch key {
		case "chef_type", "json_class":
			var ok bool
			var maybeType interface{}
			if maybeType, ok = chefTypeMap[key]; ok {
				return maybeType, nil
			}
		}
	}
	return nil, nil
}
