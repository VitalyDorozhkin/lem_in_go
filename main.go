package main

import (
	"bufio"
	"fmt"
	. "github.com/VitalyDorozhkin/lem_in_go/datastruct"
	"log"
	"math"
	"os"
)

/*


	_-/	1	-	4	\-_
0	---	2	-	5	---	7
|	-_\	3	-	6	/_-	|
\-----------------------/

*/

func main() {
	graph := NewGraph()


	return
	graph.
		Link(0, 1, 2, 3).
		Link(1, 4).
		Link(2, 5).
		Link(3, 6).
		Link(4, 7).
		Link(5, 7).
		Link(6, 7).
		Link(7, 0)

	graph.Nodes[0].X, graph.Nodes[0].Y = 15, 127
	graph.Nodes[1].X, graph.Nodes[1].Y = 31, 63
	graph.Nodes[2].X, graph.Nodes[2].Y = 31, 127
	graph.Nodes[3].X, graph.Nodes[3].Y = 31, 191
	graph.Nodes[4].X, graph.Nodes[4].Y = 63, 63
	graph.Nodes[5].X, graph.Nodes[5].Y = 63, 127
	graph.Nodes[6].X, graph.Nodes[6].Y = 63, 191
	graph.Nodes[7].X, graph.Nodes[7].Y = 95, 239

	println("Simple Ford Bellman")
	solve(*graph, 0)
	solve(*graph, 1)
	solve(*graph, 2)
	solve(*graph, 3)
	solve(*graph, 4)
	solve(*graph, 5)
	solve(*graph, 6)
	solve(*graph, 7)

	println("Ford Bellman working while changing")
	solveWhileChanging(*graph, 0)
	solveWhileChanging(*graph, 1)
	solveWhileChanging(*graph, 2)
	solveWhileChanging(*graph, 3)
	solveWhileChanging(*graph, 4)
	solveWhileChanging(*graph, 5)
	solveWhileChanging(*graph, 6)
	solveWhileChanging(*graph, 7)

	println("Ford Bellman on Queue(channels)")
	solveWithQueue(*graph, 0)
	solveWithQueue(*graph, 1)
	solveWithQueue(*graph, 2)
	solveWithQueue(*graph, 3)
	solveWithQueue(*graph, 4)
	solveWithQueue(*graph, 5)
	solveWithQueue(*graph, 6)
	solveWithQueue(*graph, 7)
}

func solve(graph Graph, startPos int) {
	dist := make([]int64, len(graph.Nodes))

	for i := range dist {
		dist[i] = math.MaxInt64 - 1
	}
	dist[startPos] = 0
	c := 0
	for i := 0; i < len(graph.Nodes)-1; i++ {
		for _, r := range graph.Links {
			c++
			if dist[r.NodeEnd.Number] > dist[r.NodeStart.Number]+1 {
				dist[r.NodeEnd.Number] = dist[r.NodeStart.Number] + 1
			}
		}
	}
	fmt.Printf("%d%v\n\n",c, dist)
}

func solveWhileChanging(graph Graph, startPos int) {
	dist := make([]int64, len(graph.Nodes))

	for i := range dist {
		dist[i] = math.MaxInt64 - 1
	}
	dist[startPos] = 0
	c := 0
	for {
		any := false
		for _, r := range graph.Links {
			c++
			if dist[r.NodeEnd.Number] > dist[r.NodeStart.Number]+1 {
				dist[r.NodeEnd.Number] = dist[r.NodeStart.Number] + 1
				any = true
			}
		}
		if !any {
			break
		}
	}
	fmt.Printf("%d%v\n\n",c, dist)
}

func solveWithQueue(graph Graph, startPos int) {
	dist := make([]int64, len(graph.Nodes))
	queue := make(chan *Node, len(graph.Nodes))
	queue <- graph.Nodes[startPos]
	for i := range dist {
		dist[i] = math.MaxInt64 - 1
	}
	dist[startPos] = 0
	c := 0
	for len(queue) > 0 {
		node := <-queue
		for _, r := range node.Nodes {
			c++
			if dist[r.Number] > dist[node.Number]+1 {
				dist[r.Number] = dist[node.Number] + 1
				queue <- r
			}
		}
	}
	close(queue)
	fmt.Printf("%d%v\n\n", c, dist)
}
