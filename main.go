package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Carregar()
	r := router.Gerar()

	fmt.Println(config.StringConexaoBanco)
	fmt.Println("Escutando na Porta:", config.Porta)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r)) //aqui o Sprintf, esta fazendo o papel do (:5000)
}

// func init() {   		AQUI FOI USADO PARA CRIAR A SECRET KEY, APENAS 1 vez
// 	chave := make([]byte, 64) //cria nosso slice de 64

// 	if _, erro := rand.Read(chave); erro != nil { //isso vai popular nosso slice de 64
// 		log.Fatal(erro)
// 	}

// 	stringBase64 := base64.StdEncoding.EncodeToString(chave)
// 	fmt.Println(stringBase64)
// }
