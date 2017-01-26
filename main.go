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
		case "new":
			//error new is a reserved word
			fmt.Println("invalid argments, must specify name of new board as additional argument")
		default:
			// open the wb board with the name of the argument
			view(args[0])
		}

	case 2:
		//edit or create a specified board

		editArg := -1
		newArg := -1
		listArg := -1
		boardArg := -1
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "edit":
				editArg = i
			case "new":
				newArg = i
			case "list":
				listArg = i
			default:
				boardArg = i
			}
		}

		switch {
		case listArg != -1:
			fmt.Println("invalid argments, list argument is reserved")
			return
		case editArg == -1 && newArg != -1:
			new(args[boardArg])
		case editArg != -1 && newArg == -1:
			edit(args[boardArg])
		case editArg != -1 && newArg != -1:
			fmt.Println("invalid argments, specified to 'edit' and create 'new'")
			return
		case editArg == -1 && newArg == -1:
			fmt.Println("invalid argments, 2 arguments without an 'edit' or 'new' arg")
			return
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

func new(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
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
