package chef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type RunList struct {
	Items []string `mapstructure:"run_list"`
}

func TestNodeFromFile(t *testing.T) {
	n1, err := NodeFromFile("test/node.json")
	if err != nil {
		t.Fatal(err)
	}

	if n1[`name`] != "testnode" {
		t.Error("Node name is incorrect")
	}
}

func TestRoleFromFile(t *testing.T) {
	r1, err := RoleFromFile("test/role.json")
	if err != nil {
		t.Fatal(err)
	}

	if r1[`name`] != "testrole" {
		t.Error("Role name is incorrect")
	}
}

func TestApiClientFromFile(t *testing.T) {
	c1, err := ApiClientFromFile("test/client.json")
	if err != nil {
		t.Fatal(err)
	}

	if c1[`name`] != "testclient" {
		t.Error("Client name is incorrect")
	}
}

func TestEnvironmentFromFile(t *testing.T) {
	e1, err := EnvironmentFromFile("test/environment.json")
	if err != nil {
		t.Fatal(err)
	}

	if e1[`name`] != "testenvironment" {
		t.Error("Environment name is incorrect")
	}
}

func TestDatabagFromFile(t *testing.T) {
	d1, err := DatabagFromFile("test/databag.json")
	if err != nil {
		t.Fatal(err)
	}

	if d1[`id`] != "testdatabag" {
		t.Error("Databag name is incorrect")
	}
}

func TestNodeWriteToFile(t *testing.T) {
	n1 := &Node{
		"chef_type":  "node",
		"json_class": "Chef::Node",
		"name":       "foo",
		"run_list":   []string{"base", "java"},
	}

	tf, _ := ioutil.TempFile("", "node-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Node can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(n1)

	if b, err := ioutil.ReadAll(n1); err != nil {
		t.Error("error during read from Node", b, err)
	}

	n2, err := NodeFromFile(tf.Name())
	if err != nil {
		t.Error("could not reserialize node from tempfile after writing it", err)
	}

	if n2[`name`] != "foo" {
		t.Error("Role name is incorrect")
	}

	if reflect.DeepEqual(*n1, n2) {
		t.Error("Starting node does not match read in node")
	}
}

func TestRoleWriteToFile(t *testing.T) {
	r1 := &Role{
		"chef_type": "role",
		"default_attributes": map[string]string{
			"ntp": "something",
			"zoo": "pandabear",
		},
		"json_class": "Chef::Role",
		"name":       "foo",
		"run_list":   []string{"base", "java"},
	}

	tf, _ := ioutil.TempFile("", "role-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Role can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(r1)

	if b, err := ioutil.ReadAll(r1); err != nil {
		t.Error("error during read from role", b, err)
	}

	r2, err := RoleFromFile(tf.Name())
	if err != nil {
		t.Error("could not reserialize role from tempfile after writing it", err)
	}

	if r2[`name`] != "foo" {
		t.Error("Role name is incorrect")
	}

	if reflect.DeepEqual(*r1, r2) {
		t.Error("Starting role does not match read in role")
	}
}

func TestEnvironmentWriteToFile(t *testing.T) {
	e1 := &Environment{
		"chef_type": "environment",
		"default_attributes": map[string]string{
			"ntp": "something",
			"zoo": "pandabear",
		},
		"json_class": "Chef::Environment",
		"name":       "foo",
	}

	tf, _ := ioutil.TempFile("", "environment-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Environment can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(e1)

	if b, err := ioutil.ReadAll(e1); err != nil {
		t.Error("error during read from environment", b, err)
	}

	e2, err := EnvironmentFromFile(tf.Name())
	if err != nil {
		t.Error("could not reserialize environment from tempfile after writing it", err)
	}

	if e2[`name`] != "foo" {
		t.Error("Environment name is incorrect")
	}

	if reflect.DeepEqual(*e1, e2) {
		t.Error("Starting environment does not match read in environment")
	}
}

func TestApiClientWriteToFile(t *testing.T) {
	c1 := &ApiClient{
		"admin":      true,
		"chef_type":  "client",
		"json_class": "Chef::ApiClient",
		"name":       "foo",
		// Reallly interesting
		// If we remove run_list for a proper client
		// DeepEqual = false
		// When it's in here
		// DeepEqual = true
		"run_list":  []string{"base", "java"},
		"validator": false,
	}

	tf, _ := ioutil.TempFile("", "client-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// ApiClient can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(c1)

	if b, err := ioutil.ReadAll(c1); err != nil {
		t.Error("error during read from client", b, err)
	}

	c2, err := ApiClientFromFile(tf.Name())
	if err != nil {
		t.Error("could not reserialize client from tempfile after writing it", err)
	}

	if c2[`name`] != "foo" {
		t.Error("ApiClient name is incorrect")
	}

	if reflect.DeepEqual(*c1, c2) {
		t.Error("Starting client does not match read in client")
		t.Error(reflect.TypeOf(*c1))
		t.Error(fmt.Sprintf("%v", *c1))
		t.Error(reflect.TypeOf(c2))
		t.Error(fmt.Sprintf("%v", c2))
	}
}

func TestDatabagWriteToFile(t *testing.T) {
	d1 := &Databag{
		"chef_type":  "databag",
		"json_class": "Chef::Databag",
		"id":         "foo",
		"data":       []string{"somestuff", "morestuff"},
	}

	tf, _ := ioutil.TempFile("", "databag-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Databag can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(d1)

	if b, err := ioutil.ReadAll(d1); err != nil {
		t.Error("error during read from databag", b, err)
	}

	d2, err := DatabagFromFile(tf.Name())
	if err != nil {
		t.Error("could not reserialize databag from tempfile after writing it", err)
	}

	if d2[`id`] != "foo" {
		t.Error("Databag name is incorrect")
	}

	if reflect.DeepEqual(*d1, d2) {
		t.Error("Starting databag does not match read in databag")
	}
}
