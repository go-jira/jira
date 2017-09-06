package bugreports

import "github.com/cheekybits/genny/generic"

type Node generic.Type

type DigraphNode struct {
	nodes map[Node][]Node
}

func NewDigraphNode() *DigraphNode {
	return &DigraphNode{
		nodes: make(map[Node][]Node),
	}
}

func (dig *DigraphNode) Add(n Node) {
	if _, exists := dig.nodes[n]; exists {
		return
	}

	dig.nodes[n] = nil
}

func (dig *DigraphNode) Connect(a, b Node) {
	dig.Add(a)
	dig.Add(b)

	dig.nodes[a] = append(dig.nodes[a], b)
}
