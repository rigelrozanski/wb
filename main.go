package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	pathL "path"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		// open the main wb
		view("wb")
	case 1:
		switch args[0] {
		case "edit":
			//edit the main wb
			edit("wb")
		case "list":
			//view the list
			list()
		case "delete":
			fmt.Println("invalid argments, must specify name of the board to delete as additional argument")
		case "new":
			fmt.Println("invalid argments, must specify name of new board as additional argument")
		default:
			// open the wb board with the name of the argument
			view(args[0])
		}

	case 2:
		//edit/delete/create-new board
		Bedit := false
		Bdelete := false
		Bnew := false
		noRsrvArgs := 0

		boardArg := -1
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "edit":
				Bedit = true
				noRsrvArgs++
			case "delete":
				Bdelete = true
				noRsrvArgs++
			case "new":
				Bnew = true
				noRsrvArgs++
			case "list":
				fmt.Println("invalid argments, list argument is reserved")
				return
			default:
				boardArg = i
			}
		}
		switch {
		case noRsrvArgs != 1:
			fmt.Println("invalid use of reserved arguments, must enter *one* of either 'new', 'edit', or 'delete'")
			return
		case Bnew:
			new(args[boardArg])
		case Bedit:
			edit(args[boardArg])
		case Bdelete:
			delete(args[boardArg])
		}
	default:
		fmt.Println("invalid number of args")
		return
	}
}

func getWbPath(wbName string) (string, error) {
	curPath, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	goPath, _ := os.LookupEnv("GOPATH")

	relBoardsPath, err := filepath.Rel(curPath, pathL.Join(goPath,
		"/src/github.com/rigelrozanski/wb/boards"))

	//create the boards directory if it doesn't exist
	os.Mkdir(relBoardsPath, os.ModePerm)

	relWbPath := pathL.Join(relBoardsPath, wbName)

	return relWbPath, err
}

func wbExists(wbPath string) bool {
	_, err := os.Stat(wbPath)
	return !os.IsNotExist(err)
}

func list() {

	boardPath, err := getWbPath("")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filepath.Walk(boardPath, visit)
}

func visit(path string, f os.FileInfo, err error) error {

	basePath := pathL.Base(path)
	if len(basePath) > 0 && strings.Contains(path, "/") {
		fmt.Println(basePath)
	}
	return nil
}

func delete(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !wbExists(wbPath) { //does the whiteboard not exist
		fmt.Println("error can't delete non-existent whiteboard")
		return
	}

	err = os.Remove(wbPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("roger, deleted successfully")
}

func new(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if wbExists(wbPath) { //does the whiteboard already exist
		fmt.Println("error whiteboard already exists")
		return
	}

	err = ioutil.WriteFile(wbPath, []byte(""), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Squeaky clean whiteboard created for", wbName)
}

func edit(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !wbExists(wbPath) {
		fmt.Println("error can't edit non-existent white board, please create it first by using 'new'")
		return
	}

	cmd2 := exec.Command("vim", wbPath)
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	err = cmd2.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func view(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch {
	case !wbExists(wbPath) && wbName == "wb":
		new("wb") //automatically create the default wb if it doesn't exist
	case !wbExists(wbPath) && wbName != "wb":
		fmt.Println("error can't view non-existent white board, please create it first by using 'new'")
	default:
		wb, err := ioutil.ReadFile(wbPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(string(wb))
	}
}
