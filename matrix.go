package main

import (
	"golang.org/x/crypto/curve25519"
)

type Eschelon struct{
	Coefficients [][]curve25519.Scalar
	Eschelon [][]curve25519.Scalar
	Transform [][]curve25519.Scalar
}

func NewEschelon(size int) *Eschelon{
	transform := make([][]curve25519.Scalar, size) // usinng make ensures that initialisation is with Scalar::Zero value
	for i := range transform {
		transform[i] = make([]curve25519.Scalar, size)
		transform[i][i] = 1 //replace by Scalar::ONE
	}
	return &Eschelon{
		Coefficients :[][]curve25519.Scalar{},
		Eschelon :[][]curve25519.Scalar{},
		Transform:transform,
	}

}
func NewIdentity(size int) *Eschelon{
	eschelon:=make([][]curve25519.Scalar,size)
	for i:=range eschelon{
		eschelon[i]=make([]curve25519.Scalar,size)
		eschelon[i][i]=1
	}
	return &Eschelon{
		Coefficients: eschelon,
		Eschelon: eschelon,
		Transform: eschelon,
	}
}

// returns true if the eschelon form is square.
func (es *Eschelon)IsFull() bool{
	return len(es.Coefficients)==len(es.Coefficients[0])
}
 func (es *Eschelon)AddRow(row []curve25519.Scalar)bool{
	for i:=range row{
		if row[i]!=0 {
			return false
		}// scalar::zero{}
	}
	currentSize:=len(es.Coefficients)
	if currentSize==len(row){
		return false
	}
	if currentSize==0{
		es.Eschelon = append(es.Eschelon, row)
		es.Coefficients = append(es.Coefficients, row)
		return true
	}

	tr:=es.Transform[currentSize]
	i:=0
	newEschelonRow:=row
	for i<currentSize{
		j := firstEntry(es.Eschelon[i])
		k := firstEntry(newEschelonRow)

		if k == -1 { // If no entry exists in the new row, return false
			return false
		}

		if j < k {
			i++
			continue
		}
		if j > k {
			break
		}

		pivot := es.Eschelon[i][j]
		f := newEschelonRow[j]

		for index := range newEschelonRow {
			newEschelonRow[index] = pivot*newEschelonRow[index] - es.Eschelon[i][index]*f
		}
		for index := range tr {
			tr[index] = pivot*tr[index] - es.Transform[i][index]*f
		}
		i++
	}
	for i:=range newEschelonRow{
		if newEschelonRow[i]!=0 {
			return false
		}// scalar::zero{}
	}

	es.Eschelon = append(es.Eschelon[:i], append(newEschelonRow, es.Eschelon[i:]...)...)
	es.Coefficients = append(es.Coefficients, row)

	if i < currentSize {
		es.Transform = append(es.Transform[:currentSize-1], es.Transform[currentSize:]...)
		es.Transform = append(es.Transform[:i], append(tr, es.Transform[i:]...)...)
		return true
	}
	es.Transform[i] = tr
	return true
}

func (es *Eschelon) compoundScalars(scalars []byte) []curve25519.Scalar{
	result := make([]curve25519.Scalar, len(es.Transform))

	// heck if i, j are in wrong position 
	for j := 0; j < len(es.Transform); j++ {
		sum := 0// scalar::zero
		for i, scalar := range scalars {
			sum += scalar* es.Coefficients[i][j] //check this multiplication
		}
		result[j] = sum
	}
	return result
}

// return the index of the first valid entry
func firstEntry(row []curve25519.Scalar) int {
for i, val := range row {
	if val != 0 {
		return i
	}
}
return -1 // Return -1 if no entry is found
}
