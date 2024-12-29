package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var (
	mode        = flag.String("m", "encrypt", "encrypt or decrypt")
	src         = flag.String("s", "", "source folder")
	destination = flag.String("d", "", "destination folder")
	pass        = flag.String("p", "", "encryption password")
)

type fileInfo struct {
	path string
	info os.FileInfo
}

func main() {

	flag.Parse()

	// Verify input parameters
	if *mode != "encrypt" && *mode != "decrypt" {
		fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'.")
		return
	}

	// Check if source and password are provided
	if *src == "" || *pass == "" {
		fmt.Println("Missing required parameters.")
		return
	}

	// Set destination folder to source folder if not provided
	if *destination == "" {
		*destination = filepath.Join(*src, "encrypted")
	}

	// Check if source folder exists
	if _, err := os.Stat(*src); os.IsNotExist(err) {
		fmt.Println("Source folder does not exist.")
		return
	}
	// Check if destination folder exists
	if _, err := os.Stat(*destination); os.IsNotExist(err) {
		fmt.Println("Destination folder does not exist.")
		if err := os.MkdirAll(*destination, 0755); err != nil {
			fmt.Println("Error creating destination folder:", err)
			return
		}
	}

	key := make([]byte, 32)
	copy(key, *pass)

	var files []fileInfo

	err := filepath.Walk(*src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, fileInfo{path, info})
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sort.Slice(files, func(i, j int) bool { // Sort files by modification time (oldest first)
		return files[i].info.ModTime().Before(files[j].info.ModTime())
	})

	fmt.Printf("Found %d valid files...\n", len(files))
	for _, path := range files {
		fmt.Println(path)
	}

	formatCount := make(map[string]map[string]int)
	processed := 0
	for _, file := range files {
		switch *mode {
		case "encrypt":
			if filepath.Ext(file.path) == ".enc" {
				continue
			}
			if err := encryptFile(file.path, *src, *destination, key, formatCount); err != nil {
				fmt.Println("Error encrypting file:", err)
			}
			processed++
		case "decrypt":
			if filepath.Ext(file.path) == ".enc" {
				if err := decryptFile(file.path, *src, *destination, key); err != nil {
					fmt.Println("Error decrypting file:", err)
				}
				processed++
			}
		}
	}

	fmt.Printf("Files processed: %d\n", processed)
	if *mode == "encrypt" {
		fmt.Printf("Encryption summary:\n")
		for subfolder, formats := range formatCount {
			fmt.Printf("Subfolder: %s\n", subfolder)
			for format, count := range formats {
				fmt.Printf("  %s: %d\n", format, count)
			}
		}
	}
}

func encryptFile(filePath string, src string, destination string, key []byte, formatCount map[string]map[string]int) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file %s has no extension", filePath)
	}
	format := ext[1:] // Remove the dot from the extension
	if format == "jpeg" {
		format = "jpg"
	}

	// Create the destination subfolder structure
	relPath, err := filepath.Rel(src, filePath)
	if err != nil {
		return err
	}
	subfolder := filepath.Dir(relPath)
	destPath := filepath.Join(destination, subfolder)
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	// Initialize the format count map for the subfolder if not already initialized
	if _, exists := formatCount[subfolder]; !exists {
		formatCount[subfolder] = make(map[string]int)
	}
	formatCount[subfolder][format]++

	// get file creation date
	datePrefix := time.Now().UTC().Format("20060102")
	inFileInfo, err := os.Stat(filePath)
	if err == nil {
		datePrefix = inFileInfo.ModTime().UTC().Format("20060102")
	}
	//creationTime := inFileInfo.ModTime().UTC().Format("20060102")

	// Today's date in format YYYYMMDD
	//datePrefix := time.Now().UTC().Format("20060102")
	newFileName := fmt.Sprintf("%s_%s_%04d%s.enc", datePrefix, format, formatCount[subfolder][format], "."+format)

	outFile, err := os.Create(filepath.Join(destPath, newFileName))
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

func decryptFile(filePath string, src string, destination string, key []byte) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Extract the original file name without the .enc extension
	originalFileName := filepath.Base(filePath[:len(filePath)-4])
	// Add the dec_ prefix to the original file name
	// newFileName := "dec_" + originalFileName

	// Create the destination subfolder structure
	relPath, err := filepath.Rel(src, filePath)
	if err != nil {
		return err
	}
	subfolder := filepath.Dir(relPath)
	destPath := filepath.Join(destination, subfolder)
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	// Create the new file path with the dec_ prefix
	outFile, err := os.Create(filepath.Join(destPath, originalFileName))
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
