package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// TryDecryptOrReturnPlainText Пытаемся расшифровать данные с использованием ключа или возвращаем PlainText
func TryDecryptOrReturnPlainText(keyFilename string, data []byte) []byte {
	dec, err := usingKeyFile(keyFilename, data)
	if err != nil {
		return data
	}

	return dec
}

// usingKeyFile Вызываем decrypt с приватным ключом
func usingKeyFile(keyFilename string, data []byte) ([]byte, error) {
	key, err := getPrivateKeyFromFile(keyFilename)
	if err != nil {
		return nil, err
	}

	return decrypt(key, data)
}

// decrypt расшифровывает открытый текст с помощью RSA и схемы заполнения из PKCS #1 v1.5
func decrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}

// getPrivateKeyFromFile Получается RSA приватный ключ из файла
func getPrivateKeyFromFile(keyFilename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(keyFilename)
	if err != nil {
		return nil, err
	}

	spkiBlock, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
