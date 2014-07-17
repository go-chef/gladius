package lib

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/go-chef/gladius/app"
	"gopkg.in/yaml.v1"
)

const TestKitchenLocalFilename = ".kitchen.local.yml"

func GenerateTestKitchenConfiguration(cfg *app.Configuration) error {
	kitchen := &app.TestKitchenConfiguration{
		Driver:      cfg.Driver,
		Provisioner: cfg.Provisioner,
		Platforms:   cfg.Platforms,
	}

	kitchenYAML, err := yaml.Marshal(&kitchen)
	if err != nil {
		return err
	}

	filename := filepath.Join(TestKitchenLocalFilename)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(kitchenYAML))
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
