package tripleDES

import (
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

func Encrypt(text, key, cypher string) (res string, err error) {

	triplekey := key + key + key
	plaintext := []byte(text)
	//Padding
	if mod := len(plaintext) % des.BlockSize; mod != 0 {
		add := des.BlockSize - mod
		for n := 0; n < add; n++ {
			plaintext = append(plaintext, []byte(" ")...)
		}
	}

	block, err := des.NewTripleDESCipher([]byte(triplekey))
	if err != nil {
		return
	}

	//fmt.Printf("%d bytes NewTripleDESCipher key with block size of %d bytes\n", len(triplekey), block.BlockSize)
	ciphertext := []byte(cypher)
	iv := ciphertext[:des.BlockSize]

	// encrypt
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plaintext))
	mode.CryptBlocks(encrypted, plaintext)
	res = string(encrypted)
	return

}
func Decrypt(text, key, cypher string) (res string, err error) {

	tripleKey := key + key + key
	plaintext := []byte(text)
	if mod := len(plaintext) % des.BlockSize; mod != 0 {
		add := des.BlockSize - mod
		for n := 0; n < add; n++ {
			plaintext = append(plaintext, []byte(" ")...)
		}
	}
	block, err := des.NewTripleDESCipher([]byte(tripleKey))
	if err != nil {
		fmt.Printf("%s \n", err.Error())
		return
	}
	fmt.Printf("%d bytes NewTripleDESCipher key with block size of %d bytes\n", len(tripleKey), block.BlockSize)
	ciphertext := []byte(cypher)
	iv := ciphertext[:des.BlockSize]
	encrypted := []byte(text)

	//decrypt
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(plaintext))
	decrypter.CryptBlocks(decrypted, encrypted)
	fmt.Printf("%x decrypt to %s\n", encrypted, decrypted)

	res = string(decrypted)
	return

}
