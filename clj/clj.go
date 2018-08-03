// simply delegate to clojure.exe
package main

import (
	"os"

	"os/exec"
)

func main() {
	var cmd *exec.Cmd
	var cmd_args []string

	cmd_args = os.Args[1:]
	cmd = exec.Command("clojure.exe", cmd_args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

}
