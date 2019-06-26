package crawler

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// struct for properties we wanna keep on each url
type SearchUrlInfo struct {
	status  int
	domain  string
	referer string
	comment string
}

// Map for found links
type SearchUrlList map[string]*SearchUrlInfo

// Gets all url links in html returned from response
func SearchCrawl(domain string, searchTerm string) SearchUrlList {
	log.Println("Crawling")
	if !strings.Contains(domain, "http://") && !strings.Contains(domain, "https://") {
		domain = "http://" + domain
	}

	urlList := make(SearchUrlList)
	// Retrieve all links in homepage body
	links := AllLinks(domain)

	// Iterate through grabbed links and insert into url list if not already in there
	for _, linkString := range links {
		if len(linkString) < 2 {
			continue
		}
		// if relative url append to Domain before insertion
		if string(linkString[0]) == "/" {
			linkString = domain + linkString
		}
		// Finally check if the formatted linkString
		if _, ok := urlList[linkString]; !ok {
			urlList[linkString] = &SearchUrlInfo{1, domain, domain, ""}
		}
	}

	// Save length for urlList
	listLen := -1
	for len(urlList) != len(subSearchCrawl(urlList, domain, searchTerm)) {
		// Save length for this iteration
		listLen = len(urlList)
		log.Println("LIST LENGTH:", listLen)
	}
	// Return crawl list and crawl ID
	return urlList
}

// Search Crawl: looks for search term input without extra operations from "all crawl"
func subSearchCrawl(urlList SearchUrlList, domain string, searchTerm string) SearchUrlList {
	// iterate through unique map of urls
	for url, urlProps := range urlList {
		// check if url is in Domain and if status code has not been set
		if urlProps.status == 1 && strings.Contains(url, domain) {
			// GET request
			response, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}

			//check response content type
			ctype := response.Header.Get("Content-Type")
			if !strings.HasPrefix(ctype, "text/html") {
				err = response.Body.Close()
				if err != nil {
					log.Println(err)
				}
				// Then skip to next iteration without looking for links inside the picture, or song, or pdf, etc.
				continue
			}

			// SEARCH: if we have a searchTerm provided, do the thing
			if len(searchTerm) > 0 {
				// html response body so we can search
				responseData, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Println(err)
				}
				// stringify response
				responseString := string(responseData)
				foundMatch, err := regexp.MatchString(searchTerm, responseString)
				if err != nil {
					log.Println(err)
				}
				// update comment value if found
				if foundMatch {
					urlProps.comment = "Found Term: " + searchTerm
				}
			}

			landingUrl := response.Request.URL.String()
			// Get landing url; check if doesn't contain the domain then we've gone off track; skip iteration
			if !strings.Contains(landingUrl, domain) {
				log.Println("HEADER:", landingUrl)
			}

			// If Connection Successful, grab status code
			code := response.StatusCode
			// Set the status so it doesn't stay = 1 and keep getting crawled
			urlProps.status = code

			// *!* After we are done with current loop's url: we must manually close connection to avoid mem leak *!*
			err = response.Body.Close()
			if err != nil {
				log.Println(err)
			}
			// Skip to next iteration if error code
			if code >= 400 {
				continue
			}

			// Grab all links on current url page
			links := AllLinks(url)
			// Iterate through grabbed links and insert into url map list if not already in there
			for _, link := range links {
				// returns val if link already in urlList, found=false if not found
				_, found := urlList[link]
				if !found {
					// insert link into map if not found
					// check if link is long enough
					if len(link) > 2 {
						// Skip over iCal downloads
						if strings.Contains(link, "iCal") {
							continue
						}
						// if relative url append to Domain before insertion
						if string(link[0]) == "/" {
							link = domain + link
							_, found := urlList[link]
							if !found {
								urlList[link] = &SearchUrlInfo{1, domain, url, ""}
							} else {
								continue
							}
						} else if strings.Contains(link, domain) {
							urlList[link] = &SearchUrlInfo{1, domain, url, ""}
						}
					}
				}
			}
			// END getting new links

			// Rest system between single page requests
			time.Sleep(50 * time.Millisecond)
			// END if linkResp Successful
		} else if !strings.Contains(url, domain) {
			// If the url string doesn't have the Domain remove it from url list
			delete(urlList, url)
		}
		// END check if url is in Domain and status code
	}
	// return updated list with new urls
	return urlList
}

// Finds element tags in page (pass in the url, what you're looking for and the element tag type to look in)
func SearchTags(url string, searchTerm string, tagType string) []string {
	links := []string{}
	col := []string{}

	// If no protocol prepend it to url
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("ERROR: All Links Failed to crawl \"" + url + "\"")
		log.Println(err)
		return links
	}
	defer resp.Body.Close() // close Body when the function returns

	page := html.NewTokenizer(resp.Body)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == tagType {
			for _, attr := range token.Attr {

				if attr.Key == "href" {
					log.Println(attr.Val)
					if strings.Contains(attr.Val, searchTerm) {
						tl := removeHash(attr.Val)
						col = append(col, tl)
						resolve(&links, col)
					}
				}
			}
		}
	}
}