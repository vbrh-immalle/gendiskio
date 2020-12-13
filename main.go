package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	ticker := time.NewTicker(250 * time.Millisecond)
	done := make(chan (bool))

	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "tempfile-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Printf("Program killed! Removing [%s].", tmpfile.Name())
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
		// https://github.com/golang/go/issues/32300
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case tckr := <-ticker.C:
				fmt.Printf("At %s: %d\n", tckr.Format("20060102T150405"), rand.Intn(100))
				if _, err := tmpfile.Write(content); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Writing to %s\n", tmpfile.Name())
			}
		}
	}()

	fmt.Scanf("%d")
	fmt.Scanf("%d")
	log.Printf("Removing [%s].", tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(tmpfile.Name()); err != nil {
		log.Fatal(err)
	}

}
