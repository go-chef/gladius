package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bigkraig/chef"
)

type ChefServer struct {
	ServerURL string `json:"server_url"`
	NodeName  string `json:"node_name"`
	ClientKey string `json:"client_key"`
	SkipSSL   bool   `json:"skip_ssl_verify"`
}

type Configuration struct {
	GitLabConfiguration      `json:"gitlab"`
	TestKitchenConfiguration `json:"kitchen"`
	ChefServers              []ChefServer `json:"chef_servers"`
}

type GitLabConfiguration struct {
	APIURL    string `json:"gitlab_api_url"`
	APISecret string `json:"gitlab_api_secret"`
}

// We may need to add support for test kitchen suites
type TestKitchenConfiguration struct {
	Driver      `json:"driver"`
	Provisioner `json:"provisioner"`
	Platforms   []Platform `json:"platforms"`
}

type Driver struct {
	Name              string   `yaml:"name,omitempty" json:"name,omitempty"`
	OpenStackUsername string   `yaml:"openstack_username,omitempty" json:"openstack_username,omitempty"`
	ApiKey            string   `yaml:"openstack_api_key,omitempty" json:"openstack_api_key,omitempty"`
	AuthUrl           string   `yaml:"openstack_auth_url,omitempty" json:"openstack_auth_url,omitempty"`
	NetworkRef        []string `yaml:"network_ref,omitempty" json:"network_ref,omitempty"`
	FlavorRef         string   `yaml:"flavor_ref,omitempty" json:"flavor_ref,omitempty"`
	KeyName           string   `yaml:"key_name,omitempty" json:"key_name,omitempty"`
	ServerName        string   `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	ImageRef          string   `yaml:"image_ref,omitempty" json:"image_ref,omitempty"`
	Username          string   `yaml:"username,omitempty" json:"username,omitempty"`
}

type Provisioner struct {
	RequireChefOmnibus bool   `yaml:"require_chef_omnibus" json:"require_chef_omnibus"`
	ChefOmnibusUrl     string `yaml:"chef_omnibus_url" json:"chef_omnibus_url"`
}

type Platform struct {
	Name         string `yaml:"name" json:"name"`
	DriverConfig Driver `yaml:"driver_config" json:"driver_config"`
}

func mockConfiguration() *Configuration {
	return &Configuration{
		GitLabConfiguration: GitLabConfiguration{
			APIURL:    "http://gitlab.example.com/api/v3",
			APISecret: "API SECRET KEY",
		},
		ChefServers: []ChefServer{
			{
				ServerURL: "https://chef-colo1.example.com",
				NodeName:  "jenkins",
				ClientKey: "CLIENT PEM",
				SkipSSL:   true,
			},
			{
				ServerURL: "https://chef-colo2.example.com",
				NodeName:  "jenkins",
				ClientKey: "CLIENT PEM",
				SkipSSL:   true,
			},
		},
		TestKitchenConfiguration: TestKitchenConfiguration{
			Driver: Driver{
				Name:              "openstack",
				OpenStackUsername: "openstack_username",
				Username:          "vm_username",
				ApiKey:            "API KEY",
				AuthUrl:           "http://openstack.example.com:5000/v2.0/tokens",
				NetworkRef:        []string{"a5eecb4d-a5f4-4e8e-975a-2c5972c5f880"},
				FlavorRef:         "m1.small",
				KeyName:           "ssh_key_name",
			},
			Provisioner: Provisioner{
				ChefOmnibusUrl:     "http://repo.example.com/chef/gochefyourselflite.sh",
				RequireChefOmnibus: true,
			},
			Platforms: []Platform{
				{
					Name: "ubuntu-14.04",
					DriverConfig: Driver{
						ImageRef: "f0527ea9-a025-4249-ac79-3e7e0883f463",
						Username: "ubuntu",
					},
				},
				{
					Name: "oel-7.0",
					DriverConfig: Driver{
						ImageRef: "6831bc97-31b3-4887-ae13-52ca046b65c2",
						Username: "cloud-user",
					},
				},
			},
		},
	}
}

func ReadConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	homedir := os.Getenv("HOME")
	file, err := os.Open(filepath.Join(homedir, ".gladius.json"))
	if err != nil {
		json, _ := json.MarshalIndent(mockConfiguration(), "", "    ")
		err = errors.New(fmt.Sprintf("Please create a ~/.gladius.json containing\n%s",
			json))
		return cfg, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfg)

	return cfg, err
}

func (c *Configuration) GenerateChefClients() ([]*chef.Client, error) {
	chefServers := make([]*chef.Client, len(c.ChefServers))

	for i, chefConfig := range c.ChefServers {
		c, err := chef.NewClient(&chef.Config{
			Name:    chefConfig.NodeName,
			Key:     chefConfig.ClientKey,
			BaseURL: chefConfig.ServerURL,
			SkipSSL: chefConfig.SkipSSL,
		})
		if err != nil {
			return chefServers, err
		}
		chefServers[i] = c
	}
	return chefServers, nil
}
