package main

import (
    "fmt"
    "os"
    "log"
    "io"
    "strings"
    "strconv"
    "sort"
    "net/http"
    "golang.org/x/net/html"
)

func main() {
    //if the caller didn't provide a URL to fetch...
    if len(os.Args) < 2 {
        //print the usage and exit with an error
        fmt.Printf("usage:\n  pagetitle <url>\n")
        os.Exit(1)
    }

    URL := os.Args[1]
    //GET the URL
    resp, err := http.Get(URL)

    //if there was an error, report it and exit
    if err != nil {
        //.Fatalf() prints the error and exits the process
        log.Fatalf("error fetching URL: %v\n", err)
    }

    //make sure the response body gets closed
    defer resp.Body.Close()

    //check response status code
    if resp.StatusCode != http.StatusOK {
        log.Fatalf("response status code was %d\n", resp.StatusCode)
    }

    //create a new tokenizer over the response body
    tokenizer := html.NewTokenizer(resp.Body)

    //loop until we find the title element and its content
    //or encounter an error (which includes the end of the stream)
    TokenPoint := make(map[string]int)
    for {
        //get the next token type
        tokenType := tokenizer.Next()

        //if it's an error token, we either reached
        //the end of the file, or the HTML was malformed
        if tokenType == html.ErrorToken {
            err := tokenizer.Err()
            if err == io.EOF {
                //end of the file, break out of the loop
                break
            }
            //otherwise, there was an error tokenizing,
            //which likely means the HTML was malformed.
            //since this is a simple command-line utility,
            //we can just use log.Fatalf() to report the error
            //and exit the process with a non-zero status code
            log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
        }
        //if this is a start tag token...
        if tokenType == html.StartTagToken {
            //get the token
            token := tokenizer.Token()
            //if the tag of the element is "span"
            if 0x9f04 == token.DataAtom {
                  // returns the unmodified string of the current token
                  tokenText := string(tokenizer.Raw())
                  if strings.Contains(tokenText, "score_") {
                        //tokenSplit := strings.Split(tokenText, "_")
                        id := tokenText[30:38]
                        //get Tokentype of next token
                        tokenNext := tokenizer.Next()
                        //if it is texttoken
                        if tokenNext == html.TextToken {
                              //take that
                              pointSl := strings.Split(string(tokenizer.Raw()), " ")
                              //fmt.Println(pointSl[0])
                              TokenPoint[id], _ = strconv.Atoi(pointSl[0])
                        }

                  //for _, tokenAttr := range token.Attr {
                      // if the second attribute contains "id score_"
                      //if strings.Contains(tokenAttr, "id score_") {
                        //report the element's attributes
                          //fmt.Println(tokenAttr)
                      //}
                  //}
                  }

        }

        }
      }
      var points []int
      for _, pt := range TokenPoint {
          points = append(points, pt)
        }
      sort.Sort(sort.Reverse(sort.IntSlice(points)))
      var top string
      for j, k := range TokenPoint {
        if k == points[0] {
            top = j
        }
      }
      fmt.Printf("https://news.ycombinator.com/item?id=%s\n", top)

}

//Convert between types
