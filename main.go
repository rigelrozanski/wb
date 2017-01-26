package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
			view("list")
		case "new":
			//error new is a reserved word
			fmt.Println("invalid argments, must specify name of new board as additional argument")
		default:
			// open the wb board with the name of the argument
			view(args[0])
		}

	case 2:
		//edit or create a specified board

		//get edit position
		editArg := -1
		newArg := -1
		boardArg := -1
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "edit":
				editArg = i
			case "new":
				newArg = i
			default:
				boardArg = i
			}
		}

		switch {
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

	relPath, err := filepath.Rel(curPath, path.Join(goPath,
		"/src/github.com/rigelrozanski/wb", wbName))
	return relPath, err
}

func wbExists(wbPath string) bool {
	_, err := os.Stat(wbPath)
	return !os.IsNotExist(err)
}

func addNameToList(wbName string) {
	//TODO: complete funtionality
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

	addNameToList(wbName)
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

	if !wbExists(wbPath) { //does the whiteboard file not yet exist?
		fmt.Println("error can't view non-existent white board, please create it first by using 'new'")
	} else {
		wb, err := ioutil.ReadFile(wbPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(string(wb))
	}
}
