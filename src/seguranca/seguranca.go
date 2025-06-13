// go get golang.org/x/crypto/bcrypt
package seguranca

import "golang.org/x/crypto/bcrypt"

// Hash recebe um string e coloca um hash nela
func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	//custo da operação, enquanto maior mais dificil de quebrar a operação (consome recursos comptacionais)
}

// VerificarSenha compara um senha e um hash e retorna se elas são iguais
func VerificarSenha(senhaComHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaComHash), []byte(senhaString))
}
