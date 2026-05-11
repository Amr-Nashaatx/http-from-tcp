package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	// go routine that reads a line
	go func() {
		var line string
		buf := make([]byte, 8)

		for {
			n, readErr := f.Read(buf)
			if n > 0 {
				tmp := string(buf[:n])
				parts := strings.Split(tmp, "\n")
				// if read string contains no new line "\n"
				// we just concat the line string
				if len(parts) == 1 {
					line += parts[0]
					continue
				}
				// if read string contains multiple new line character "\n"
				// we add each part and print the whole line string and reset it, for each part except the last one we ignore it.
				for i := range len(parts) {
					if i == len(parts)-1 {
						break
					}
					line = line + parts[i]
					lines <- line
					line = ""
				}
			}
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					close(lines)
					f.Close()
					return
				} else {
					panic(readErr)
				}

			}
		}
	}()

	return lines
}
func main() {
	f, err := os.Open("./message.txt")

	if err != nil {
		panic(fmt.Errorf("could not open file %w", err))
	}

	linesChan := getLinesChannel(f)
	for line := range linesChan {
		fmt.Println("read: ", line)
	}
}
