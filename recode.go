// Copyright 2014 James Hutchinson. All rights reserved.
//
// Use of this source code is governed by the Apache v2
// license that can be found in the LICENSE file.

// Simple module to cast fields in one struct into another

package recode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Recode(in, out interface{}) error {
	// Check the input is a struct with arbitrary amount of indirection
	inVal := reflect.ValueOf(in)
	for inVal.Kind() == reflect.Ptr {
		if inVal.IsNil() {
			return errors.New("recode: Recode() expects a non-nil value for in")
		}
		inVal = inVal.Elem()
	}

	if inVal.Kind() != reflect.Struct {
		return fmt.Errorf("recode: Recode() expects a struct for in. Got %v", inVal.Kind())
	}
	// Output should be a struct with one level of indirection
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("recode: Recode() expects a pointer for out. Got %v", outVal.Kind())
	}
	if outVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("recode: Recode() expects a pointer to a struct for out. Got pointer to %v", outVal.Elem().Kind())
	}
	return recodeStruct(inVal, outVal.Elem())
}

func recodeStruct(in, out reflect.Value) error {
	t := out.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" { // unexported
			continue
		}
		inF := in.FieldByName(f.Name)
		outF := out.FieldByName(f.Name)
		if inF.Kind() == reflect.Invalid {
			return fmt.Errorf("recode: Recode() struct %v has no field %v", in.Type().Name(), f.Name)
		}
		err := recodeField(inF, outF)
		if err != nil {
			return err
		}

	}
	return nil
}

func recodeField(in, out reflect.Value) error {
	switch out.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return recodeInt(in, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return recodeUint(in, out)
	case reflect.String:
		return recodeString(in, out)
	case reflect.Float32, reflect.Float64:
		return recodeFloat(in, out)
	case reflect.Struct:
		return recodeStruct(in, out)
	case reflect.Interface, reflect.Ptr:
		return recodeField(in, out.Elem())
	case reflect.Slice:
		return recodeSlice(in, out)
	default:
		return fmt.Errorf("recode: Recode() cannot map to %v", out.Kind())
	}
}

func recodeInt(in, out reflect.Value) error {
	switch in.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetInt(in.Int())
	case reflect.String:
		i, err := strconv.ParseInt(in.String(), 10, 64)
		if err != nil {
			return err
		}
		out.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.SetInt(int64(in.Uint()))
	case reflect.Float32, reflect.Float64:
		out.SetInt(int64(in.Float()))
	case reflect.Interface, reflect.Ptr:
		return recodeInt(in.Elem(), out)
	default:
		return fmt.Errorf("recode: Recode() cannot map %v to %v", in.Kind(), out.Kind())
	}
	return nil
}

func recodeUint(in, out reflect.Value) error {
	switch in.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.SetUint(in.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetUint(uint64(in.Int()))
	case reflect.Float32, reflect.Float64:
		out.SetUint(uint64(in.Float()))
	case reflect.String:
		i, err := strconv.ParseUint(in.String(), 10, 64)
		if err != nil {
			return err
		}
		out.SetUint(i)
	case reflect.Interface, reflect.Ptr:
		return recodeUint(in.Elem(), out)
	default:
		return fmt.Errorf("recode: Recode() cannot map %v to %v", in.Kind(), out.Kind())
	}
	return nil
}

func recodeString(in, out reflect.Value) error {
	switch in.Kind() {
	case reflect.String:
		out.SetString(in.String())
	case reflect.Interface, reflect.Ptr:
		return recodeString(in.Elem(), out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.SetString(strconv.FormatUint(in.Uint(), 10))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetString(strconv.FormatInt(in.Int(), 10))
	case reflect.Float32, reflect.Float64:
		out.SetString(strconv.FormatFloat(in.Float(), 'f', -1, 64))
	default:
		return fmt.Errorf("recode: Recode() cannot map %v to %v", in.Kind(), out.Kind())
	}
	return nil
}

func recodeFloat(in, out reflect.Value) error {
	switch in.Kind() {
	case reflect.Float32, reflect.Float64:
		out.SetFloat(in.Float())
	case reflect.String:
		f, err := strconv.ParseFloat(in.String(), 64)
		if err != nil {
			return err
		}
		out.SetFloat(f)
	case reflect.Interface, reflect.Ptr:
		return recodeFloat(in.Elem(), out)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetFloat(float64(in.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.SetFloat(float64(in.Uint()))
	default:
		return fmt.Errorf("recode: Recode() cannot map %v to %v", in.Kind(), out.Kind())
	}
	return nil
}

func recodeSlice(in, out reflect.Value) error {
	switch in.Kind() {
	case reflect.Slice:
		out.Set(reflect.MakeSlice(out.Type(), in.Len(), in.Len()))
		for i := 0; i < in.Len(); i++ {
			err := recodeField(in.Index(i), out.Index(i))
			if err != nil {
				return err
			}
		}
	case reflect.Interface, reflect.Ptr:
		return recodeSlice(in.Elem(), out)
	default:
		return fmt.Errorf("recode: Recode() cannot map %v to %v", in.Kind(), out.Kind())
	}
	return nil
}
