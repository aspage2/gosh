package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func mainloop() int {

	reader := bufio.NewReader(os.Stdin)
	userHome := os.Getenv("HOME")
	if userHome == "" {
		userHome = "/"
	}
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
		varsReplaced := GetVars(input, EnvVarSet{})
		flds := strings.Fields(varsReplaced)
		if len(flds) == 0 {
			continue
		}
		cmdName := flds[0]
		args := flds[1:]
		switch cmdName {
		case "exit":
			exitVal := 0
			if len(args) >= 1 {
				exitVal, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
			}
			return exitVal
		case "cd":
			tgt := userHome
			if len(args) > 0 {
				tgt = args[0]
			}
			if err := os.Chdir(tgt); err != nil {
				fmt.Println(err.Error())
			}
			lastExit = 1
			continue
		case "lastexit":
			fmt.Printf("%d\n", lastExit)
			continue
		}
		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err = cmd.Run(); err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				fmt.Printf("error: %s\n", err)
			} else {
				lastExit = exitErr.ExitCode()
			}
		}
	}
}

func main() {
	os.Exit(mainloop())
}
