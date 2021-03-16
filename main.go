package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	token   = flag.String("t", "", "Your User Token (a.k.a token_v2)")
	pageid  = flag.String("p", "", "Your Page ID")
	spaceid = flag.String("s", "", "Your Space ID (Optional)")
	debug   = flag.Bool("v", false, "Verbose Mode")
)

func PrintStruct(emp interface{}) {
	empJSON, err := json.MarshalIndent(emp, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
}

func main() {
	flag.Parse()
	files := flag.Args()

	if *debug {
		log.Printf("acPasstoken = %s", *token)
		log.Printf("page_id = %s", *pageid)
		log.Printf("space_id = %s", *spaceid)
		log.Printf("verbose = true")
		log.Printf("files = %s", files)
	}

	if *token == "" || *pageid == "" {
		fmt.Println("token or pageid is missing")
		printUsage()
		return
	}

	for _, v := range files {
		fmt.Printf("Local: %s\n", v)
		link := upload(v)
		fmt.Printf("Download Link: %s", link)
	}
}

func printUsage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
