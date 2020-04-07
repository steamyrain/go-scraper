package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func getHref(t html.Token) (ok bool, href string) {
	for _,a := range t.Attr  {
		ok = a.Key == "href"
		if ok {
			href = a.Val
		}
	}
	return
}

func  crawl(url string,ch chan string, chFinished chan bool) {
	resp,err := http.Get(url)
	defer func() {
		chFinished <- true
	}()
	if err != nil {
		fmt.Println("ERROR:",err)
	}	
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch {
			case tt == html.ErrorToken:
				return
			case tt == html.StartTagToken:
				t := z.Token()
				isAnchor := t.Data == "a"
				if !isAnchor {continue}
				ok, urls := getHref(t)
				if !ok {continue}
				nhasProto := strings.Index(urls,"http") != 0 
				if !nhasProto { ch <- urls }
		}
	}
}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]
	ch := make(chan string)
	chFinished := make(chan bool)
	for _, url : range seedsUrls {
		go crawl(url,ch,chFinished)
	}	
	for c:=0; c<len(seedUrls); {
		select {
			case url := <-ch:
				found[url] = true
			case <- chFinished:
				c++
		}
	}
	fmt.Println("\nFound", len(foundUrls), "unique urls")
	for url,_ := range foundUrls {
		fmt.Println("-",url)
	}
	close(ch)
}
