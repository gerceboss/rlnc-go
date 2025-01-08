package components

import (
	"errors"
	"fmt"

	"github.com/bwesterb/go-ristretto"
)

type Echelon struct{
	Coefficients [][]ristretto.Scalar
	Echelon [][]ristretto.Scalar
	Transform [][]ristretto.Scalar
}

func NewEschelon(size int) *Echelon{
	transform := make([][]ristretto.Scalar, size)
	for i := range transform {
		transform[i] = make([]ristretto.Scalar, size)
		transform[i][i].SetOne() 
	}
	return &Echelon{
		Coefficients :[][]ristretto.Scalar{},
		Echelon :[][]ristretto.Scalar{},
		Transform:transform,
	}

}
func NewIdentity(size int) *Echelon{
	echelon:=make([][]ristretto.Scalar,size)
	for i:=range echelon{
		echelon[i]=make([]ristretto.Scalar,size)
		echelon[i][i].SetOne()
	}
	return &Echelon{
		Coefficients: echelon,
		Echelon: echelon,
		Transform: echelon,
	}
}

// returns true if the echelon form is square.
func (es *Echelon)IsFull() bool{
	return len(es.Coefficients)==len(es.Coefficients[0])
}
 func (es *Echelon)AddRow(row []ristretto.Scalar)bool{
	var chk int
	// var z ristretto.Scalar
	for i:=range row{
		if row[i].IsNonZeroI()==1 {
			chk++
		}
	}
	if chk==0{
		return false
	}
	currentSize:=len(es.Coefficients)
	if currentSize==len(row){
		return false
	}
	if currentSize==0{
		es.Echelon = append(es.Echelon, row)
		es.Coefficients = append(es.Coefficients, row)
		return true
	}

	tr:=es.Transform[currentSize]
	i:=0
	newEschelonRow:=row
	for i<currentSize{
		j := firstEntry(es.Echelon[i])
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

		pivot := es.Echelon[i][j]
		f := newEschelonRow[j]

		for index := range newEschelonRow {
			var r1,r2 ristretto.Scalar
			r1.Mul(&pivot,&newEschelonRow[index]) 
			r2.Mul(&es.Echelon[i][index],&f)
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

	chk=0
	for i:=range newEschelonRow{
		if newEschelonRow[i].IsNonZeroI()==1 {
			chk++
		}
	}
	if chk==0{
		return false
	}

	newEschelonRowSlice:=[][]ristretto.Scalar{newEschelonRow}
	es.Echelon = append(es.Echelon[:i], append(newEschelonRowSlice, es.Echelon[i:]...)...)
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

func (es *Echelon) CompoundScalars(scalars []byte) []ristretto.Scalar{
	result := make([]ristretto.Scalar, len(es.Transform))

	// check if i, j are in wrong position 
	for j := 0; j < len(es.Transform); j++ {
		var sum ristretto.Scalar
		sum.SetZero()
		for i:=0;i<len(scalars);i++ {
			var s ristretto.Scalar
			s.SetZero()
			scalarSlice:=[]byte{scalars[i]}
			if i<len(es.Coefficients){
				var temp1,temp2 ristretto.Scalar
				temp1.Set(&sum)
				temp2.Set(s.Derive(scalarSlice))
				sum.MulAdd(&temp2, &es.Coefficients[i][j],&temp1) //check this multiplication
			}
		}
		fmt.Println(&sum)
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
func (es *Echelon)Inverse()([][]ristretto.Scalar,error){
	if len(es.Coefficients)==0{
		return nil,errors.New("no coefficients to decode")
	}
	if len(es.Echelon)!=len(es.Coefficients[0]){
		return nil,errors.New("the echelon form is not square")
	}
	inverse := make([][]ristretto.Scalar, len(es.Transform))
	for i := range es.Transform {
		inverse[i] = make([]ristretto.Scalar, len(es.Transform[i]))
		copy(inverse[i], es.Transform[i])
	}

	for i := len(es.Echelon) - 1; i >= 0; i-- {
		var pivot ristretto.Scalar
		pivot.Inverse(&es.Echelon[i][i])
		for k := range inverse[i] {
			inverse[i][k].Mul(&pivot,&inverse[i][k])
		}
		for j := i + 1; j < len(es.Echelon); j++ {
			var diff ristretto.Scalar
			diff.Mul(&es.Echelon[i][j],&pivot)
			for k := range es.Echelon {
				var actualDiff ristretto.Scalar
				actualDiff.Mul(&inverse[j][k],&diff)
				inverse[i][k].Sub(&inverse[i][k],&actualDiff)
			}
		}
	}

	return inverse, nil
}
