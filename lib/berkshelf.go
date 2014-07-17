package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-chef/gladius/app"
)

type BerkshelfChefConfiguration struct {
	ChefServerURL string `json:"chef_server_url"`
	NodeName      string `json:"node_name"`
	ClientKey     string `json:"client_key"`
}

type BerkshelfSSLConfiguration struct {
	Verify bool `json:"verify"`
}

type BerkshelfConfiguration struct {
	BerkshelfChefConfiguration `json:"chef"`
	BerkshelfSSLConfiguration  `json:"ssl"`
}

func NewBerkshelfConfiguration(server *app.ChefServer) *BerkshelfConfiguration {
	cfg := &BerkshelfConfiguration{
		BerkshelfChefConfiguration: BerkshelfChefConfiguration{
			ChefServerURL: server.ServerURL,
			NodeName:      server.NodeName,
			ClientKey:     server.ClientKey,
		},
		BerkshelfSSLConfiguration: BerkshelfSSLConfiguration{
			Verify: !server.SkipSSL,
		},
	}
	return cfg
}

func NeedBerkshelfInstall(path string) bool {
	_, err := os.Stat(fmt.Sprintf("%s%c%s", path, os.PathSeparator, "Berksfile.lock"))
	if err != nil {
		return true
	}
	return false
}

func GenerateBerkshelfConfiguration(path string, chefServer *app.ChefServer) (string, error) {
	berks := NewBerkshelfConfiguration(chefServer)
	berksJSON, _ := json.MarshalIndent(berks, "", "    ")

	filename := filepath.Join(path, ".berksconfig.json")

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(berksJSON))
	if err != nil {
		return "", err
	}

	return filename, nil
}
