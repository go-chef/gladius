package chef

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
	// "strings"
	// "reflect"
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
		"name":     "foo",
		"run_list": []string{"base", "java"},
	}
	// spew.Dump(n1)

	tf, _ := ioutil.TempFile("", "node-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Node can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(n1)

	if b, err := ioutil.ReadAll(n1); err != nil {
		// spew.Dump(b)
		t.Error("error during read from Node", b, err)
	}

	if _, err := NodeFromFile(tf.Name()); err != nil {
		t.Error("could not reserialize node from tempfile after writing it", err)
	} /*else {
		spew.Dump(node)
	}*/
}

func TestRoleWriteToFile(t *testing.T) {
	r1 := &Role{
		"name":       "foo",
		"json_class": "Chef::Role",
		"chef_type":  "role",
		// "default_attributes": []byte{
		// 	"ntp": {"servers": {"ntp1", "ntp2"}},
		// },
		"run_list": []string{"base", "java"},
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

	// eq := reflect.DeepEqual(*r1, r2)
	// t.Error(eq)
	// // if eq {
	// 	t.Error("Imported role does not match exported role")
	// }
	// } else {
	// 	fmt.Println("They're unequal.")
	// }
	// if r1 != r2 {
	// 	t.Error("Imported role does not match exported role")
	// 	t.Error(r1)
	// 	t.Error(&r2)
	// }

	if r2[`name`] != "foo" {
		t.Error("Role name is incorrect")
	}

	// var rl1, rl2 RunList
	// err = mapstructure.Decode(r2, &rl1)
	// err = mapstructure.Decode(r1, &rl2)
	// if rl1 != rl2 {
	// 	t.Error("run_lists do not match")
	// }
	// t.Error(rl1.Items)
	// t.Error(rl2.Items)
	// if reflect.DeepEqual(rl1.Items, rl2.Items) {
	// if rl1.Items == rl2.Items {
	// 	t.Error("run_lists do not match")
	// }
	// t.Error(r1[`name`])
	var r RunList
	err = mapstructure.Decode(r2, r)
	for _, v := range r.Items {
		if v != "base" && v != "java" {
			t.Error("run_list item is invalid: " + v)
		}
	}
}
