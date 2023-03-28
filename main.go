package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	rootDir := "./" // gyökérkönyvtár

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ha az elérési út a Go modul gyökérkönyvtára, akkor futtassuk a 'go mod tidy' parancsot
		if info.IsDir() && info.Name() == "go.mod" {
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = filepath.Dir(path)

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Output: %s\n", output)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
