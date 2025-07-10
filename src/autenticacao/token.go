// go get github.com/dgrijalva/jwt-go
package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usário
func CriarToken(usuarioID uint64) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 16).Unix() //exp se refere a expiração, logo tempo de duraçao do token
	permissoes["usuarioId"] = usuarioID
	//até aqui foram as permissoes do nosso token, agora temos que gerar ele e depois assinar ele digitalmente
	//para tal vamos usar a chave secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes) //SigningMethodHS256 é o metodo de assinatura
	return token.SignedString([]byte(config.SecretKey))            //secret , deve ser seguro isso aqui para evitar vasamento de dados
}

// ValidarToken verifica se o token passado na requisição é valido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retronarChaveDeVerificacao) //jwt = jason web token
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { //vai retornar todos os maps criados (permissoes),
		return nil //ignoramos o primeiro retorno pois so queremos saber se eles existem
	}

	return errors.New("token inválido")
}

// ExtrairUsuarioID retorna o usuarioId que está salvo no token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r) // Chama uma função auxiliar para extrair o token JWT do cabeçalho da requisição HTTP
	// Faz o parsing (análise/validação) do token JWT
	// A função 'retronarChaveDeVerificacao' será usada para verificar a assinatura do token
	token, erro := jwt.Parse(tokenString, retronarChaveDeVerificacao) //jwt = jason web token
	if erro != nil {
		return 0, erro
	}

	// Se o token foi validado corretamente e seus dados (claims) podem ser acessados,
	// convertemos as claims para o tipo jwt.MapClaims (um map[string]interface{})
	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Aqui acessamos a claim "usuarioId" que está dentro do token
		// Como os claims do JWT são do tipo genérico interface{}, o valor vem como float64
		// Então usamos fmt.Sprintf para transformar esse float64 em uma string sem casas decimais (%.0f)
		// Depois usamos strconv.ParseUint para converter essa string para uint64
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64) //tivemos que transformar em string
		if erro != nil {
			return 0, erro
		}
		// Se ocorreu tudo bem até aqui retornamos o usuarioID extraido do token
		return usuarioID, nil

	}

	// Se o token não é válido ou não conseguimos acessar as claims corretamente
	return 0, errors.New("token inválido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization") //header = cabeçalho, entao o programa pega o valor do cabeçalho http

	if len(strings.Split(token, " ")) == 2 { // deve vir Bearer e mais meu token
		return strings.Split(token, " ")[1]
	}
	return ""
}

func retronarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	//vamos efetuar a verificação do método de assinatura agora, usamos o SigningMethodHS256, ele deve estar dentro de: SigningMethodHMAC
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}
	//então se nao der erro, a gente pode seguramente retornar a nossa chave secreta pelo return config.SecretKey, nil
	return config.SecretKey, nil
}
