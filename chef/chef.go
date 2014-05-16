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
	Request http.Request
	ClientConfig
}

type Object struct{}

// chef.Object is any chef object (
type ObjectReadWriter interface {
	Read()
	Write()
}

// Dynamic Json data for node
type rawJSON map[string]interface{}

type Node struct {
	Name string `json:"name"`
	rawJSON
}

func (n *Node) Read() error {
	// Example from gladius
	//  n  := chef.NewNode("NodeName")
	//  chef.Upload(n)

	// Append into nodes
	// Use makerequest to pull
	// MAGIC
	return nil
}

// where does it downlad too ?
//  memory returns the buff then you can do what you will
func Download(orw *ObjectReadWriter) *Object {
	//MAGIC
	return &Object{}
}

func NewNode(name string) *Node {
	data := make(rawJSON)
	return &Node{name, data}
}
