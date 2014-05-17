package chef

import (
	"os"
	"path"
)

func NodeFromFile(filePath string) (node Node, err error) {
	// sanitize the path
	filePath = path.Clean(filePath)
	fileBuff, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	node, err = jsonDecoder(fileBuff)
	return node, err
}
