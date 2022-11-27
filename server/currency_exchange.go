package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const exchangeURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type DollarExchange struct {
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

func GetExchange(db *sql.DB) *DollarExchange {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	exchange, err := getExchange(ctx)
	if err != nil {
		panic(err)
	}
	exchange.ID = uuid.New().String()

	err = saveExchange(db, exchange)
	if err != nil {
		panic(err)
	}

	return exchange
}

func getExchange(ctx context.Context) (*DollarExchange, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", exchangeURL, nil)
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

	var exchange DollarExchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return nil, err
	}

	return &exchange, nil
}

func saveExchange(db *sql.DB, exchange *DollarExchange) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, `INSERT INTO exchanges
	(id, code, codein, name, high, low, varbid, pctchange, bid, ask, timestamp, createdate) 
	values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(exchange.ID, exchange.Usdbrl.Code, exchange.Usdbrl.Codein, exchange.Usdbrl.Name,
		exchange.Usdbrl.High, exchange.Usdbrl.Low, exchange.Usdbrl.VarBid, exchange.Usdbrl.PctChange,
		exchange.Usdbrl.Bid, exchange.Usdbrl.Ask, exchange.Usdbrl.Timestamp, exchange.Usdbrl.CreateDate)

	if err != nil {
		return err
	}
	return nil
}

func GetAllExchangesSaved(db *sql.DB) []DollarExchange {
	rows, err := db.Query(`select id, code, name from exchanges`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var exchanges []DollarExchange

	for rows.Next() {
		var exchange DollarExchange

		err = rows.Scan(&exchange.ID, &exchange.Usdbrl.Code, &exchange.Usdbrl.Name)
		if err != nil {
			panic(err)
		}
		exchanges = append(exchanges, exchange)
	}
	return exchanges
}
