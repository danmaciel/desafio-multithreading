package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
 Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
 As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://brasilapi.com.br/api/cep/v1/01153000 + cep

http://viacep.com.br/ws/" + cep + "/json/

Os requisitos para este desafio são:
- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

*/

func main() {

	fmt.Println("Digite o CEP para pesquisar")
	reader := bufio.NewReader(os.Stdin)
	zipcode, errorOnRead := reader.ReadString('\n')

	if errorOnRead != nil {
		print("Entre com uma cep válido!\n")
		return
	}
	cepReadyFromSearch := treatZipCode(zipcode)

	_, errConvertToInt := strconv.Atoi(cepReadyFromSearch)
	if errConvertToInt != nil {
		print("Entre com uma cep válido!\n")
		return
	}

	apiSearch(cepReadyFromSearch)
}

func apiSearch(zipcode string) {
	viacepChannel := make(chan string)
	brasilapiChannel := make(chan string)

	go getApiData("https://brasilapi.com.br/api/cep/v1/"+zipcode, viacepChannel)
	go getApiData("http://viacep.com.br/ws/"+zipcode+"/json/", brasilapiChannel)

	select {
	case msg1 := <-viacepChannel:
		print("\nViaCep => " + msg1 + "\n")
	case msg2 := <-brasilapiChannel:
		print("\nBrasilAPI => " + msg2 + "\n")
	case <-time.After(time.Second):
		print("TimeOut\n")
	}
}

func getApiData(url string, c chan<- string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Não foi possível criar a requisição: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição: %s\n", err)
	}

	resBody, err := io.ReadAll(res.Body)

	if err == nil {
		c <- string(resBody)
	}

}

func treatZipCode(zipcode string) string {
	aux := strings.TrimSpace(zipcode)
	aux = strings.Replace(aux, ".", "", -1)
	aux = strings.Replace(aux, "-", "", -1)

	return aux
}
