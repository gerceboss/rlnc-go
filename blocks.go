package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/bwesterb/go-ristretto"
)

type Committer struct{
	Generators []ristretto.Point
}

func NewCommitter(n int) *Committer{
	return &Committer{
		Generators: generators(n),
	}
}

func (c *Committer ) Len() int{
	return len(c.Generators)
}

// Commit computes the commitment for the given scalars.
func (c *Committer) Commit(scalars []ristretto.Scalar) (ristretto.Point, error) {
	if len(scalars) > len(c.Generators) {
		return ristretto.Point{}, errors.New("chunk size is too large")
	}

	// Multiscalar multiplication
	result := curve25519.NewIdentityPoint() //fix this
	for i, scalar := range scalars {
		term := curve25519.NewIdentityPoint()
		term.Mul(&c.Generators[i], &scalar)
		result.Add(result, term)
	}

	return *result, nil
}


func chunk_to_scalars(chunk []byte) ([]ristretto.Scalar,error){

	if len(chunk)%32!=0{
		return nil,errors.New("Chunk size is not divisible by 32")
	}

	scalars := []ristretto.Scalar{}
	for i := 0; i < len(chunk); i += 32 {
		var scalar ristretto.Scalar
		var temp [32]byte
		copy(temp[:], chunk[i:i+32])
		scalar.SetBytes(&temp)
		scalars = append(scalars, scalar)
	}

	return scalars, nil
}

func block_to_chunks(block []byte,num_chunks int)([][]byte,error){
	if len(block)%num_chunks!=0{
		return nil,errors.New("Block size is not divisible by num_chunks")
	}
	chunkSize:=len(block)/num_chunks
	var chunks [][]byte
	for i := 0; i < len(block); i += chunkSize {
		chunks = append(chunks, block[i:i+chunkSize])
	}
	return chunks, nil
}

func RandomU8Slice(len int)[]uint8{
	rand.Seed(time.Now().UnixNano())
	ret := make([]uint8, len)
	for i := 0; i < len; i++ {
		ret[i] = uint8(rand.Intn(256)) 
	}
	for i:=31;i<len;i+=32{
		ret[i]=0
	}
	return ret
}

func generators(n int)[]ristretto.Point{
	result:=make([]ristretto.Point,n)
	for i:=0;i<n;i++{
		var c ristretto.Point
		var r ristretto.Scalar
		r.Rand()
		result[i]= *c.PublicScalarMultBase(&r) //scalar multiplication with the ristretto_basepoint
	}
	return result
}
