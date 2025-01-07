package main

import "testing"
func TestSourceNode(t *testing.T) {
    numChunks:=3
	chunkSize:=4
	committer:=NewCommitter(chunkSize)
	block:=RandomU8Slice(numChunks*chunkSize*32)
	sourceNode,err:=NewSource(*committer,block,numChunks)
	if err!=nil{
		if len(sourceNode.Chunks())!=numChunks{
			t.Errorf("source node test fail")
		}
		if len(sourceNode.Commitments())!=numChunks{
			t.Errorf("source node test fail")
		}
	}
}

func TestSendReceive(t *testing.T){
    numChunks:=3
	chunkSize:=4
	committer:=NewCommitter(chunkSize)
	block:=RandomU8Slice(numChunks*chunkSize*32)
	sourceNode,_:=NewSource(*committer,block,numChunks)
	message,err:=sourceNode.Send()
	if err!=nil{
		t.Errorf("error in sending messages")
	}
	destinationNode:=NewNode(*committer,numChunks)

	destinationNode.Receive(message)
	// if rcvErr!=*LinearlyDependentChunk && rcvErr!=ReceiveError{} {
	// 	t.Errorf("")
	// }// fix this

	
	if len(destinationNode.Chunks())!=1{
		t.Errorf("length not 1")
	}
	if len(destinationNode.Commitments())!=numChunks{
		t.Errorf("length not equal to numChunks")
	}

	destinationNode.Send()
}

func TestDecode(t *testing.T){
	numChunks:=3
	chunkSize:=4
	committer:=NewCommitter(chunkSize)
	block:=RandomU8Slice(numChunks*chunkSize*32)
	sourceNode,_:=NewSource(*committer,block,numChunks)
	message1,err1:=sourceNode.Send()
	message2,err2:=sourceNode.Send()
	message3,err3:=sourceNode.Send()
	if err1!=nil || err2!=nil || err3!=nil{
		t.Errorf("error in sending a message")
	}
	destinationNode:=NewNode(*committer,numChunks)

	//receive has some problem
	destinationNode.Receive(message1)
	destinationNode.Receive(message2)
	destinationNode.Receive(message3)

	// if rcvErr!=*LinearlyDependentChunk && rcvErr!=ReceiveError{} {
	// 	t.Errorf("")
	// }// fix this

	decoded,err:=destinationNode.Decode()
	if err!=nil{
		t.Errorf("error in decoding")
	}
	if len(decoded)!=len(block){
		t.Errorf("decode test fail")
	}
	
	for i:=range decoded{
		if decoded[i]!=block[i]{
			t.Errorf("decode test fail")
		}
	}
}