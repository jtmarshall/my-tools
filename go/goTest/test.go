package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type UrlInfo struct {
	status  int
	domain  string
	referer string
}

func main() {
	var urlProps UrlInfo
	url := "www.test.com"
	domain := "www.test.com"

	for i := 1; i <= 5; i++ {
		log.Println(i)
		if i > 0 {
			urlErr := ""
			redirects := 0

			// Start request
			preResp, respErr := http.NewRequest("GET", url, nil)
			if respErr != nil {
				log.Println("PreResp Err, Continuing: ", respErr)
				urlErr = "Pre No Response"
				fmt.Println(urlErr)
			}

			// Need to separate out into transport request since we want redirects
			linkResp, linkErr := http.DefaultTransport.RoundTrip(preResp)
			if linkErr != nil {
				log.Println("LinkResp Err, Continuing: ", linkErr)
				urlErr = "No Response"
				fmt.Println(urlErr)
			} else {
				// If Connection Successful
				code := linkResp.StatusCode
				fmt.Println(code, linkResp.Request.URL.String(), linkResp.Header.Get("Location"))

				// Set the status so it doesn't stay one and keep getting crawled
				urlProps.status = code

				// *!* After we are done with current loop's url: we must manually close connection to avoid mem leak *!*
				linkResp.Body.Close()

				if code >= 400 {
					fmt.Println("4xx Status:", url, code)
					fmt.Println("Response: ", linkResp)
					fmt.Println("Response Body: ", linkResp.Status)
					linkResp.Body.Close()
					continue
				}

				nextURL := url
				// Check Redirects if 300 status, up to 10
				if 299 < code && code < 400 {
					for ; redirects <= 10; redirects++ {
						// Reset Pre response to next url
						preResp, respErr := http.NewRequest("GET", nextURL, nil)
						if respErr != nil {
							log.Println("PreResp Redirect Err: ", respErr)
							// Skip to next url if error
							break
						}
						// Next request
						resp, err := http.DefaultTransport.RoundTrip(preResp)
						if err != nil {
							log.Println("Redirect Break: ", err)
							urlErr = "Redirect Break: " + nextURL
							code = 300
							fmt.Println(domain, url, code, urlProps.referer)
							// Skip to next url if error
							break
						}

						if resp.StatusCode == 200 {
							code = resp.StatusCode
							resp.Body.Close()
							break
						} else {
							// Set next url to check
							nextURL = resp.Header.Get("Location")
							// Ignore http to https redirects
							if strings.Trim(resp.Request.URL.String(), "http") == strings.Trim(nextURL, "https") {
								// fmt.Println(resp.Request.URL.String(), nextURL)
								redirects--
							}

							// Ignore redirects for trailing slash
							subNext := nextURL
							if last := len(subNext) - 1; last >= 0 && subNext[last] == '/' {
								// If nextURL up has a trailing slash remove it
								subNext = subNext[:last]
							}
							if subNext == resp.Request.URL.String() {
								// Then if the current url and the nextURL(w/o the slash) are the same, ignore redirect
								redirects--
							}
						}
						resp.Body.Close()
						if redirects > 10 || code > 399 {
							// Log in DB outages, then skip to next iteration
							fmt.Println(domain, url, code, urlProps.referer)
							break
						}
					}
				}
				// END Check redirects
			}
		}
	}
}
