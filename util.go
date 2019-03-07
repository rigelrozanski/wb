package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	pathL "path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rigelrozanski/wb/lib"

	cmn "github.com/rigelrozanski/common"
)

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

type stat struct {
	name      string
	additions int
	deletions int
}

func listStats() error {
	wbDir, err := lib.GetWbBackupRepoPath()
	if err != nil {
		return err
	}

	var stats []stat

	iterFn := func(name, relPath string) (stop bool) {

		getStatsCmd := `git -C ` + wbDir + ` log` + ` --pretty=tformat: --numstat ` + relPath
		out, err := cmn.Execute(getStatsCmd)
		if err != nil {
			panic(err)
		}
		lines := strings.Split(out, "\n")
		additionsTot, deletionsTot := 0, 0
		for _, line := range lines {
			line := strings.Split(line, "\t")
			if len(line) != 3 {
				continue
			}
			additions, err := strconv.Atoi(line[0])
			if err != nil {
				continue
			}
			deletions, err := strconv.Atoi(line[1])
			if err != nil {
				continue
			}
			additionsTot += additions
			deletionsTot += deletions
		}
		stats = append(stats, stat{name, additionsTot, deletionsTot})
		if err != nil {
			panic(err)
		}
		return false
	}
	lib.IterateWBs(iterFn)

	// sort and print
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].additions < stats[j].additions
	})
	fmt.Println("Add\tDel\tName")
	for _, stat := range stats {

		// ignore some special wb
		if strings.Contains(stat.name, ".") ||
			stat.name == lsWB ||
			stat.name == logWB {

			continue
		}
		fmt.Printf("%v\t%v\t%v\n", stat.additions,
			stat.deletions, stat.name)
	}

	return nil
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
	err := lib.MoveWbToTrash(wbName)
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

func recoverWb(wbName string) error {
	err := lib.RecoverWbFromTrash(wbName)
	if err != nil {
		return err
	}

	err = lib.AddToLS(lsWB, wbName)
	if err != nil {
		return err
	}

	fmt.Println("roger, recovered")
	return nil
}

func emptyTrash() error {
	err := lib.EmptyTrash()
	if err != nil {
		return err
	}
	fmt.Println("roger, emptied the trash")
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
	_, _, err = edit(wbName)
	return err
}

func edit(wbName string) (modified bool, nameUpdate string, err error) {

	origContent, found := lib.GetWB(wbName)
	if !found {
		// check for shortcuts
		shortcuts, foundSC := lib.GetWB(shortcutsWB)
		if !foundSC {
			return false, wbName,
				fmt.Errorf("that wb is not found (nor the shortcuts)")
		}

		shortcutFound := false
		for _, shortcut := range shortcuts {
			str := strings.Fields(shortcut)
			if len(str) < 3 {
				continue
			}
			if str[0] == wbName && str[1] == "->" { // shortcut found
				shortcutFound = true
				wbName = str[2]
				break
			}
		}
		if !shortcutFound {
			return false, wbName,
				fmt.Errorf("error can't edit non-existent white board, please create it first by using %v", keyNew)
		}
	}

	wbPath, err := lib.GetWbPath(wbName)
	if err != nil {
		return false, wbName, err
	}

	cmd := exec.Command("vim", "-c", "+normal 1G1|", wbPath) //start in the upper left corner nomatter
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return false, wbName, err
	}

	// determine if was modified
	newContent, found := lib.GetWB(wbName)
	if !found {
		panic("wuz found now isn't")
	}

	if len(newContent) != len(origContent) {
		return true, wbName, nil
	}

	for i, line := range origContent {
		if line != newContent[i] {
			return true, wbName, nil
		}
	}
	return false, wbName, nil
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
