package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var wbName string

var RootCmd = &cobra.Command{
	Use:   "wb",
	Short: "read whiteboard",
	Run:   rootRun,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&wbName, "name", "b", "wb", "specific whiteboard")
}

func getWbPath() (string, error) {
	curPath, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	goPath, _ := os.LookupEnv("GOPATH")

	relPath, err := filepath.Rel(curPath, path.Join(goPath,
		"/src/github.com/rigelrozanski/wb",
		wbName))
	return relPath, err
}

func rootRun(cmd *cobra.Command, args []string) {

	wbPath, err := getWbPath()
	if err != nil {
		fmt.Println(err.Error())
	}

	if _, err := os.Stat(wbPath); os.IsNotExist(err) { //does the whiteboard file not yet exist?
		err = ioutil.WriteFile(wbName, []byte("Squeaky clean whiteboard for "+wbName+"\n"), 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	wb, err := ioutil.ReadFile(wbPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Print(string(wb))
}
