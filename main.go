package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rigelrozanski/wb/lib"
)

//keywords used throughout wb
const (
	keyNew        = "new"
	keyCopy       = "cp"
	keyView       = "cat"
	keyRemove     = "rm"
	keyRecover    = "recover"
	keyEmptyTrash = "empty-trash"
	keyList       = "ls"
	keyLog        = "log"
	keyStats      = "stats"
	keyPush       = "push"

	keyHelp1 = "--help"
	keyHelp2 = "-h"

	defaultWB   = "wb"
	lsWB        = "lsls"
	logWB       = "loglog"
	shortcutsWB = "shortcuts"

	help = `
/|||||\ |-o-o-~|

Usage: 

wb [name]            -> vim into a wb
wb [name] [entry]    -> fast entry appended as new line in wb
wb new [name]        -> create a new wb
wb cp [copy] [name]  -> duplicate a wb
wb cat [name]        -> print wb contents to console
wb rm [name]         -> remove a wb (add to trash)
wb recover [name]    -> remove a wb from trash
wb empty-trash       -> empty trash
wb ls                -> list all the wb in console
wb log               -> list the log
wb stats             -> list git statistics per wb
wb push [msg]        -> git push the boards directory

notes:
- if the [name] is not provided, 
  the default board named 'wb' will be used
- special reserved wb names: wb, lsls, loglog

`
)

func main() {
	args := os.Args[1:]

	var err error
	switch len(args) {
	case 0:
		err = openDefaultWB()
	case 1:
		err = handle1Args(args)
	case 2:
		err = handle2Args(args)
	case 3:
		err = handle3Args(args)
	default:
		name := args[0]
		entry := strings.Join(args[1:], " ")
		err = fastEntry(name, entry)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func openDefaultWB() error {
	// open the main wb
	modified, err := edit(defaultWB)
	if err != nil {
		return err
	}
	if modified {
		log("modified wb", defaultWB)
	}
	return nil
}

func handle1Args(args []string) error {
	if len(args) != 1 {
		panic("improper args")
	}
	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyPush:
		err := push(fmt.Sprintf("%v", time.Now()))
		if err != nil {
			return err
		}
		lib.MustClearWB(logWB)
		log("pushed", "n/a")
	case keyView:
		err := view(defaultWB)
		if err != nil {
			return err
		}
	case keyList:
		err := list()
		if err != nil {
			return err
		}
	case keyEmptyTrash:
		err := emptyTrash()
		if err != nil {
			return err
		}
		log("emptied trash", "n/a")
	case keyLog:
		err := listLog()
		if err != nil {
			return err
		}
	case keyStats:
		err := listStats()
		if err != nil {
			return err
		}
	case keyNew, keyRemove:
		fmt.Println("invalid argments, must specify name of board")
	default:
		// open the wb board with the name of the argument
		name := args[0]
		modified, err := edit(name)
		if err != nil {
			return err
		}
		if modified {
			log("modified wb", name)
		}
	}

	return nil
}

// TODO this is spagetti - fix!
func handle2Args(args []string) error {
	if len(args) != 2 {
		panic("improper args")
	}

	//edit/delete/create-new board
	Bview, Bdelete, Brecover, Bnew := false, false, false, false
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
		case keyRecover:
			Brecover = true
			noRsrvArgs++
		case keyNew:
			Bnew = true
			noRsrvArgs++
		case keyList, keyEmptyTrash:
			break
		default:
			boardArg = i
		}
	}

	switch {
	case Bnew:
		name := args[boardArg]
		err := freshWB(name)
		if err != nil {
			return err
		}
		log("created wb", name)
	case Bview:
		err := view(args[boardArg])
		if err != nil {
			return err
		}
	case Bdelete:
		name := args[boardArg]
		err := remove(name)
		if err != nil {
			return err
		}
		log("deleted wb", name)
	case Brecover:
		name := args[boardArg]
		err := recoverWb(name)
		if err != nil {
			return err
		}
		log("recovered wb", name)
	default:
		name := args[0]
		entry := args[1]
		return fastEntry(name, entry)
	}

	return nil
}

func handle3Args(args []string) error {
	if len(args) != 3 {
		panic("improper args")
	}

	switch args[0] {
	case keyCopy:
		err := duplicate(args[1], args[2])
		if err != nil {
			return err
		}
		log("duplicated from "+args[1], args[2])
	default:
		name := args[0]
		entry := strings.Join(args[1:], " ")
		return fastEntry(name, entry)
	}

	return nil
}
