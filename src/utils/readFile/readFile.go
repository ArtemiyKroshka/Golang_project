package readFile

import (
	"bufio"
	"log"
	"os"
	"sync"
)

func GetStrings(fileName string, mu *sync.RWMutex) []string {

	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}

	return lines

}
