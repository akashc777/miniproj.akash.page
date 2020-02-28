package main

import (
	"fmt"
	//"io/ioutil"
	"log"
	"math/rand"
	//"os"
	"miniproj.akash.page/node"
	"time"
)

func main()  {
  rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ldate | log.Lmicroseconds)

  nodes, leader_node := startCluster(3)

	for{
		time.Sleep(250 * time.Millisecond)
		if countLeaders(nodes) == 1 {
			fmt.Println("leaders should still be one who is %v", leader_node.ID)

		}
	}

}

func gimmeNodes(num int) []*node.Node {
	var nodes []*node.Node

	for i := 0; i < num; i++ {
		transport := &node.HTTPTransport{Address: "127.0.0.1:0"}
		logger := &node.Log{}
		applyer := &node.StateMachine{}
		node := node.NewNode(fmt.Sprintf("%d", i), transport, logger, applyer)
		nodes = append(nodes, node)
		nodes[i].Serve()
	}

	// let them start serving
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if j != i {
				nodes[i].AddToCluster(nodes[j].Transport.String())
			}
		}
	}

	for _, node := range nodes {
		node.Start()
	}

	return nodes
}




func countLeaders(nodes []*node.Node) int {
	leaders := 0
	for i := 0; i < len(nodes); i++ {
		nodes[i].RLock()
		if nodes[i].State == node.Leader {
			leaders++
		}
		nodes[i].RUnlock()
	}
	return leaders
}



func findLeader(nodes []*node.Node) *node.Node {
	for i := 0; i < len(nodes); i++ {
		nodes[i].RLock()
		if nodes[i].State == node.Leader {
			nodes[i].RUnlock()
			return nodes[i]
		}
		nodes[i].RUnlock()
	}
	return nil
}

func startCluster(num int) ([]*node.Node, *node.Node) {
	nodes := gimmeNodes(num)
	for {
		time.Sleep(50 * time.Millisecond)
		if countLeaders(nodes) == 1 {
			break
		}
	}
	leader := findLeader(nodes)
	return nodes, leader
}

func stopCluster(nodes []*node.Node) {
	for _, node := range nodes {
		node.Exit()
	}
}
