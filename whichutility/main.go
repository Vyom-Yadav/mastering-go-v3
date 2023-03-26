package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Provide more arguments for which utility")
		return
	}
	executable := args[1]
	path := os.Getenv("PATH")
	pathSplit := filepath.SplitList(path)
	result := make([]string, 0, 0)
	for _, dir := range pathSplit {
		filePath := filepath.Join(dir, executable)
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			mode := fileInfo.Mode()
			if mode.IsRegular() {
				// Is an executable?
				if mode&0111 != 0 {
					result = append(result, filePath)
				}
			}
		}
	}

	if len(result) != 0 {
		for i, v := range result {
			fmt.Printf("%d: %s\n", i, v)
		}
	}
}
