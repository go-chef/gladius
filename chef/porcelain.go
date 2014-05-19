package chef

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// NodeFromFile reads, decodes and returns a node from filePath, or an error
func NodeFromFile(filename string) (node Node, err error) {
	// sanitize the path
	filename = path.Clean(filename)
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(b, &node)
	}
	return node, err
}

func RoleFromFile(filename string) (role Role, err error) {
	filename = path.Clean(filename)
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(b, &role)
	}
	return role, err
}

func ApiClientFromFile(filename string) (apiclient ApiClient, err error) {
	filename = path.Clean(filename)
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(b, &apiclient)
	}
	return apiclient, err
}

// Cookbook probably isn't a good fit for this
// func ClientFromFile(filename string) (client Client, err error) {
// 	filename = path.Clean(filename)
// 	var b []byte
// 	if b, err = ioutil.ReadFile(filename); err == nil {
// 		json.Unmarshal(b, &client)
// 	}
// 	return Client, err
// }

func EnvironmentFromFile(filename string) (environment Environment, err error) {
	filename = path.Clean(filename)
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(b, &environment)
	}
	return environment, err
}

func DatabagFromFile(filename string) (databag Databag, err error) {
	filename = path.Clean(filename)
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(b, &databag)
	}
	return databag, err
}
