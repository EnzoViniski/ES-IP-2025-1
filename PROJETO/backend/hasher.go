package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Altere a senha abaixo para a senha que você deseja hashear
	senha := "senha_segura_123"

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erro ao gerar hash:", err)
		return
	}

	fmt.Println("--- COPIE A LINHA DE HASH ABAIXO ---")
	fmt.Println(string(hash))
	fmt.Println("------------------------------------")
}
