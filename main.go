package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

// RetweetResponse struct
type RetweetResponse struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

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
	usernames, err := retweeters(TwitterClient, "1171044354359726082")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(usernames)
}

func retweeters(client *http.Client, tweetID string) ([]string, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", tweetID)
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var retweets []RetweetResponse
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&retweets)
	if err != nil {
		return nil, err
	}
	var usernames []string
	for _, retweet := range retweets {
		usernames = append(usernames, retweet.User.ScreenName)
	}
	return usernames, nil
}
