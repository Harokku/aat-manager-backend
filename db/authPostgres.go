package db

import (
	"aat-manager/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

const (
	GsuiteToken = "gtoken"
)

type Token struct {
}

// encryptToken encrypts a token using AES encryption with a given key.
// It takes a plainText string and a key as input parameters.
// It returns the encrypted token as a string and an error if any.
func encryptToken(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesgcm.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

// decryptToken decrypts a token using AES encryption with a given key.
func decryptToken(cipherText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ct, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce, ct := ct[:12], ct[12:]
	plainText, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// SaveToken saves the provided token into the database after encrypting it using the AES encryption key.
// It requires a valid database connection obtained from the pgConnect function.
// The encryption key is fetched from the environment using the ReadEnvOrPanic function.
// The token is stored in the tokens table in the database with the name "gtoken".
// The function returns an error if there is any issue encrypting or saving the token.
// Note: The pgConnect and ReadEnvOrPanic functions must be properly implemented and available.
// The encryptToken function is used to encrypt the token.
// The GsuiteToken constant specifies the name under which the token is stored in the database.
func (t Token) SaveToken(token string) error {
	db := pgConnect()                               // Acquire db connection
	aesKey := utils.ReadEnvOrPanic(utils.AESSECRET) // Acquire aes secret from env
	aesByteKey, err := hex.DecodeString(aesKey)     // Decode string to byte
	if err != nil {
		return err
	}
	encryptedToke, err := encryptToken(token, aesByteKey) // Encrypt token
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO tokens(name, value) VALUES ($1, $2)", GsuiteToken, encryptedToke)
	if err != nil {
		return err
	}

	return nil
}

// GetToken retrieves the encrypted token from the database and decrypts it using the AES encryption key.
// It returns the decrypted token as a string and an error if any.
// The encryption key is fetched from the environment using the ReadEnvOrPanic function.
// It requires a valid database connection obtained from the pgConnect function.
// The token is stored in the tokens table in the database with the name "gtoken".
// The function log any errors encountered during the retrieval or decryption process.
func (t Token) GetToken() (string, error) {
	db := pgConnect()

	var encryptedToken string
	err := db.QueryRow("SELECT value FROM tokens where name = $1", GsuiteToken).Scan(&encryptedToken)
	if err != nil {
		log.Printf("Failed to select token from database: %v", err)
		return "", err
	}

	// Fetch encryption key from env
	aesKey := utils.ReadEnvOrPanic(utils.AESSECRET)
	aesByteKey, err := hex.DecodeString(aesKey) // Decode string to byte
	if err != nil {
		return "", err
	}

	// Decrypt the token and return it
	decryptedToken, err := decryptToken(encryptedToken, aesByteKey)
	if err != nil {
		log.Printf("Failed to decrypt token: %v", err)
		return "", err
	}

	return decryptedToken, nil
}
