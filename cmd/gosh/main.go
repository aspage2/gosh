package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func mainloop() {

	reader := bufio.NewReader(os.Stdin)
	userHome := os.Getenv("HOME")
	lastExit := 0
	for {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		psDir := strings.Replace(cwd, userHome, "~", 1)
		fmt.Printf("%s $ ", psDir)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		flds := strings.Fields(input)
		if len(flds) == 0{
			continue
		}
		cmdName := flds[0]
		args := flds[1:]
		switch cmdName {
		case "exit":
			return
		case "lastexit":
			fmt.Printf("%d\n", lastExit)
		}
		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err = cmd.Run(); err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				fmt.Printf("error: %s\n",  err)
			} else {
				lastExit = exitErr.ExitCode()
			}
		}
	}
}

func main() {
	mainloop()

}
