package datastruct

type Link struct {
	NodeStart *Node
	NodeEnd   *Node
}

type Graph struct {
	Nodes map[int]*Node
	Links []*Link
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[int]*Node),
	}
}

func (g *Graph) Link(root int, linked ...int) *Graph {
	rootNode, ok := g.Nodes[root]
	if !ok {
		rootNode = NewNode(root)
		g.Nodes[root] = rootNode
	}

	linkedNodes := make(Nodes, 0, len(linked))
	for _, v := range linked {
		node, ok := g.Nodes[v]
		if !ok {
			node = NewNode(v)
			g.Nodes[v] = node
		}
		linkedNodes = append(linkedNodes, node)
	}

	g.Links = append(g.Links, rootNode.Link(linkedNodes...)...)
	return g
}

type Nodes []*Node

type Node struct {
	Number int
	Nodes  Nodes
	Status string
}

func (n *Node) Link(nodes ...*Node) (links []*Link) {
	links = make([]*Link, 0, len(nodes))
	n.Nodes = append(n.Nodes, nodes...)
	for _, node := range nodes {
		links = append(links, &Link{
			NodeStart: n,
			NodeEnd:   node,
		})
	}
	return links
}

func NewNode(number int) *Node {
	return &Node{
		Number: number,
	}
}
