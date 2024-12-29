package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func encryptFile(filePath string, key []byte, formatCount map[string]int) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	ext := filepath.Ext(filePath)
	format := ext[1:] // Remove the dot from the extension
	if format == "jpeg" {
		format = "jpg"
	}
	formatCount[format]++
	// todays date in format YYYYMMDD
	datePrifix := time.Now().UTC().Format("20060102")
	newFileName := fmt.Sprintf("%s_%s_%04d%s.enc", datePrifix, format, formatCount[format], "."+format)
	outFile, err := os.Create(filepath.Join(filepath.Dir(filePath), newFileName))
	if err != nil {
		return err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, nonce)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	if _, err := outFile.Write(nonce); err != nil {
		return err
	}

	if _, err := io.Copy(writer, inFile); err != nil {
		return err
	}

	return nil
}

func decryptFile(filePath string, key []byte) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Extract the original file name without the .enc extension
	originalFileName := filepath.Base(filePath[:len(filePath)-4])
	// Add the dec_ prefix to the original file name
	newFileName := "dec_" + originalFileName
	// Create the new file path with the dec_ prefix
	outFile, err := os.Create(filepath.Join(filepath.Dir(filePath), newFileName))
	if err != nil {
		return err
	}
	defer outFile.Close()

	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, nonce); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, nonce)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	if _, err := io.Copy(outFile, reader); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: filecrypt <encrypt|decrypt> <folder> <password>")
		return
	}

	action := os.Args[1]
	folder := os.Args[2]
	password := os.Args[3]

	key := make([]byte, 32)
	copy(key, password)

	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	fmt.Printf("Found %d valid files...\n", len(files))

	formatCount := make(map[string]int)
	processed := 0
	for _, path := range files {
		switch action {
		case "encrypt":
			if filepath.Ext(path) == ".enc" {
				continue
			}
			if err := encryptFile(path, key, formatCount); err != nil {
				fmt.Println("Error encrypting file:", err)
			}
			processed++
		case "decrypt":
			if filepath.Ext(path) == ".enc" {
				if err := decryptFile(path, key); err != nil {
					fmt.Println("Error decrypting file:", err)
				}
				processed++
			}
		}
	}

	fmt.Printf("Files processed: %d\n", processed)
	if action == "encrypt" {
		fmt.Printf("Encryption summary:\n")
		for format, count := range formatCount {
			fmt.Printf("  %s: %d\n", format, count)
		}
	}

}
