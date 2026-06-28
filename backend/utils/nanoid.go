package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePublicID() (string, error) {
	id, err := gonanoid.Generate(alphabet, 8)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GenerateCustomId(size int) (string, error) {
	id, err := gonanoid.Generate(alphabet, size)
	if err != nil {
		return "", err
	}
	return id, nil
}
