package chef

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	testNodeJSON = "test/node.json"
	// FML
	testNodeMapStringInterfaceLol = NewNode(&Reader{
		"name":       "test",
		"run_list":   []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		"chef_type":  "node",
		"json_class": "Chef::Node",
	})
)

func TestNodeName(t *testing.T) {
	n1 := testNodeMapStringInterfaceLol // (*Node)
	name, err := n1.Name()
	if err != nil {
		t.Error("unexpected error from Name() call", err)
	} else if name != "test" {
		t.Error("use a fucking assertion library you plonker")
	}
}

func TestNodeFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testNodeJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testNodeJSON)
	} else {
		dec := json.NewDecoder(file)
		var n Node
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// TestNewNode checks the NewNode Reader chain for Type
func TestNewNode(t *testing.T) {
	var v interface{}
	v = testNodeMapStringInterfaceLol
	switch v.(type) {
	case *Node:
		t.Log(v, "was correctly identified as pointer to Node type")
	default:
		spew.Dump("v", v)
		t.Error(v, "was not a Node type")
	}
}

// TestNodeReadIntoFile tests that Read() can be used to read by io.Readers
func TestNodeReadIntoFile(t *testing.T) {
	n1 := testNodeMapStringInterfaceLol // (*Node)
	spew.Dump(n1)

	tf, _ := ioutil.TempFile("test", "node-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())
	// Copy to tempfile (I use Read() internally)
	io.Copy(tf, n1)
}
