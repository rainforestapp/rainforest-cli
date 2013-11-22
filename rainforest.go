package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Issuing Run")
	client := new(http.Client)

	req, err := http.NewRequest("POST", "https://app.rainforestqa.com/api/1/runs?tests=all&conflict=abort", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "RainforestCli/1.0 (https://docs.rainforestqa.com/bot.html)")
	req.Header.Add("CLIENT_TOKEN", "<key>")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	j, err := simplejson.NewJson(response)
	if err != nil {
		log.Fatal(err)
	}

	run_id, err := j.Get("id").Int()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" - started #%d\n", run_id)
}
