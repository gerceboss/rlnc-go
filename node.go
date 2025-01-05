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
