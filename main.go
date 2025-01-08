package main

import (
	"fmt"
	"rlnc/components"
)

func main() {
    numNodes:=10000
    chunkSize:=1
    committer:=components.NewCommitter(chunkSize)
    meshSize:=10
    network:=components.NewNetwork(*committer,numNodes,meshSize)
    for i:=0; network.Timestamp<100 && !network.AllNodesFull();i++{
        network.Round()
        fmt.Println("Timestamp: ",network.Timestamp,"Full nodes: ",network.FullNodes,"Wasted Bandwidth: ",network.WastedBandwidth)
    }
}
