package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "wb",
	Short: "read whiteboard",
	Run:   rootRun,
}

func init() {
}

func getWbPath() (string, error) {
	curPath, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	goPath, _ := os.LookupEnv("GOPATH")

	relPath, err := filepath.Rel(curPath, path.Join(goPath, "/src/github.com/rigelrozanski/wb/whiteboard"))
	return relPath, err
}

func rootRun(cmd *cobra.Command, args []string) {

	wbPath, err := getWbPath()
	if err != nil {
		fmt.Println(err.Error())
	}

	if _, err := os.Stat(wbPath); os.IsNotExist(err) { //does the whiteboard file not yet exist?
		err = ioutil.WriteFile("whiteboard", []byte("Squeaky Clean Whiteboard \n"), 0644)
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
