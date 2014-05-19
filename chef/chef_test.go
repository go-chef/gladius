package chef

import (
	"encoding/json"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	//  "io"
	"io/ioutil"
	"os"
	"testing"
)

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
