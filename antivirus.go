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

	for _, file := range extractKeys(files) {
		if file == "" {
			wg.Done()
			continue
		}
		input := fmt.Sprintf(file + "/" + files[file][0])
		fmt.Println("--------------")
		fmt.Println(input)
		fmt.Println("--------------")
		// fmt.Println(input)
		go func() {
			out2, err2 := exec.Command("md5", input).Output()
			if err2 != nil {
				wg.Done()
				fmt.Println("Finished with (failed) " + input)
				return
				// log.Fatal(err2)
			}
			fileHash := strings.TrimSpace(string(strings.Split(string(out2), " = ")[1]))

			if (listContains(fileHash, maliciousHashes)) {
				fmt.Printf("Found suspicious file: %s", string(out2))
			}
			fmt.Println("Finished with " + input)
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

func RetrieveFilenames(wg *sync.WaitGroup) map[string][]string {
	fmt.Println("Getting target files...")
	m := make(map[string][]string)

	out, err := exec.Command("ls", "-Ra", "./Target").Output()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n./")

	for _, directory := range lines[1:] {
		dirname := strings.Split(directory, "\n")[0]
		dirname = dirname[0:len(dirname)-1]

		m[dirname] = strings.Split(directory, "\n")[3:]
	}
	fmt.Println(m)
	wg.Add(len(lines))


	fmt.Println("Found target files...")
	return m
}

func extractKeys(m map[string][]string) []string {
	keys := []string{}
	for k := range m {
    keys = append(keys, k)
	}
	return keys
}
