package bnp

import "reflect"

// InSlice checks if value is in the given array
func InSlice(val interface{}, slic interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(slic).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slic)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// SliceContains checks if value is in the given slice
func SliceContains(slic interface{}, val interface{}) bool {
	exists, _ := InSlice(val, slic)
	return exists
}

// SliceIndexOf returns the first index found for the given value, or -1
func SliceIndexOf(slic interface{}, val interface{}) int {
	exists, index := InSlice(val, slic)
	if exists {
		return index
	}
	return -1
}
