package chef

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	//  "io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNodeFromFile(t *testing.T) {
	if n1, err := NodeFromFile("test/node.json"); err != nil {
		t.Fatal(err)
	} else {
		spew.Dump(n1)
	}
}

func TestNodeWriteToFile(t *testing.T) {
	n1 := &Node{
		"name":     "foo",
		"run_list": []string{"base", "java"},
	}
	spew.Dump(n1)

	tf, _ := ioutil.TempFile("", "node-to-file")
	defer tf.Close()
	defer os.Remove(tf.Name())

	// Node can just be Encoded directly now that it implements Read()
	enc := json.NewEncoder(tf)
	enc.Encode(n1)

	if b, err := ioutil.ReadAll(n1); err != nil {
		spew.Dump(b)
		t.Error("error during read from Node", b, err)
	}

	if node, err := NodeFromFile(tf.Name()); err != nil {
		t.Error("could not reserialize node from tempfile after writing it", err)
	} else {
		spew.Dump(node)
	}

}
