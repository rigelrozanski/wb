package lib

import (
	"fmt"
	"io/ioutil"
	pathL "path"
	"strings"

	cmn "github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var BoardsDir = "boards"

// get the contents of a local wb
func GetWB(name string) (content []string, found bool) {

	path, err := GetWbPath(name)
	if err != nil {
		return content, false
	}
	if !cmn.FileExists(path) {
		return content, false
	}

	content, err = cmn.ReadLines(path)
	if err != nil {
		return content, false
	}

	return content, true
}

func todoStr(name string) string {
	return fmt.Sprintf("TODO add: [%v]", name)
}

// get the contents of a local wb
func RemoveFromLS(lsname, remove string) error {

	path, err := GetWbPath(lsname)
	if err != nil {
		return err
	}
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fileStr := string(bz)

	fileStr = strings.Replace(fileStr, "\n"+todoStr(remove), "", 1)
	fileStr = strings.Replace(fileStr, fmt.Sprintf("[%s]", remove), "", 1)
	err = ioutil.WriteFile(path, []byte(fileStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

// get the contents of a local wb
func AddToLS(lsName, newWB string) error {
	return PrependWB(lsName, todoStr(newWB))
}

// nolint
func MustPrependWB(wbName, entry string) {
	err := PrependWB(wbName, entry)
	if err != nil {
		panic(err)
	}
}

// prepend a string to a new top line within a wb
func PrependWB(wbName, entry string) error {

	path, err := GetWbPath(wbName)
	if err != nil {
		return err
	}
	if !cmn.FileExists(path) {
		return err
	}

	content, err := cmn.ReadLines(path)
	if err != nil {
		return err
	}

	content = append([]string{entry}, content...)

	err = cmn.WriteLines(content, path)
	if err != nil {
		return err
	}

	return nil

}

// nolint
func MustClearWB(wbName string) {
	err := ClearWB(wbName)
	if err != nil {
		panic(err)
	}
}

// remove all the content of wb
func ClearWB(wbName string) error {
	path, err := GetWbPath(wbName)
	if err != nil {
		return err
	}
	if !cmn.FileExists(path) {
		return err
	}

	err = cmn.WriteLines([]string{}, path)
	if err != nil {
		return err
	}

	return nil
}

// get the contents of a local wb
func WbExists(name string) (found bool) {

	path, err := GetWbPath(name)
	if err != nil {
		return false
	}
	if !cmn.FileExists(path) {
		return false
	}
	return true
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

// get the full path of a wb backup repo
func GetBoardsGITDir() (string, error) {
	wbBackupRepoPath, err := GetWbBackupRepoPath()
	if err != nil {
		return "", err
	}
	return pathL.Join(wbBackupRepoPath, ".git"), nil
}

// function for iterating over the wbs
type IterateFn func(name, relPath string) (stop bool)

// perform the provided iterateFn for all wbs
func IterateWBs(iterFn IterateFn) error {

	boardPath, err := GetWbPath("")
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(boardPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		relPath, err := GetWbPath(name)
		if err != nil {
			return err
		}
		if iterFn(name, relPath) {
			break
		}
	}
	return nil
}
