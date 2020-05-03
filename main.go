package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		flds := strings.Fields(input)

		proc := exec.Command(flds[0], flds[1:]...)
		proc.Stdout = os.Stdout
		err = proc.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}
