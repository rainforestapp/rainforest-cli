package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const endpoint_url = string("https://app.rainforestqa.com/")

func print_usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("rainforest --run --token 12345 --tags run-me --fg --fail-fast")
	fmt.Println("")
}

func exit_and_print_usage(message string) {
	fmt.Println("Error:", message, "\n")
	print_usage()
	os.Exit(1)
}

func main() {
	run := flag.Bool("run", false, "Start a new run")
	tests := flag.String("tests", "", "Comma seperated list of tests to run. Pass 'all' to run everything.")
	tags := flag.String("tags", "", "Run tests tagged with this")
	token := flag.String("token", "", "Your Rainforest API token (get it a test's API tab)")
	abort := flag.Bool("abort", false, "Abort any existing runs (helps to save steps, or if you have webhooks)")
	fg := flag.Bool("fg", false, "Run in the foreground, outputting progress every 5 seconds")
	fail_fast := flag.Bool("fail-fast", true, "Return a failure as soon as it happens")

	flag.Parse()

	if *token == "" {
		exit_and_print_usage("You must specify your API token with --token <your token>")
	}

	if *run != true {
		exit_and_print_usage("You must specify --run (it is the only thing supported at the moment)")
		os.Exit(1)
	}

	if *tests == "" && *tags == "" {
		exit_and_print_usage("You must specify either --tags <some comma seperated tags> or --tests <all or comma seperated list of ids> to run")
		os.Exit(1)
	}

	log.Println("Issuing run")
	client := new(http.Client)

	params := ""

	if *tests != "" {
		params += "tests=" + *tests
	} else {
		params += "tags=" + *tags
	}

	if *abort == true {
		params += "&conflict=abort"
	}

	req, err := http.NewRequest("POST", endpoint_url+"/api/1/runs?"+params, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "RainforestCli/1.0 (https://docs.rainforestqa.com/bot.html)")
	req.Header.Add("CLIENT_TOKEN", *token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if res.StatusCode != 201 {
		if res.StatusCode == 404 {
			log.Fatal("Error: We could not find that account - please check your token.")
		} else {
			log.Fatal("Error: Issuing run failed with code ", res.StatusCode)
		}
		os.Exit(1)
	}

	j, err := simplejson.NewJson(response)
	if err != nil {
		fmt.Printf("Error resp\n\n%s\n", response)
		log.Fatal(err)
		os.Exit(1)
	}

	run_id, err := j.Get("id").Int()

	if err != nil {
		fmt.Printf("Error resp\n\n%s\n", response)
		log.Fatal(err)
	} else {
		log.Println(fmt.Sprintf("Run %d started. See %sruns/%d for web progress.", run_id, endpoint_url, run_id))
	}

	if *fg == true {
		for true {
			time.Sleep(1000 * time.Millisecond)

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/1/runs/%d", endpoint_url, run_id), nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("User-Agent", "RainforestCli/1.0 (https://docs.rainforestqa.com/bot.html)")
			req.Header.Add("CLIENT_TOKEN", *token)

			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			response, err := ioutil.ReadAll(res.Body)
			res.Body.Close()

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			if res.StatusCode != 200 {
				log.Println(fmt.Sprintf("P error: %s", response))
			} else {
				j, err := simplejson.NewJson(response)
				if err == nil {
					current_progress, current_progress_err := j.Get("current_progress").Map()
					state, state_err := j.Get("state").String()
					result, result_err := j.Get("result").String()

					if state_err == nil && result_err == nil && current_progress_err == nil {
						if state == "in_progress" {
							log.Printf("Run %d is %s and has %s. Current progress : %s%%", run_id, state, result, current_progress["percent"])
						} else {
							log.Printf("Run %d is %s.", run_id, state)
						}

						if *fail_fast == true && result == "failed" {
							log.Fatalln("Run failed")
							os.Exit(1)
						} else if state == "timed_out" || state == "aborted" {
							log.Fatalln("Run", state)
							os.Exit(1)
						} else if state == "complete" {
							log.Println("Run complete and was", result)
							os.Exit(0)
						}
					}
				}
			}
		}
	}

	os.Exit(0)
}
