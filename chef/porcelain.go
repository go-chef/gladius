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
