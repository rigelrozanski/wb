package lib

import (
	"os"
	pathL "path"
	"path/filepath"

	"github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var BoardsDir = "boards"

// get the contents of a local wb
func GetWB(name string) []string {

	path, err := GetWbPath(name)
	if err != nil {
		panic(err)
	}
	if !WbExists(path) {
		panic("wb no exist")
	}

	lines, err := common.ReadLines(path)
	if err != nil {
		panic(err)
	}

	return lines
}

// get the full path of a wb
func GetWbPath(wbName string) (string, error) {
	return GetRelPath(pathL.Join("/src/github.com/rigelrozanski/wb", BoardsDir), wbName)
}

// get the path for the key file
func GetKeyPath() (string, error) {
	return GetRelPath("/src/github.com/rigelrozanski/wb", "key.json")
}

// get the relative path to current loco
func GetRelPath(absPath, file string) (string, error) {
	curPath, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	goPath, _ := os.LookupEnv("GOPATH")

	relBoardsPath, err := filepath.Rel(curPath, pathL.Join(goPath,
		absPath))

	//create the boards directory if it doesn't exist
	os.Mkdir(relBoardsPath, os.ModePerm)

	relWbPath := pathL.Join(relBoardsPath, file)

	return relWbPath, err
}

// does the wb at this path exist
func WbExists(wbPath string) bool {
	_, err := os.Stat(wbPath)
	return !os.IsNotExist(err)
}
