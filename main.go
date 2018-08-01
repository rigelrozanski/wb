package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	pathL "path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rigelrozanski/wb/lib"

	cmn "github.com/rigelrozanski/common"
)

//keywords used throughout wb
const (
	keyNew    = "new"
	keyCopy   = "cp"
	keyView   = "cat"
	keyRemove = "rm"
	keyList   = "ls"
	keyLog    = "log"
	keyPush   = "push"

	keyHelp1 = "--help"
	keyHelp2 = "-h"

	defaultWB = "wb"
	lsWB      = "lsls"
	logWB     = "loglog" //TODO implement

	help = `
/|||||\ |-o-o-~|

Usage: 

wb [name]            -> vim into a wb
wb new [name]        -> create a new wb
wb cp [copy] [name]  -> duplicate a wb
wb cat [name]        -> print wb contents to console
wb rm [name]         -> remove a wb
wb ls                -> list all the wb in console
wb log               -> list the log
wb push [msg]        -> git push the boards directory

notes:
- if the [name] is not provided, 
  the default board named 'wb' will be used
- special reserved wb names: wb, lsls, loglog 

`
)

func main() {
	args := os.Args[1:]
	errBadArgs := errors.New("invalid number of args")
	var err error
	var modified bool

	switch len(args) {
	case 0:
		// open the main wb
		modified, err = edit(defaultWB)
		if modified {
			log("modified wb", defaultWB)
		}
		break
	case 1:
		switch args[0] {
		case keyHelp1, keyHelp2:
			fmt.Println(help)
		case keyPush:
			err = push(fmt.Sprintf("%v", time.Now()))
			if err == nil {
				lib.MustClearWB(logWB)
				log("pushed", "n/a")
			}
			break
		case keyView:
			err = view(defaultWB)
			break
		case keyList:
			err = list()
			break
		case keyLog:
			err = listLog()
			break
		case keyNew, keyRemove:
			fmt.Println("invalid argments, must specify name of board")
		default:
			// open the wb board with the name of the argument
			name := args[0]
			modified, err = edit(name)
			if modified {
				log("modified wb", name)
			}
			break
		}
	case 2:

		//edit/delete/create-new board
		Bview, Bdelete, Bnew := false, false, false
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
			case keyList:
				break
			default:
				boardArg = i
			}
		}

		switch {
		case noRsrvArgs != 1:
			err = errBadArgs
			break
		case Bnew:
			name := args[boardArg]
			err = freshWB(name)
			log("created wb", name)
			break
		case Bview:
			err = view(args[boardArg])
			break
		case Bdelete:
			name := args[boardArg]
			err = remove(name)
			log("deleted wb", name)
			break
		default:
			err = errBadArgs
		}

	case 3:
		if args[0] != keyCopy {
			err = errBadArgs
		}
		err = duplicate(args[1], args[2])
		log("duplicated from "+args[1], args[2])
		break

	default:
		err = errBadArgs
	}

	if err != nil {
		fmt.Println(err)
	}
}

func log(action, wbName string) {
	lib.MustPrependWB(logWB, fmt.Sprintf("time: %v\taction: %v\t wbName: %v", time.Now(), action, wbName))
}

func list() error {

	if lib.WbExists(lsWB) {
		return view(lsWB)
	}
	boardPath, err := lib.GetWbPath("")
	if err != nil {
		return err
	}
	filepath.Walk(boardPath, visit)
	return nil
}

func listLog() error {
	return view(logWB)
}

// TODO ioutils instead of visit
func visit(path string, f os.FileInfo, err error) error {

	basePath := pathL.Base(path)
	basePath = strings.Replace(basePath, lib.BoardsDir, "", 1) //remove the boards dir
	if len(basePath) > 0 {
		fmt.Println(basePath)
	}
	return nil
}

func remove(wbName string) error {
	wbPath, err := lib.GetWbPath(wbName)
	if err != nil {
		return err
	}

	if !cmn.FileExists(wbPath) { //does the whiteboard not exist
		return errors.New("error can't delete non-existent whiteboard")
	}

	err = os.Remove(wbPath)
	if err != nil {
		return err
	}
	err = lib.RemoveFromLS(lsWB, wbName)
	if err != nil {
		return err
	}

	fmt.Println("roger, deleted successfully")
	return nil
}

func freshWB(wbName string) error {
	wbPath, err := lib.GetWbPath(wbName)
	if err != nil {
		return err
	}

	if cmn.FileExists(wbPath) { //does the whiteboard already exist
		return errors.New("error whiteboard already exists")
	}

	//create the blank canvas to work from
	err = ioutil.WriteFile(wbPath, []byte(""), 0644)
	if err != nil {
		return err
	}
	err = lib.AddToLS(lsWB, wbName)
	if err != nil {
		return err
	}

	fmt.Println("Squeaky clean whiteboard created for", wbName)

	//now edit the wb
	_, err = edit(wbName)
	return err
}

func edit(wbName string) (modified bool, err error) {

	origContent, found := lib.GetWB(wbName)
	if !found {
		return false, fmt.Errorf("error can't edit non-existent white board, please create it first by using %v", keyNew)
	}

	wbPath, err := lib.GetWbPath(wbName)
	if err != nil {
		return false, err
	}

	//if !cmn.FileExists(wbPath) {
	//return false, fmt.Errorf("error can't edit non-existent white board, please create it first by using %v", keyNew)
	//}

	//cmd := exec.Command("vim", "-c", "startreplace | +normal 25G70|", wbPath) //start with replace
	cmd := exec.Command("vim", "-c", "+normal 1G1|", wbPath) //start in the upper left corner nomatter
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return false, err
	}

	// determine if was modified
	newContent, found := lib.GetWB(wbName)
	if !found {
		panic("wuz found now isn't")
	}

	if len(newContent) != len(origContent) {
		return true, nil
	}

	for i, line := range origContent {
		if line != newContent[i] {
			return true, nil
		}
	}
	return false, nil
}

func duplicate(copyWB, newWB string) error {
	copyPath, err := lib.GetWbPath(copyWB)
	if err != nil {
		return err
	}
	if !cmn.FileExists(copyPath) {
		return fmt.Errorf("error can't copy non-existent white board, please create it first by using %v", keyNew)
	}

	newPath, err := lib.GetWbPath(newWB)
	if err != nil {
		return err
	}
	if cmn.FileExists(newPath) {
		return errors.New("error i will not overwrite an existing wb! ")
	}

	err = cmn.Copy(copyPath, newPath)
	if err != nil {
		return err
	}
	err = lib.AddToLS(lsWB, newWB)
	if err != nil {
		return err
	}
	fmt.Printf("succesfully copied from %v to %v\n", copyWB, newWB)
	return nil
}

func view(wbName string) error {
	wbPath, err := lib.GetWbPath(wbName)
	if err != nil {
		return err
	}

	switch {
	case !cmn.FileExists(wbPath) && wbName == defaultWB:
		err := freshWB(defaultWB) //automatically create the default wb if it doesn't exist
		if err != nil {
			return err
		}
	case !cmn.FileExists(wbPath) && wbName != defaultWB:
		fmt.Println("error can't view non-existent white board, please create it first by using ", keyNew)
	default:
		wb, err := ioutil.ReadFile(wbPath)
		if err != nil {
			return err
		}
		fmt.Println(string(wb))
	}
	return nil
}

func push(commitMsg string) error {
	if commitMsg == "" {
		commitMsg = "wb commit"
	}
	wbBackupDir, err := lib.GetWbBackupRepoPath()
	if err != nil {
		return err
	}
	shPath, err := cmn.GetRelPath("/src/github.com/rigelrozanski/wb", "push.sh")
	if err != nil {
		return err
	}
	cmn.WriteLines([]string{`#!/bin/bash
git -C "` + wbBackupDir + `" add -A
git -C "` + wbBackupDir + `" commit -m "` + commitMsg + `"
git -C "` + wbBackupDir + `" push
		`}, shPath)
	cmd := exec.Command("/bin/bash", shPath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println("backup git push complete!")
	return nil
}
