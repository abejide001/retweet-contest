package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
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
	var keyFile, usersFile, tweetID string
	flag.StringVar(&keyFile, "key", "keys.json", "the file where you store consumer key and secret")
	flag.StringVar(&usersFile, "users", ".users.csv", "the file where users who retweeted the tweet are stored")
	flag.StringVar(&tweetID, "tweet", "1171681896511758336", "the id of the tweet")
	flag.Parse()
	key, secret, err := keys(keyFile)
	if err != nil {
		fmt.Println(err)
	}
	client, err := TwitterClient(key, secret)
	if err != nil {
		fmt.Println(err)
	}
	usernames, err := retweeters(client, tweetID)
	if err != nil {
		fmt.Println(err)
	}
	existingUsers := ExistingUsers(usersFile)
	allUsernames := merge(usernames, existingUsers)
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, username := range allUsernames {
		if err := w.Write([]string{username}); err != nil {
			fmt.Println(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println(err)
	}
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

func keys(keyFile string) (key, secret string, err error) {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	f, err := os.Open(keyFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys) // decodes the string
	return keys.Key, keys.Secret, nil
}

// TwitterClient function, takes in a string.
func TwitterClient(key, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}
	var conf oauth2.Config
	TwitterClient := conf.Client(context.Background(), &token)
	return TwitterClient, nil
}

// ExistingUsers function
func ExistingUsers(userFile string) []string {
	f, err := os.Open(userFile)
	if err != nil {
		return nil
	}
	defer f.Close()
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}
	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, user := range a {
		uniq[user] = struct{}{}
	}

	for _, user := range b {
		uniq[user] = struct{}{}
	}
	var ret = make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}
	return ret
}
