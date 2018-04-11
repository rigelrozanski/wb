package main

import "github.com/rigelrozanski/common"

// get the contents of a local wb
func GetWB(name string) []string {

	path, err := getWbPath(name)
	if err != nil {
		panic(err)
	}
	if !wbExists(path) {
		panic("wb no exist")
	}

	lines, err := common.ReadLines(path)
	if err != nil {
		panic(err)
	}

	return lines
}
