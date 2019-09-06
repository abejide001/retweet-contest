package main

import (
	"fmt"
	"os"
	"encoding/json"
)

func main() {
	var keys struct {
		Key     string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	f, err := os.Open("keys.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys) //decodes the string
	fmt.Println(keys)
}
