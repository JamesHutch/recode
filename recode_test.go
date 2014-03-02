// Copyright 2014 James Hutchinson. All rights reserved.
//
// Use of this source code is governed by the Apache v2
// license that can be found in the LICENSE file.

// Simple module to cast fields in one struct into another

package recode

import (
	"testing"
)

type InputStruct struct {
	OutputString interface{}
	OutputInt    interface{}
	OutputFloat  interface{}
	OutputUint   interface{}
}

type OutputStruct struct {
	OutputString string
	OutputInt    int
	OutputFloat  float64
	OutputUint   uint
}

type RecodeTest struct {
	InputData  *InputStruct
	ExpectData *OutputStruct
}

type EmbedInputStruct struct {
	Data InputStruct
}

type EmdedOutputStruct struct {
	Data OutputStruct
}

var recodeTests = []RecodeTest{
	{&InputStruct{"string", uint(1), 1.1, 1}, &OutputStruct{"string", 1, 1.1, 1}},
	{&InputStruct{"1.5", "1", "1.1", "1"}, &OutputStruct{"1.5", 1, 1.1, 1}},
	{&InputStruct{"1.5", "-10", -1, 1.0}, &OutputStruct{"1.5", -10, -1, 1}},
	{&InputStruct{"1.5", "-10", uint(1), 1.0}, &OutputStruct{"1.5", -10, 1, 1}},
	{&InputStruct{1.5, "-10", uint(1), 1.0}, &OutputStruct{"1.5", -10, 1, 1}},
	{&InputStruct{1, "-10", uint(1), 1.0}, &OutputStruct{"1", -10, 1, 1}},
}

func TestRecode(t *testing.T) {
	for _, test := range recodeTests {
		var out OutputStruct
		err := Recode(test.InputData, &out)
		if err != nil {
			t.Fatalf("Recode: unexpected error, %v", err)
		}
		if out.OutputInt != test.ExpectData.OutputInt || out.OutputFloat != test.ExpectData.OutputFloat || out.OutputString != test.ExpectData.OutputString || out.OutputUint != test.ExpectData.OutputUint {
			t.Errorf("Recode: got %d %f \"%s\" %d, expected %d %f \"%s\" %d", out.OutputInt, out.OutputFloat, out.OutputString, out.OutputUint, test.ExpectData.OutputInt, test.ExpectData.OutputFloat, test.ExpectData.OutputString, test.ExpectData.OutputUint)
		}
		iEmb := EmbedInputStruct{*test.InputData}
		var oEmb EmdedOutputStruct
		err = Recode(iEmb, &oEmb)
		if err != nil {
			t.Fatalf("Recode: unexpected error, %v", err)
		}
	}
}
