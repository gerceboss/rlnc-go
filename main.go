package main

import (
	"fmt"
    "rlnc_go/components"
)

func main() {
    numNodes:=10000
    chunkSize:=1
    committer:=NewCommitter(chunkSize)
    meshSize:=10
    network:=NewNetwork(*committer,numNodes,meshSize)
    for i:=0; network.timestamp<100 && !network.AllNodesFull();i++{
        network.Round()
        fmt.Println("Timestamp: ",network.timestamp,"Full nodes: ",network.fullNodes,"Wasted Bandwidth: ",network.wastedBandwidth)
    }
}
