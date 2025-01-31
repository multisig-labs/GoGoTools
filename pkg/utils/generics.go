// Copyright (C) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import "reflect"

// Safe nil checking on an interface, that does not panic
func IsNil(v interface{}) bool {
	tv := reflect.ValueOf(v)
	switch tv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return tv.IsNil()
	}
	return false
}

// AppendSlices appends multiple slices into a single slice.
func AppendSlices[T any](slices ...[]T) []T {
	totalLength := 0
	for _, slice := range slices {
		totalLength += len(slice)
	}
	result := make([]T, 0, totalLength)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func Unique[T comparable](arr []T) []T {
	visited := map[T]bool{}
	unique := []T{}
	for _, e := range arr {
		if !visited[e] {
			unique = append(unique, e)
			visited[e] = true
		}
	}
	return unique
}

func Any[T any](input []T, f func(T) bool) bool {
	for _, e := range input {
		if f(e) {
			return true
		}
	}
	return false
}

func Find[T any](input []T, f func(T) bool) *T {
	for _, e := range input {
		if f(e) {
			return &e
		}
	}
	return nil
}

func Belongs[T comparable](input []T, elem T) bool {
	for _, e := range input {
		if e == elem {
			return true
		}
	}
	return false
}

func RemoveFromSlice[T comparable](input []T, toRemove T) []T {
	output := make([]T, 0, len(input))
	for _, e := range input {
		if e != toRemove {
			output = append(output, e)
		}
	}
	return output
}

func Filter[T any](input []T, f func(T) bool) []T {
	output := make([]T, 0, len(input))
	for _, e := range input {
		if f(e) {
			output = append(output, e)
		}
	}
	return output
}

func Map[T, U any](input []T, f func(T) U) []U {
	output := make([]U, 0, len(input))
	for _, e := range input {
		output = append(output, f(e))
	}
	return output
}

func MapWithError[T, U any](input []T, f func(T) (U, error)) ([]U, error) {
	output := make([]U, 0, len(input))
	for _, e := range input {
		o, err := f(e)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}
