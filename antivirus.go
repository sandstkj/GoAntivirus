package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	maliciousHashesFile := "maliciousHashes.txt"

	var wg sync.WaitGroup

	files := RetrieveFilenames(&wg)

	maliciousHashes := RetrieveMalwareSignatures(maliciousHashesFile)

	for _, file := range files {
		if file == "" {
			wg.Done()
			continue
		}
		input := fmt.Sprintf("./FilePool/%s", file)
		fmt.Println(input)
		go func() {
			out2, err2 := exec.Command("md5", input).Output()
			if err2 != nil {
				log.Fatal(err2)
			}
			fileHash := strings.TrimSpace(string(strings.Split(string(out2), " = ")[1]))

			if (listContains(fileHash, maliciousHashes)) {
				fmt.Printf("Found suspicious file: %s", string(out2))
			}

			wg.Done()
		}()
	}
	wg.Wait()
}

func listContains(element string, list []string) bool {
    for _, current := range list {
        if current == element {
            return true
        }
    }
    return false
}

func RetrieveMalwareSignatures(filename string) []string {
	file, error := os.Open(filename)
    if error != nil {
        return []string{""}
    }
    defer file.Close()

    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := file.Read(buf[0:])
        result = append(result, buf[0:n]...)
        if err != nil {
            if err == io.EOF {
                break
            }
            return []string{""}
        }
    }
		return strings.Split(string(result), "\n")
}

func RetrieveFilenames(wg *sync.WaitGroup) []string {
	out, err := exec.Command("ls", "./FilePool").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	wg.Add(len(lines))
	return lines
}
