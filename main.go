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

//keywords used throughout wb
const (
	keyNew     = "nu"
	keyNew2    = "nu2"
	keyView    = "cat"
	keyRemove  = "rm"
	keyBackup  = "backup"
	keyRestore = "restore"
	keyList    = "list"
	keyHelp1   = "--help"
	keyHelp2   = "-h"

	defaultWB = "wb"
	boardsDir = "boards"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		// open the main wb
		edit(defaultWB)
	case 1:
		switch args[0] {
		case keyHelp1, keyHelp2:
			printHelp()
		case keyBackup:
			backup()
		case keyRestore:
			restore()
		case keyView:
			view(defaultWB)
		case keyList:
			list()
		case keyNew, keyNew2, keyRemove:
			fmt.Println("invalid argments, must specify name of board")
		default:
			// open the wb board with the name of the argument
			edit(args[0])
		}

	case 2:
		//edit/delete/create-new board
		Bview := false
		Bdelete := false
		Bnew := false
		Bnew2 := false
		noRsrvArgs := 0

		boardArg := -1
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case keyView:
				Bview = true
				noRsrvArgs++
			case keyRemove:
				Bdelete = true
				noRsrvArgs++
			case keyNew:
				Bnew = true
				noRsrvArgs++
			case keyNew2:
				Bnew2 = true
				noRsrvArgs++
			case keyList:
				fmt.Println("invalid argments, list argument is reserved")
				return
			default:
				boardArg = i
			}
		}
		switch {
		case noRsrvArgs != 1:
			fmt.Printf("invalid use of reserved arguments, must enter *one*"+
				" of either %v, %v, or %v\n", keyNew, keyView, keyRemove)
			return
		case Bnew:
			freshWB(args[boardArg])
		case Bnew2:
			freshWBFilled(args[boardArg])
		case Bview:
			view(args[boardArg])
		case Bdelete:
			remove(args[boardArg])
		}
	default:
		fmt.Println("invalid number of args")
		return
	}
}

func printHelp() {
	fmt.Println(`
/|||||\ |-o-o-~|

Usage: 

wb [boardname]       -> edit the wb
wb nu [boardname]    -> create a new wb
wb nu2 [boardname]   -> create a new filled wb
wb cat [boardname]   -> view a wb
wb rm [boardname]    -> remove a wb
wb backup            -> aws backup all wbs
wb restore           -> aws restore all wbs
wb list              -> list all the wbs

note - if the [boardname] is not provided, 
the default board named 'wb' will be used.
`)
}

func getWbPath(wbName string) (string, error) {
	return getRelPath(pathL.Join("/src/github.com/rigelrozanski/wb", boardsDir), wbName)
}

func getKeyPath() (string, error) {
	return getRelPath("/src/github.com/rigelrozanski/wb", "key.json")
}

func getRelPath(absPath, file string) (string, error) {
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

func wbExists(wbPath string) bool {
	_, err := os.Stat(wbPath)
	return !os.IsNotExist(err)
}

func list() {

	boardPath, err := getWbPath("")
	if err != nil {
		fmt.Println(err)
		return
	}

	filepath.Walk(boardPath, visit)
}

func visit(path string, f os.FileInfo, err error) error {

	basePath := pathL.Base(path)
	basePath = strings.Replace(basePath, boardsDir, "", 1) //remove the boards dir
	if len(basePath) > 0 {
		fmt.Println(basePath)
	}
	return nil
}

func remove(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err)
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

func freshWBFilled(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if wbExists(wbPath) { //does the whiteboard already exist
		fmt.Println("error whiteboard already exists")
		return
	}

	//create the blank canvas to work from
	blankLine := fmt.Sprintf("%140s\n", " ")
	var blankStr string
	for i := 1; i <= 48; i++ {
		blankStr += blankLine
	}

	err = ioutil.WriteFile(wbPath, []byte(blankStr), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Squeaky clean whiteboard created for", wbName)

	//now edit the wb
	edit(wbName)
}

func freshWB(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if wbExists(wbPath) { //does the whiteboard already exist
		fmt.Println("error whiteboard already exists")
		return
	}

	//create the blank canvas to work from
	err = ioutil.WriteFile(wbPath, []byte(""), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Squeaky clean whiteboard created for", wbName)

	//now edit the wb
	edit(wbName)
}

func edit(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !wbExists(wbPath) {
		fmt.Println("error can't edit non-existent white board, please create it first by using ", keyNew)
		return
	}

	//cmd2 := exec.Command("vim", "-c", "startreplace | +normal 25G70|", wbPath) //start with replace
	cmd2 := exec.Command("vim", "-c", "+normal 1G1|", wbPath) //start in the upper left corner nomatter
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	err = cmd2.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func view(wbName string) {
	wbPath, err := getWbPath(wbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch {
	case !wbExists(wbPath) && wbName == defaultWB:
		freshWB(defaultWB) //automatically create the default wb if it doesn't exist
	case !wbExists(wbPath) && wbName != defaultWB:
		fmt.Println("error can't view non-existent white board, please create it first by using ", keyNew)
	default:
		wb, err := ioutil.ReadFile(wbPath)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(string(wb))
	}
}
