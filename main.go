package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func main() {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	f, err := os.Open("keys.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys) // decodes the string
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(keys.Key, keys.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	var token oauth2.Token
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		fmt.Println(err)
		return
	}
	var conf oauth2.Config
	TwitterClient := conf.Client(context.Background(), &token)
	response, err := TwitterClient.Get("https://api.twitter.com/1.1/statuses/retweets/1171044354359726082.json")
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	io.Copy(os.Stdout, response.Body)
}
