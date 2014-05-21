package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Node has a Reader, hey presto
type Node struct {
	*Reader
}

// NativeNode represents the native Go version of the deserialized Node type
// BUG(fujin): Add other destructing fields like run_list, override attributes
type NativeNode struct {
	Name string `mapstructure:"name"`
}

// Name method is pretty cool if you like giving names to stuff
// Declare a temporary NativeNode, decode the Reader into it and return a copy of the nodes Name
func (n *Node) Name() (name string, err error) {
	var node NativeNode
	return node.Name, mapstructure.Decode(n.Reader, &node)
}

// NewNode wraps a Node around a pointer to a Reader
func NewNode(reader *Reader) *Node {
	return &Node{reader}
}
