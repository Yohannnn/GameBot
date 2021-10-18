package core

import "reflect"

//Contains
//Checks if an array contains an element
func Contains(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

//IntArray
//Returns an array that starts at one value and end at another
func IntArray(x, y int) []int {
	var a []int
	z := x
	for z <= y {
		a = append(a, z)
		z++
	}
	return a
}
