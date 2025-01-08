package main

import (
	"testing"

	"github.com/bwesterb/go-ristretto"
)
func TestAddRow(t *testing.T){
	echelon:=NewEschelon(3)
	var s1,s2,s3 ristretto.Scalar
	row1:=[]ristretto.Scalar{*s1.SetZero(),*s2.SetZero(),*s3.SetZero()}
	if echelon.AddRow(row1){
		t.Errorf("AddRow failed")
	}
	row2:=[]ristretto.Scalar{*s1.SetUint64(0),*s2.SetUint64(0),*s3.SetUint64(1)}
	if !echelon.AddRow(row2){
		t.Errorf("AddRow failed")
	}
	// add more tests
}

func TestInverse(t *testing.T){
	echelon:=NewEschelon(3)
	_,err:=echelon.Inverse()
	if err==nil{
		t.Errorf("No error detected on emty inverse")
	}
	var s1,s2,s3 ristretto.Scalar
	row1:=[]ristretto.Scalar{*s1.SetUint64(1),*s2.SetUint64(0),*s3.SetUint64(0)}
	row2:=[]ristretto.Scalar{*s1.SetUint64(0),*s2.SetUint64(1),*s3.SetUint64(0)}
	row3:=[]ristretto.Scalar{*s1.SetUint64(0),*s2.SetUint64(0),*s3.SetUint64(1)}
	echelon.AddRow(row1)
	echelon.AddRow(row2)
	echelon.AddRow(row3)
	inverse,err:=echelon.Inverse()
	if err!=nil{
		t.Errorf("error in inverting an identity matrix")
	}
	if inverse[0][0]!=*s1.SetUint64(1) || inverse[0][1]!=*s1.SetUint64(0){
		t.Errorf("inverse is calculated incorrectly")
	}

	// add one more exampple to invert
}

// fails at 2nd, either bug in AddRow or CompoundScalars
func TestCompoundScalars(t *testing.T){
	echelon:=NewEschelon(3)
	scalars:=echelon.CompoundScalars([]byte{1,2,3})
	for i:=range scalars{
		if scalars[i].IsNonZeroI()==1{
			t.Errorf("compound scalars fail 1")
		}
	}
	
	echelon=NewEschelon(3)
	var s1,s2,s3 ristretto.Scalar
	row:=[]ristretto.Scalar{*s1.SetUint64(6),*s2.SetUint64(15),*s3.SetUint64(5)}
	row2:=[]ristretto.Scalar{*s1.SetUint64(2),*s2.SetUint64(0),*s3.SetUint64(0)}
	echelon.AddRow(row2)
	row3:=[]ristretto.Scalar{*s1.SetUint64(0),*s2.SetUint64(3),*s3.SetUint64(1)}
	echelon.AddRow(row3)

	scalars=echelon.CompoundScalars([]byte{3,5})
	for i:=range scalars{
		// fmt.Println(scalars[i].Bytes(),"!=",row[i].Bytes())
		if scalars[i].EqualsI(&row[i])==0{
			t.Errorf("coumpound scalars fail 2")
		}
	}
}