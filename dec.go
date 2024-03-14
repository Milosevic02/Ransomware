package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	folder := "pdfs"
	const aes_key = "12345678123456781234567812345678" //Your Key

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path: ", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc") {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file: ", path, err)
				return err
			}

			newFileName := strings.TrimSuffix(path, ".enc")
			file, err := os.Create(newFileName)
			if err != nil {
				fmt.Println("Error creating file: ", path, err)
				return err
			}
			defer file.Close()

			aes_key_byted := []byte(aes_key)

			aes_key_cipher, _ := aes.NewCipher(aes_key_byted)

			gcm, err := cipher.NewGCM(aes_key_cipher)
			if err != nil {
				fmt.Println("Some error")
				return err
			}

			nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

			plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				fmt.Println("Error opening file: ", path, err)
				return err
			}

			_, err = file.Write(plainText)
			if err != nil {
				fmt.Println("Error writing file: ", path, err)
				return err
			}

			err = os.Remove(path)
			if err != nil {
				fmt.Println("Error removing file: ", path, err)
				return err
			}

			fmt.Println("Decryption successfull")
		}
		return err
	})
	if err != nil {
		fmt.Println("Error walking through the files/folders", err)
		return
	}
}
