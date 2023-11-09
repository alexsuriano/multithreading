package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type APICEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	StatusText string `json:"statusText"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chViaCEP := make(chan string)
	chAPICEP := make(chan string)

	cep := 88514670

	urlViaCEP := fmt.Sprintf("https://viacep.com.br/ws/%d/json/", cep)
	urlAPICEP := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%d.json", cep)

	go GetCEP(ctx, chViaCEP, urlViaCEP)
	go GetCEP(ctx, chAPICEP, urlAPICEP)

	select {
	case data := <-chViaCEP:
		fmt.Printf("Dados de ViaCEP: %v\n", data)
	case data := <-chAPICEP:
		fmt.Printf("Dados de APICEP: %v\n", data)
	case <-time.After(1 * time.Second):
		cancel()
		fmt.Println("Timeout")
	}
}

func GetCEP(ctx context.Context, ch chan string, url string) {

	select {
	case <-ctx.Done():
		return

	default:
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			log.Println(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		ch <- string(body)
	}
}
