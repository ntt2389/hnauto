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
    "net/smtp"
    "golang.org/x/net/html"
    "crypto/tls"
)

func main() {
  //if the caller didn't provide a URL to fetch...
  if len(os.Args) < 2 {
      //print the usage and exit with an error
      fmt.Printf("usage:\n  pagetitle <url>\n")
      os.Exit(1)
  }

  URL := os.Args[1]

  //top := getTop(URL)

  //fmt.Printf("https://news.ycombinator.com/item?id=%s\n", top)

  top := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", getTop(URL))

  title := getTitle(top)

  fmt.Println(title, top)

  sendMail(title, top)

}

func getTop(URL string) string {

    //Get the URL using func getResponse
    resp := getResponse(URL)

    //make sure the response body gets closed
    defer resp.Close()

    TokenPoint := make(map[string]int)

    //create a new tokenizer over the response body
    tokenizer := html.NewTokenizer(resp)

    //loop until we find the title element and its content
    //or encounter an error (which includes the end of the stream)
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
                        //take the id from "score_id"
                        id := tokenText[30:38]
                        //get Tokentype of next token
                        tokenNext := tokenizer.Next()
                        //if it is texttoken
                        if tokenNext == html.TextToken {
                              //take that
                              pointSl := strings.Split(string(tokenizer.Raw()), " ")
                              //add to the map id:point
                              TokenPoint[id], _ = strconv.Atoi(pointSl[0])
                        }
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
      return top


}

func getTitle(URL string) string {

    //Get the URL using func getResponse
    resp := getResponse(URL)

    //make sure the response body gets closed
    defer resp.Close()

    var title string

    //create a new tokenizer over the response body
    tokenizer := html.NewTokenizer(resp)

    //loop until we find the title element and its content
    //or encounter an error (which includes the end of the stream)
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
            if 0x1 == token.DataAtom {
  // returns the unmodified string of the current token
                  tokenText := string(tokenizer.Raw())
                  if strings.Contains(tokenText, "storylink") {
                        tokenNext := tokenizer.Next()
                              if tokenNext == html.TextToken {
                              title = fmt.Sprintln(string(tokenizer.Raw()))
                              }
                  }
            }
        }
    }
    return title

}


func getResponse(URL string) io.ReadCloser {
    resp, err := http.Get(URL)

    //if there was an error, report it and exit
    if err != nil {
        //.Fatalf() prints the error and exits the process
        log.Fatalf("error fetching URL: %v\n", err)
    }


    //check response status code
    if resp.StatusCode != http.StatusOK {
        log.Fatalf("response status code was %d\n", resp.StatusCode)
    }

    return resp.Body

}

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func sendMail(title string, top string) {
	mail := Mail{}
	mail.senderId = "trieungant89@gmail.com"
	mail.toIds = []string{"ntt2389@gmail.com", "jeannie2389@gmail.com"}
	mail.subject = "HN Today"
	mail.body = title + " " + top

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, "23890019", smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}


//Convert between types
