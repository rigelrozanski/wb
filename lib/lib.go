package lib

import (
	"fmt"
	pathL "path"

	cmn "github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var BoardsDir = "boards"

// get the contents of a local wb
func GetWB(name string) []string {

	path, err := GetWbPath(name)
	if err != nil {
		panic(err)
	}
	if !cmn.FileExists(path) {
		panic("wb no exist")
	}

	lines, err := cmn.ReadLines(path)
	if err != nil {
		panic(err)
	}

	return lines
}

// get the full path of a wb
func GetWbPath(wbName string) (string, error) {
	wbBackupRepoPath, err := GetWbBackupRepoPath()
	if err != nil {
		return "", err
	}
	return pathL.Join(wbBackupRepoPath, BoardsDir, wbName), nil
}

// get the full path of a wb backup repo
func GetWbBackupRepoPath() (string, error) {
	configPath, err := cmn.GetRelPath("/src/github.com/rigelrozanski/wb", "config.txt")
	if err != nil {
		return "", fmt.Errorf("missing config.txt file in root of repo, error: %v", err)
	}
	lines, err := cmn.ReadLines(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config, error: %v", err)
	}
	return lines[0], nil
}
