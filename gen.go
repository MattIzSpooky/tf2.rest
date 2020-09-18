// +build exclude

package main

import (
	"fmt"
	classPkg "github.com/MattIzSpooky/tf2.rest/class"
	"github.com/MattIzSpooky/tf2.rest/codegen"
	"os"
	"sync"
	"time"
)



func main() {
	var wg sync.WaitGroup
	errChan := make(chan error)


	for _, class := range classPkg.Classes {
		wg.Add(1)
		go runSingle(class, &wg, errChan)
	}

	go func() {
		for err := range errChan {
			fmt.Println(err.Error())
		}
	}()

	wg.Wait()
	close(errChan)
}

func runSingle(class string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	scraper := codegen.NewScraper(class)
	err := scraper.FetchDocument()

	if err != nil {
		errChan <- err
		return
	}

	rspSlice, err := scraper.Run()

	if err != nil {
		errChan <- err
		return
	}

	tmpl := codegen.NewResponseTemplate()
	fileName := fmt.Sprintf("responses/%s.go", class)
	f, err := os.Create(fileName)

	if err != nil {
		errChan <- err
		return
	}

	err = tmpl.Execute(f, codegen.ResponseTemplateData{
		Timestamp: time.Now(),
		URL:       scraper.GetURL(),
		Responses: rspSlice,
		Class:     class,
	})

	if err != nil {
		errChan <- err
		return
	}

	fmt.Println(fmt.Sprintf("Generated %s", fileName))

	if f.Close() != nil {
		errChan <- err
	}
}
