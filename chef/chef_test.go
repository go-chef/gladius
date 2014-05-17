package chef

import (
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNodeFromFile(t *testing.T) {
	n1, err := NodeFromFile("test/node.json")
	spew.Dump(n1)
	if err != nil {
		spew.Dump(err)
		t.Fatal(err)
	}
}

func TestNodeWriteToFile(t *testing.T) {
	n1 := &Node{
		"name":     "foo",
		"run_list": []string{"base", "java"},
	}

	tf, _ := ioutil.TempFile("test", "gladius-chef-node")
	defer tf.Close()
	defer os.Remove(tf.Name())

	spew.Dump(n1)

	// try to read n1 into b
	var b []byte
	if _, err := n1.Read(b); err != nil {
		spew.Dump(b)
	}

	// because Node has a io.Reader Read() compliant implementation, we can copy out of it
	// This hangs -- why?
	if _, err := io.Copy(tf, n1); err != nil {
		t.Error("could not copy node into tempfile", err)
	}

	if _, err := NodeFromFile(tf.Name()); err != nil {
		t.Error("could not reserialize node from tempfile after writing it", err)
	}

}
