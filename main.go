package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type apiR struct {
	Response Response
}

type coin struct {
	USD float64 `json:"דUSD"`
}

type Response struct {
	Ada coin `json:"ADA"`
}

type JSONType struct {
	CoinMap map[string]ValsType `json:"data"`
}

type ValsType struct {
	Id        string `json:"Id"`
	CoinName  string `json."CoinName`
	Url       string `json:"Url"`
	ImageUrl  string `json:"ImageUrl"`
	Name      string `json:"Name"`
	Symbol    string `json:"Symbol"`
	FullName  string `json:"FullName"`
	Algorithm string `json:"Algorithm"`
	ProofType string `json:"ProofType"`
	ToUSD     float64
	PriceTime time.Time
}

func initIMDB() *JSONType {
	// HTTP call
	resp, err := http.Get("https://min-api.cryptocompare.com/data/all/coinlist")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Decode JSON
	dec := json.NewDecoder(resp.Body)
	data := &JSONType{}
	if err := dec.Decode(data); err != nil {
		fmt.Println(err)
	}

	return data
}

func getCoinPrice() float64 {
	// HTTP call
	resp, err := http.Get("https://min-api.cryptocompare.com/data/pricemulti?fsyms=ADA&tsyms=USD")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Decode JSON
	dec := json.NewDecoder(resp.Body)
	data := &apiR{}
	if err := dec.Decode(data); err != nil {
		fmt.Println(err)
	}

	fmt.Println(data.Response.Ada.USD)
	return data.Response.Ada.USD
}

func main() {

	imdb := initIMDB()
	ttl := 1

	fmt.Println(imdb.CoinMap["ADA"].FullName)

	todayWith := time.Now().Add(time.Hour * time.Duration(ttl))
	if todayWith.After(imdb.CoinMap["ADA"].PriceTime) {
		fmt.Println("need to updated price")
		currentPrice := getCoinPrice()
		fmt.Println(currentPrice)

	} else {
		fmt.Println(imdb.CoinMap["ADA"].PriceTime)
		fmt.Println(imdb.CoinMap["ADA"].ToUSD)
	}
}
