package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
)

type BadDate struct {
	Date, Mailer string
}

var (
	maildirGlob = flag.String("maildir",
		os.ExpandEnv("${HOME}/Maildir/.*/*/*"),
		"Pattern for email sent files")
	csvFn = flag.String("output", "mailers.csv", "Output for bad dates/mailers")
)

func main() {
	stats := map[string]int{}
	badDates := []BadDate{}

	flag.Parse()

	files, err := filepath.Glob(*maildirGlob)
	if err != nil {
		log.Fatalf("Failed to glob %q: %s", *maildirGlob, err)
	}

	for idx, fn := range files {
		if idx != 0 && idx % 1000 == 0 {
			log.Printf("Processing %d/%d", idx, len(files))
		}
		r, err := os.Open(fn)
		if err != nil {
			stats["failed-open"]++
			continue
		} else {
			stats["success-open"]++
		}

		msg, err := mail.ReadMessage(r)
		if err != nil {
			stats["failed-mail-parse"]++
			r.Close()
			continue
		} else {
			stats["success-mail-parse"]++
		}

		_, err = msg.Header.Date()
		if err != nil {
			stats["failed-date-parse"]++
			mailer := msg.Header.Get("X-Mailer")
			if mailer == "" {
				mailer = fn
			}
			badDates = append(badDates, BadDate{
				Date: msg.Header.Get("Date"),
				Mailer: mailer,
			})
			continue
		}
	}
	fmt.Println("Stats", stats)

	if count, ok := stats["failed-date-parse"]; ok || count != 0 {
		log.Println("Saving bad dates to ", *csvFn)
		f, err := os.OpenFile(*csvFn, os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Failed to create csv file %q: %s", *csvFn, err)
		}
		defer f.Close()
		csvf := csv.NewWriter(f)
		for _, bm := range badDates {
			err = csvf.Write([]string{
				bm.Mailer,
				bm.Date,
			})
			if err != nil {
				log.Fatalf("Failed to write record %#v", bm)
			}
		}
	} else {
		log.Println("No bad dates found")
	}
}
