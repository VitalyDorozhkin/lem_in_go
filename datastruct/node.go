package datastruct

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Link struct {
	NodeStart *Node
	NodeEnd   *Node
}

type Step struct {
	LeminNumber int
	NodeEnd     *Node
}

type Graph struct {
	Nodes map[string]*Node
	Links []*Link
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func (g *Graph) Link(root string, linked ...string) *Graph {
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
	Name   string
	Nodes  Nodes
	Status string
	X      int
	Y      int
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

func NewNode(name string) *Node {
	return &Node{
		Name: name,
	}
}

func NewDefaultGraph() *Graph {
	graph := NewGraph()
	graph.
		Link("0", "1", "2", "3").
		Link("1", "4").
		Link("2", "5").
		Link("3", "6").
		Link("4", "7").
		Link("5", "7").
		Link("6", "7").
		Link("7", "0")

	graph.Nodes["0"].X, graph.Nodes["0"].Y = 60, 500
	graph.Nodes["1"].X, graph.Nodes["1"].Y = 30, 240
	graph.Nodes["2"].X, graph.Nodes["2"].Y = 120, 500
	graph.Nodes["3"].X, graph.Nodes["3"].Y = 120, 680
	graph.Nodes["4"].X, graph.Nodes["4"].Y = 240, 240
	graph.Nodes["5"].X, graph.Nodes["5"].Y = 240, 500
	graph.Nodes["6"].X, graph.Nodes["6"].Y = 240, 680
	graph.Nodes["7"].X, graph.Nodes["7"].Y = 400, 760
	return graph
}

func NewReadedGraph() (ants []*Node, graph*Graph, stepsList[][]Step) {
	graph = NewGraph()
	scanner := bufio.NewScanner(os.Stdin)
	prev := ""
	var count int
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("|%s|\n", line)
		if len(line) == 0 {
			break
		}
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "##") {
				prev = strings.TrimPrefix(line, "##")
			}
			continue
		}
		if newCount, err := scanCount(line); err == nil {
			count = newCount
			println("count:", count)
		} else if node := scanNode(line); node != nil {
			node.Status = prev
			graph.Nodes[node.Name] = node
			println("node:", node.Name)
			if node.Status == "start" {
				println("sccscs", count)
				ants = make([]*Node, count, count)
				for i := range ants {
					ants[i] = node
				}
			}
		} else if from, to := scanLink(line); from != "" && to != "" {
			println(from, "-", to)
			graph.Link(from, to)
		}
		prev = ""
	}
	stepsList = make([][]Step, 0, 0)
	for scanner.Scan(){
		line := scanner.Text()
		fmt.Printf("|%s|\n", line)
		if len(line) == 0 {
			break
		}
		dos := strings.Split(line, " ")
		steps := make([]Step, 0, len(dos))
		for _, do := range dos {
			res := strings.Split(do, "-")
			fmt.Printf("%+v\n", res)
			if len(res) != 2 || !strings.HasPrefix(res[0], "L") || !isInt(strings.TrimPrefix(res[0], "L")) {
				return nil, nil, nil
			}
			idx, _ := strconv.Atoi(strings.TrimPrefix(res[0], "L"))
			println(idx)
			steps = append(steps, Step{
				LeminNumber: idx,
				NodeEnd: graph.Nodes[res[1]],
			})
		}
		stepsList = append(stepsList, steps)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func scanCount(line string) (int, error) {
	if !isInt(line) {
		return 0, errors.New("not int")
	}
	return strconv.Atoi(line)
}

func scanNode(line string) (node *Node) {
	if res := strings.Split(line, " "); len(res) == 3 && isInt(res[1]) && isInt(res[2]) {
		node = &Node{
			Name: res[0],
		}
		node.X, _ = strconv.Atoi(res[1])
		node.Y, _ = strconv.Atoi(res[2])
	}
	return
}

func scanLink(line string) (from string, to string) {
	if res := strings.Split(line, "-"); len(res) == 2 {
		from = res[0]
		to = res[1]
	}
	return
}

func ifMin(a *int, b int) {
	if b < *a {
		*a = b
	}
}

func ifMax(a *int, b int) {
	if b > *a {
		*a = b
	}
}
func MoveGraph(g *Graph, side int, padding int) {
	minX, minY, max := 10000000, 10000000, -1
	for _, v := range g.Nodes {
		ifMin(&minY, v.Y)
		ifMin(&minX, v.X)
	}
	for _, v := range g.Nodes {
		v.Y -= minY
		v.X -= minX
		ifMax(&max, v.Y)
		ifMax(&max, v.X)
	}
	c := float64(side-2*padding) / float64(max)
	for _, v := range g.Nodes {
		v.Y = int(float64(v.Y)*c) + padding
		v.X = int(float64(v.X)*c) + padding
	}
}
