package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func encryptFile(filePath string, key []byte) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(filePath + ".enc")
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

	outFile, err := os.Create(filePath[:len(filePath)-4])
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

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			switch action {
			case "encrypt":
				if filepath.Ext(path) == ".enc" {
					return nil
				}
				return encryptFile(path, key)
			case "decrypt":
				if filepath.Ext(path) == ".enc" {
					return decryptFile(path, key)
				}
				return nil
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Operation completed successfully.")
	}
}
