package chef

import (
	"github.com/davecgh/go-spew/spew"
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

//func TestNodeWriteToFile(t *testing.T) {
//	n1 := Node{
//		"name":     "foo",
//		"run_list": []string{"base", "java"},
//	}
//	n1.ToFile("path")
//}
