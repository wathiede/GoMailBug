package main

import (
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
)

type Contact struct {
	mail.Address
	fn string
}

var emailPat *string = flag.String("email_pat", "", "Glob pattern for email to parse")
var output *string = flag.String("output", "", "File to store output")

func init() {
	flag.Parse()
}

func main() {
	results := make(chan *Contact)

	go getEmails(results)

	save(results)
}

func getAddresses(fn string, results chan<- *Contact) {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Failed to open %q: %s\n", fn, err)
	}

	msg, err := mail.ReadMessage(file)
	if err != nil {
		log.Printf("Failed to parse message %q: %s\n", fn, err)
		return
	}

	headers := []string{"to", "cc"}
	for _, header := range headers {
		addrs, err := msg.Header.AddressList(header)
		if err != nil {
			//log.Printf("Failed to parse %s: header in %q: %s", header, fn, err)
			continue
		}

		for _, addr := range addrs {
			c := &Contact{*addr, fn}
			results <- c
		}
	}
}

func getEmails(results chan<- *Contact) {
	var files sort.StringSlice
	files, err := filepath.Glob(*emailPat)
	if err != nil {
		log.Fatalf("Failed to glob %q: %s\n", *emailPat, err)
	}
	sort.Sort(files)

	for _, fn := range files {
		getAddresses(fn, results)
	}
	close(results)
}


func save(results <-chan *Contact) {
	file, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", *output, err)
	}
	defer file.Close()

	for c := range results {
		fmt.Fprintf(file, "%v|%s|%s\n", c.Address.Name, c.Address.Address, c.fn)
	}
}
