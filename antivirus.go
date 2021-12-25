package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("ls", "./FilePool").Output()
	if err != nil {
		log.Fatal(err)
	}
	s := string(out)
	lines := strings.Split(s, "\n")
	fmt.Println(lines)
	for _, line := range lines {
		if line == "" {
			continue
		}
		input := fmt.Sprintf("./FilePool/%s", line)
		fmt.Println(input)
		go func() {
			out2, err2 := exec.Command("md5", input).Output()
			if err2 != nil {
				log.Fatal(err2)
			}
			fmt.Println(string(out2))
		}()
	}
	for true {

	}

	fmt.Println(lines)
}
