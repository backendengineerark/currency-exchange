package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

const serverURL = "http://localhost:8080/cotacao"

type DollarExchangeInputDTO struct {
	ID     string
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	exchange, err := getExchange()
	if err != nil {
		panic(err)
	}

	err = saveBid(*exchange)
	if err != nil {
		panic(err)
	}

}

func getExchange() (*DollarExchangeInputDTO, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var exchange DollarExchangeInputDTO
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return nil, err
	}

	return &exchange, nil
}

func saveBid(exchange DollarExchangeInputDTO) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("Dólar: " + exchange.Usdbrl.Bid + "\n")
	if err != nil {
		return err
	}

	return nil
}
