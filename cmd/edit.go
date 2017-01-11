package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit whiteboard",
	Run:   editRun,
}

func init() {
	RootCmd.AddCommand(editCmd)
}

func editRun(cmd *cobra.Command, args []string) {

	wbPath, err := getWbPath()
	if err != nil {
		fmt.Println(err.Error())
	}

	cmd2 := exec.Command("vim", wbPath)
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	err = cmd2.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
