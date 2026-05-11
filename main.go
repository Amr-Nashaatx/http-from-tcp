package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("./message.txt")

	if err != nil {
		panic(fmt.Errorf("could not open file %w", err))
	}

	buf := make([]byte, 8)

	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			fmt.Printf("Read: %s\n", string(buf[:n]))
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				fmt.Println("File ended!")
				return
			} else {
				panic(readErr)
			}

		}
	}
}
