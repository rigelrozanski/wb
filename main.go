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
		default:
			// open the wb board with the name of the argument
			view(args[0])
		}

	case 2:
		//edit a specified board

		//get edit position
		editArg := -1
		boardArg := -1
		for i := 0; i < len(args); i++ {
			if args[i] == "edit" {
				editArg = i
			} else {
				boardArg = i
			}
		}
		if editArg == -1 {
			fmt.Println("invalid argments, 2 arguments without an 'edit' arg")
			return
		}

		edit(args[boardArg])
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

func edit(wbName string) {

	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err.Error())
	}

	if !wbExists(wbPath) {
		addNameToList(wbName)
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
		addNameToList(wbName)
		err = ioutil.WriteFile(wbPath, []byte(""), 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Squeaky clean whiteboard created for", wbName)
	} else {
		wb, err := ioutil.ReadFile(wbPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(string(wb))
	}
}
