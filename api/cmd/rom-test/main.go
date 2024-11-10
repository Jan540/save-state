package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type NDSROMHeader struct {
	GameTitle [12]byte
	GameCode  [4]byte
	MakerCode [2]byte
}

func main() {
	fmt.Printf("backup%s.sav\n", time.Now().Format("020106-150405"))

	path := "/home/jann/games/roms"

	saveDir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("failed to read dir")
		os.Exit(1)
	}

	for _, dirEntry := range saveDir {
		if dirEntry.IsDir() {
			continue
		}

		ext := filepath.Ext(dirEntry.Name())

		if ext != ".nds" {
			continue
		}

		file, err := os.Open(filepath.Join(path, dirEntry.Name()))

		if err != nil {
			fmt.Println("failed to open file")
			os.Exit(1)
		}

		defer file.Close()

		var header NDSROMHeader

		if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
			fmt.Println("failed to read header")
			os.Exit(1)
		}

		fmt.Printf("Game Title: %s; Game Code: %s\n", header.GameTitle, header.GameCode)
	}
}
