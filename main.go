package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	Coins map[string]coin `json:""`
}

type coin struct {
	USD float64 `json:"USD"`
}

type JSONType struct {
	CoinMap map[string]ValsType `json:"Data"`
}

type ValsType struct {
	Id        string `json:"Id"`
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

var imdb = initIMDB()
var ttl = 1

func getCoinPrice(c string) float64 {
	// HTTP call
	params := url.Values{
		"fsyms": {c},
		"tsyms": {"USD"},
	}
	reqUrl := "https://min-api.cryptocompare.com/data/pricemulti?" + params.Encode()
	fmt.Println(reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	mp := make(map[string]coin)

	// Decode JSON into our map
	bytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &mp)
	if err != nil {
		println(err)
	}

	return mp[c].USD
}

func getReqCoinPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Gophers!")
	key := strings.ToUpper(r.URL.Path[10:]) // Get Coin symbol
	if _, ok := imdb.CoinMap[key]; ok {
		fmt.Fprintln(w, "You Asked for : "+key)
		timeWithTTL := imdb.CoinMap["ADA"].PriceTime.Add(time.Minute * time.Duration(ttl))
		if time.Now().After(timeWithTTL) {
			fmt.Fprintln(w, "need to updated price")
			t := imdb.CoinMap["ADA"]
			currentPrice := getCoinPrice("ADA")
			t.ToUSD = currentPrice
			t.PriceTime = time.Now()
			imdb.CoinMap["ADA"] = t

			fmt.Fprintln(w, fmt.Sprintf("%f", imdb.CoinMap["ADA"].ToUSD))
		} else {
			fmt.Fprintln(w, fmt.Sprintf("%f", imdb.CoinMap["ADA"].ToUSD))
		}
	} else {
		fmt.Fprintln(w, "You Asked for : "+key+" But there is no Coin like this")
	}

}

func main() {
	http.HandleFunc("/getPrice/", getReqCoinPrice)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
