package main

import (
	"errors"

	"github.com/bwesterb/go-ristretto"
)

type Eschelon struct{
	Coefficients [][]ristretto.Scalar
	Eschelon [][]ristretto.Scalar
	Transform [][]ristretto.Scalar
}

func NewEschelon(size int) *Eschelon{
	transform := make([][]ristretto.Scalar, size)
	for i := range transform {
		transform[i] = make([]ristretto.Scalar, size)
		transform[i][i].SetOne() 
	}
	return &Eschelon{
		Coefficients :[][]ristretto.Scalar{},
		Eschelon :[][]ristretto.Scalar{},
		Transform:transform,
	}

}
func NewIdentity(size int) *Eschelon{
	eschelon:=make([][]ristretto.Scalar,size)
	for i:=range eschelon{
		eschelon[i]=make([]ristretto.Scalar,size)
		eschelon[i][i].SetOne()
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
 func (es *Eschelon)AddRow(row []ristretto.Scalar)bool{
	for i:=range row{
		if row[i].IsNonZeroI()==1 {
			return false
		}
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
			var r1,r2 ristretto.Scalar
			r1.Mul(&pivot,&newEschelonRow[index]) 
			r2.Mul(&es.Eschelon[i][index],&f)
			newEschelonRow[index].Sub(&r1,&r2)
		}
		for index := range tr {
			var r1,r2 ristretto.Scalar
			r1.Mul(&pivot,&tr[index] ) 
			r2.Mul(&es.Transform[i][index],&f)
			tr[index].Sub(&r1,&r2)
		}
		i++
	}
	for i:=range newEschelonRow{
		if newEschelonRow[i].IsNonZeroI()==1 {
			return false
		}
	}
	newEschelonRowSlice:=[][]ristretto.Scalar{newEschelonRow}
	es.Eschelon = append(es.Eschelon[:i], append(newEschelonRowSlice, es.Eschelon[i:]...)...)
	es.Coefficients = append(es.Coefficients, row)

	if i < currentSize {
		trSlice:=[][]ristretto.Scalar{tr}
		es.Transform = append(es.Transform[:currentSize-1], es.Transform[currentSize:]...)
		es.Transform = append(es.Transform[:i], append(trSlice, es.Transform[i:]...)...)
		return true
	}
	es.Transform[i] = tr
	return true
}

func (es *Eschelon) compoundScalars(scalars []byte) []ristretto.Scalar{
	result := make([]ristretto.Scalar, len(es.Transform))

	// check if i, j are in wrong position 
	for j := 0; j < len(es.Transform); j++ {
		var sum ristretto.Scalar
		sum.SetZero()
		for i, scalar := range scalars {
			var s ristretto.Scalar
			scalarSlice:=[]byte{scalar}
			sum.MulAdd(s.Derive(scalarSlice), &es.Coefficients[i][j],&sum) //check this multiplication
		}
		result[j] = sum
	}
	return result
}

// return the index of the first valid entry
func firstEntry(row []ristretto.Scalar) int {
for i, val := range row {
	if val.IsNonZeroI()==1 {
		return i
	}
}
return -1 // Return -1 if no entry is found
}
func (es *Eschelon)inverse()([][]ristretto.Scalar,error){
	if len(es.Coefficients)==0{
		return nil,errors.New("no coefficients to decode")
	}
	if len(es.Eschelon)!=len(es.Coefficients[0]){
		return nil,errors.New("the eschelon form is not square")
	}
	inverse := make([][]ristretto.Scalar, len(es.Transform))
	for i := range es.Transform {
		inverse[i] = make([]ristretto.Scalar, len(es.Transform[i]))
		copy(inverse[i], es.Transform[i])
	}

	for i := len(es.Eschelon) - 1; i >= 0; i-- {
		var pivot ristretto.Scalar
		pivot.Inverse(&es.Eschelon[i][i])
		for k := range inverse[i] {
			inverse[i][k].Mul(&pivot,&inverse[i][k])
		}
		for j := i + 1; j < len(es.Eschelon); j++ {
			var diff ristretto.Scalar
			diff.Mul(&es.Eschelon[i][j],&pivot)
			for k := range es.Eschelon {
				var actualDiff ristretto.Scalar
				actualDiff.Mul(&inverse[j][k],&diff)
				inverse[i][k].Sub(&inverse[i][k],&actualDiff)
			}
		}
	}

	return inverse, nil
}
