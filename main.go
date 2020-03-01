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

var currentnode int = 0

func main()  {
  rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ldate | log.Lmicroseconds)


	var e make(chan)

	var prip = []string{"192.168.99.104:59964", "192.168.99.105:59964", "192.168.99.106:59964"}

	transport := &node.HTTPTransport{Address: prip[currentnode]}
	logger := &node.Log{}
	node := node.NewNode(fmt.Sprintf("%d", currentnode), transport, logger)

	node.Serve()

	// let node start serving
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < len(prip); i++ {

		if currentnode != i {
			node.AddToCluster(prip[i])
		}

	}

	node[currentnode].Start()

	select{
		case <-e:
			stopNode(node)

	}

}


func stopNode(node *node.Node) {

	node.Exit()
}
