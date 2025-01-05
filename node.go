package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

type Message struct{
	chunk Chunk
	commitments []curve25519.Point //ristretto point
}
type Chunk struct{
	data []curve25519.Scalar
	coefficients []curve25519.Scalar
}

type Node struct{
	chunks [][]curve25519.Scalar
	commitments []curve25519.Point //ristretto point
	eschelon Eschelon
	committer Committer
}



type ReceiveError struct {
	Type    string
	Message string
}

func (e *ReceiveError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewReceiveError(errorType, message string) *ReceiveError {
	return &ReceiveError{
		Type:    errorType,
		Message: message,
	}
}

var (
	ExistingCommitmentsMismatch = func(msg string) *ReceiveError {
		return NewReceiveError("ExistingCommitmentsMismatch", msg)
	}
	ExistingChunksMismatch = func(msg string) *ReceiveError {
		return NewReceiveError("ExistingChunksMismatch", msg)
	}
	InvalidMessage = func(msg string) *ReceiveError {
		return NewReceiveError("InvalidMessage", msg)
	}
	LinearlyDependentChunk = NewReceiveError("LinearlyDependentChunk", "The chunk is linearly dependent")
)

// use ristretto
func NewMessage(chunk Chunk,commitments []curve25519.Point)(*Message){
	return &Message{
		chunk:chunk,
		commitments: commitments,
	}
}

func (m *Message) CoefficientsToScalars() []curve25519.Scalar{
	return m.chunk.coefficients
}

func (m *Message) Verify(committer Committer)error{
	 msm:=MSM(m.CoefficientsToScalars(), m.commitments) // finc the necessary package for ristretto
	commitment,err:=committer.Commit(m.chunk.data)
	if err!=nil{
		return err
	}
	if msm!=commitment{
		return errors.New("The commitment does not match")

	}
	return nil
}
func (m *Message) Coefficients() []curve25519.Scalar{
	return m.chunk.coefficients
}


func NewNode(committer Committer,numChunks int) *Node{
	eschelon:=NewEschelon(numChunks)
	return &Node{
		chunks: [][]curve25519.Scalar{},
		commitments: []curve25519.Point{}, //ristretto point
		eschelon : *eschelon, 
		committer :committer,
	}
}

func NewSource(committer Committer,block []byte,numChunks int)(*Node,error){
	chunkies,err:=block_to_chunks(block,numChunks)
	if err!=nil{
		return nil,err
	}
	var chunks []curve25519.Scalar
	for i:=range chunkies{
		it,err:=chunk_to_scalars(chunkies[i])
		if err!=nil{
			return nil,err
		}
		chunks = append(chunks,it)
	}
	var commitments []curve25519.Scalar
	for i:=range chunks{
		res,err:=committer.Commit(chunks[i])
		if err!=nil{
			return nil,err
		}
		commitments = append(commitments, res)
	}
	return &Node{
		chunks: chunks,
		commitments: commitments,
		eschelon: *NewIdentity(numChunks),
		committer: committer,
	},nil
}


func (n *Node) CheckExistingCommitments(commitments []curve25519.Point)error{
	if len(n.commitments)!=0{
		if len(n.commitments)!=len(commitments){
			return errors.New("The number of commitments is different")
		}
		for i:=range commitments{
			if n.commitments[i]!=commitments[i]{
				return errors.New("The commitments donot match")
			}
		}
	}
	return nil
}

func (n *Node) CheckExistingChunks(chunk Chunk) error{
	if len(n.chunks)!=0{
		if len(n.chunks[0])!=len(chunk.data){
			return errors.New("The chunk size iis different")
		}
	}
	return nil
}
// return in a better manner form this function
func (n *Node) Receive(message Message)ReceiveError{
	err:=n.CheckExistingCommitments(message.commitments)
	if err!=nil{
		return *ExistingChunksMismatch(err.Error())
	}

	err2:= n.CheckExistingChunks(message.chunk)
	if err2!=nil{
		return *ExistingChunksMismatch(err2.Error())
	}

	err3:=message.Verify(n.committer)
	if err3!=nil{
		return *InvalidMessage(err3.Error())
	}

	//Verify linear independence
	if !n.eschelon.AddRow(message.chunk.coefficients){
		return *LinearlyDependentChunk
	}

	n.chunks = append(n.chunks, message.chunk.data)

	if len(n.commitments)==0{
		n.commitments=message.commitments
	}
	return ReceiveError{}
}

func (n *Node) Send() (Message,error){
	if len(n.chunks)==0{
		return Message{},errors.New("There are no chunks to send")
	}
	scalars:=// generate random coefficeints
	chunk:=n.LinearCombChunk(scalars)
	message:=NewMessage(chunk,n.commitments)
	err:=message.Verify(n.committer)
	if err!=nil{
		return Message{},err
	}
	return message,nil
}

func (n *Node) LinearCombChunk(scalars []byte)Chunk{
	coefficients:=n.eschelon.compoundScalars(scalars)
	data:=n.LinearCombData(scalars)
	return Chunk{
		coefficients: coefficients,
		data:data,
	}

}
func (n *Node) LinearCombData(scalars []byte)[]curve25519.Scalar{

	result := make([]curve25519.Scalar, len(n.chunks[0]))
	for i := 0; i < len(n.chunks[0]); i++ {
		sum := 0
		for j, scalar := range scalars {
			if j >= len(n.chunks) {
				return nil, errors.New("Mismatch between scalars and chunks")
			}
			sum += scalar * n.chunks[j][i] //Scalar multiplication
		}
		result[i] = sum
	}
	return result
}
