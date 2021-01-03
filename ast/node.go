package ast

import (
	"fmt"
	"strings"
)

// ParseNode represents a function to parse ast nodes.
type ParseNode func(p *Parser) (*Node, error)

// Capture is a structure to indicate that the value should be converted to a
// node. If one of the children returns a node, then that node gets returned
type Capture struct {
	// Type of the node.
	Type int
	// Value is the expression to capture the value of the node.
	Value interface{}

	// Convert is an optional functions to change the type of the parsed value.
	// e.g. convert "1" to an integer instead of the string itself.
	Convert func(i string) interface{}
}

// Node is a simple node in a tree with double linked lists instead of slices to
// keep track of its siblings and children. A node is either a value or a
// parent node.
type Node struct {
	// Type of the node.
	Type int
	// Value of the node. Only possible if it has no children.
	Value interface{}

	// Parent is the parent node.
	Parent *Node
	// PreviousSibling is the previous sibling of the node.
	PreviousSibling *Node
	// NextSibling is the next sibling of the node.
	NextSibling *Node
	// FirstChild is the first child of the node.
	FirstChild *Node
	// LastChild is the last child of the node.
	LastChild *Node
}

func (n *Node) String() string {
	if n.IsParent() {
		var children []string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			children = append(children, c.String())
		}
		return fmt.Sprintf("[%03d]: [%v]", n.Type, strings.Join(children, ", "))
	}
	return fmt.Sprintf("[%03d]: %v", n.Type, n.Value)
}

// IsParent returns whether the node has children and thus is not a value node.
func (n *Node) IsParent() bool {
	return n.FirstChild != nil
}

// Children returns all the children of the node.
func (n *Node) Children() (cs []*Node) {
	if n.FirstChild == nil {
		return //  Node has no children.
	}
	cs = append(cs, n.FirstChild)
	for c := n.FirstChild; c != nil; cs = append(cs, c) {
		c = c.NextSibling
	}
	return
}

// SetPrevious inserts the given node as the previous sibling.
func (n *Node) SetPrevious(sibling *Node) {
	sibling.Parent = n.Parent
	if n.PreviousSibling != nil {
		// Already has a sibling.
		// 1. Copy over previous of node.
		// 2. Add node as next.
		// 3. Update next of previous node.
		// 4. Assign sibling as previous.
		sibling.PreviousSibling = n.PreviousSibling // (1)
		sibling.NextSibling = n                     // (2)
		n.PreviousSibling.NextSibling = sibling     // (3)
		n.PreviousSibling = sibling                 // (4)
	}
	// Does not have a previous sibling yet.
	// 1. Reference each other.
	// 2. Update references of parent.
	n.PreviousSibling = sibling // (1)
	sibling.NextSibling = n
	if n.Parent != nil { // (2)
		n.Parent.FirstChild = sibling
	}
	return

}

// SetNext inserts the given node as the next sibling.
func (n *Node) SetNext(sibling *Node) {
	sibling.Parent = n.Parent
	// (a) <-> (b) | a.AddSibling(b)
	// (a) <-> (c) <-> (b)
	if n.NextSibling != nil {
		// Already has a sibling.
		// 1. Copy over next of node.
		// 2. Add node as previous.
		// 3. Update previous of next node.
		// 4. Assign sibling as next.
		sibling.NextSibling = n.NextSibling     // (1)
		sibling.PreviousSibling = n             // (2)
		n.NextSibling.PreviousSibling = sibling // (3)
		n.NextSibling = sibling                 // (4)
		return
	}
	// Does not have a next sibling yet.
	// 1. Reference each other.
	// 2. Update references of parent.
	sibling.PreviousSibling = n // (1)
	n.NextSibling = sibling
	if n.Parent != nil { // (2)
		n.Parent.LastChild = sibling
	}
}

// SetFirst inserts the given node as the first child of the node.
func (n *Node) SetFirst(child *Node) {
	// Set the first child.
	if n.FirstChild != nil {
		n.FirstChild.SetPrevious(child)
		return
	}
	// No children present.
	child.Parent = n
	n.FirstChild = child
	n.LastChild = child
	return
}

// SetLast inserts the given node as the last child of the node.
func (n *Node) SetLast(child *Node) {
	// Set the last child.
	if n.FirstChild != nil {
		n.LastChild.SetNext(child)
		return
	}
	// No children present.
	child.Parent = n
	n.FirstChild = child
	n.LastChild = child
	return
}