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
// Return a probably not-useful NodeName destructuring struct (say that 3 times fast)
// Additionally, return the name and any errors
// SATISFY ALL OF THE INTERFACE
func (n *Node) Name() (name string, err error) {
	var nodeName NativeNode
	err = mapstructure.Decode(n.Reader, &nodeName)
	name = nodeName.Name
	return
}

// NewNode wraps a Node around a pointer to a Reader
func NewNode(reader *Reader) *Node {
	return &Node{reader}
}
