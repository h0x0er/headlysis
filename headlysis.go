package headlysis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	singleUrl uint = iota
	multUrl
	preConError
)

const (
	tag string = "[Headlysis]:"
)

var mainOutput Output // used online for mutli urls.
var mutex = new(sync.Mutex)
var counter = 1

func preCondition(opts *Options) uint {

	if opts.Url == "" && opts.UrlFile != "" && opts.OutputFile != "" {
		return multUrl
	}

	if opts.Url == "" && opts.UrlFile == "" && opts.OutputFile == "" {
		return preConError
	}

	if opts.Url == "" && opts.UrlFile != "" && opts.OutputFile == "" {
		return preConError
	}

	if opts.Url == "" && opts.UrlFile == "" && opts.OutputFile != "" {
		return multUrl
	}

	if opts.Url != "" {
		return singleUrl
	}

	return preConError

}

func makeRequest(url string) (*http.Response, error) {

	client := http.DefaultClient

	return client.Get(url)
}

func Headlysis(opts *Options) {
	switch preCondition(opts) {

	case singleUrl:
		// fmt.Println(tag, "Welcome to headlysis... We are analyzing your url... please wait for some time...")
		resp, err := makeRequest(opts.Url)
		if err != nil {
			log.Fatal(tag + " Error while making request")
		}
		np, p := GetMissingHeaders(resp.Header)

		var out Output

		tempOut := UrlAnalysisOutput{}
		tempOut.Url = opts.Url

		for _, sh := range p {
			tempOut.PresentHeaders = append(tempOut.PresentHeaders, PresentHeader{Name: sh.Name, Value: resp.Header.Get(sh.Name)})
		}

		for _, h := range np {
			tempOut.NotPresentHeaders = append(tempOut.NotPresentHeaders, NotPresentHeader{Name: h.Name, InfoUrl: h.GetUrl()})
		}
		out.Headlysis = append(out.Headlysis, tempOut)

		bytes, _ := json.Marshal(out)
		fmt.Println(string(bytes))

	case multUrl:

		urls, _ := getUrl(opts.UrlFile)

		max_jobs := opts.Threads
		done := make(chan bool, max_jobs)

		for i, url := range urls {
			if i%20 == 0 {
				time.Sleep(time.Second * 2)
			}
			go handleUrl(url, opts.Verbose, done)
		}

		for status := 0; status < len(urls); status++ {
			<-done
		}

		bytes, _ := json.Marshal(mainOutput)

		ioutil.WriteFile(opts.OutputFile, bytes, 0666)

		fmt.Println(string(bytes))

	case preConError:
		log.Fatal("Please check supplied options. Pass --help for viewing help")

	}

}

func handleUrl(url string, verbose bool, done chan bool) {

	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	resp, err := makeRequest(url)

	if err != nil {
		if verbose {
			fmt.Println("[", err.Error(), "]", url)
		}

	}

	if resp != nil {
		mutex.Lock()
		defer mutex.Unlock()

		fmt.Println(counter, ": ", url)
		counter++

		np, p := GetMissingHeaders(resp.Header)

		var tempOut UrlAnalysisOutput
		tempOut.Url = url

		for _, sh := range p {
			tempOut.PresentHeaders = append(tempOut.PresentHeaders, PresentHeader{Name: sh.Name, Value: resp.Header.Get(sh.Name)})
		}

		for _, h := range np {
			tempOut.NotPresentHeaders = append(tempOut.NotPresentHeaders, NotPresentHeader{Name: h.Name, InfoUrl: h.GetUrl()})
		}

		mainOutput.Headlysis = append(mainOutput.Headlysis, tempOut)
	}
	done <- true

}

func getUrl(file string) ([]string, error) {

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	return lines, nil
}
