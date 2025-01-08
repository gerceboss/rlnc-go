package components

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type SimulationNode struct{
    node Node
    neighbours []int //usize
    sentMessage bool
}

type Network struct{
    nodes []SimulationNode
    Timestamp uint32
    WastedBandwidth uint32
    FullNodes int//usize
    roundMessages []Message
    roundDestinations []int //usize
}

func NewSimulationNode(committer Committer,numChunks int) *SimulationNode{
    return &SimulationNode{
        node: *NewNode(committer,numChunks),
        neighbours: []int{},
        sentMessage: false,
    }
}

func NewSourceSim(committer Committer,block []byte,numChunks int) (SimulationNode,error){
    node,err:=NewSource(committer,block,numChunks)
    if err!=nil{
        return SimulationNode{},errors.New(err.Error())
    }
    return SimulationNode{
        node: *node,
        neighbours: []int{},
        sentMessage: false,
    },nil
}


func CreateNodes(committer Committer,num int,numChunks int,meshSize int,block []byte)([]SimulationNode,error){
    ret:=make([]SimulationNode,num)
    sourceNode,err:=NewSourceSim(committer,block,numChunks)
    if err!=nil{
        return []SimulationNode{},err
    }
    ret = append(ret, sourceNode)
    for i:=1;i<num;i++{
        ret = append(ret, *NewSimulationNode(committer,numChunks))
    }

    for i:=0;i<num;i++{
        neighbours:=make([]int,meshSize)
        rand.Seed(time.Now().UnixNano()) 
        for j:=0;j<meshSize;j++{
            neighbours = append(neighbours,rand.Intn(256)%num )
        }
        ret[i].neighbours=neighbours
    }
    return ret,nil
}

func NewNetwork(committer Committer,numNodes int,meshSize int) *Network{
    numChunks:=10
    nodes,_:=CreateNodes(committer,numNodes,numChunks,meshSize,RandomU8Slice(committer.Len()*numChunks*32))
    return &Network{
        nodes: nodes,
        Timestamp: 0,
        WastedBandwidth: 0,
        FullNodes: 1,
        roundDestinations:[]int{} ,
        roundMessages: []Message{},
    }
}

func (net *Network)Round(){
    net.Timestamp+=1
    net.roundMessages=[]Message{}
    net.roundDestinations=[]int{}
    for i:=0;i<len(net.nodes);i++{
        source:=net.nodes[i]
        for j:=range source.neighbours{
            if j==i{
                continue
            }
            msg,err:=source.node.Send()
            if err==nil{
                source.sentMessage = true
                net.roundMessages = append(net.roundMessages, msg)
                net.roundDestinations = append(net.roundDestinations, j)
            }
        }
    }

    for i,message:=range net.roundMessages{
        j:=net.roundDestinations[i]
        destination:=net.nodes[j]
        err:=destination.node.Receive(message)
        chk:= ReceiveError{
            Type: "",
            Message: "",
        }
        switch err{
            case chk:
                if destination.node.IsFull(){
                    net.FullNodes+=1
                }
            case *LinearlyDependentChunk:
                    net.WastedBandwidth+=1
            default:
                panic(fmt.Sprintf("Unhandled error: %v", err))
        }
    }
}

func (net *Network) AllNodesFull()bool{
    return net.FullNodes==len(net.nodes)
}

