package crawler

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	// Go get...
	"golang.org/x/net/html"
)

var skipMap = map[string]bool{
	"www.test.com": true,
}
var domainList []string
// struct for properties we wanna keep on each url
type UrlInfo struct {
	status  int
	domain  string
	referer string
}
// Map for found links
type UrlList map[string]*UrlInfo

// schema for urls/phoneNumbers
var phoneMap = map[string][]string{}

// Populate Domain list on init
func init() {
	if len(domainList) < 1 {
		getDomainList()
	}
	// Sort list alphabetically
	sort.Strings(domainList)
}

// RunAllCrawl starts the all crawl
func RunAllCrawl() {
	// Safe to start crawl
	log.Println("All Crawl Started.")
	// Reset Domain list
	getDomainList()
	// Start all crawl
	allCrawl(domainList)
	log.Println("All Crawl End.", "Writing to CSV")
}

// allCrawl Everything
func allCrawl(domains []string) {
	// Re-create file so we start with a clean slate
	wFile, err := os.Create("./PhoneCrawl.csv")
	if err != nil {
		log.Fatal("ERR Creating File:", err)
	}
	wFile.Close()

	// Channel buffer for concurrency
	chanBuff := make(chan bool, 8)
	// Prime queue
	for len(chanBuff) < 8 {
		chanBuff <- true
	}

	// START Crawling!
	for i := 0; i < len(domains); {

		// Check if the Domain is in the skip list
		if skipMap[domains[i]] {
			// If it is: then increment i and continue, so we can skip it
			i++
			continue
		}

		// Check if we have a val in buffer
		if <-chanBuff {
			go func(i int) {
				log.Println("START", i, domains[i])
				StartCrawl(domains[i], chanBuff)
			}(i)

			// Increment outside go routine, or else we can't continue until that Domain finishes
			i++
		}
		// Rest the system
		time.Sleep(500 * time.Millisecond)
	}
}

// Gets all url links returned from html response
func StartCrawl(domain string, c chan bool) {
	// urlList := map[string]int{}  // Declare and Initialize Map for found links
	urlList := make(UrlList)

	defer func() {
		// Reset channel as true once Domain is finished, notify that we're done after this function
		c <- true
	}()

	// Retrieve all links in homepage body
	links := AllLinks(domain)

	// Iterate through grabbed links and insert into url list if not already in there
	for _, linkString := range links {
		if len(linkString) < 2 {
			continue
		}

		// if relative url append to Domain before insertion
		if string(linkString[0]) == "/" {
			linkString = "http://" + domain + linkString
		}

		// Finally check if the formatted linkString
		if _, ok := urlList[linkString]; !ok {
			urlList[linkString] = &UrlInfo{1, domain, domain}
		}
	}

	// Save length for urlList
	listLen := -1

	// Crawl Loop
	// Keep (re)crawling while there are new links being added
	for len(urlList) != len(subCrawl(urlList, domain)) {

		// If we go too high, or if list length doesn't change from last iteration break out
		if len(urlList) > 499 || listLen == len(urlList) {
			break
		}
		// Save length for this iteration
		listLen = len(urlList)
	}
}

// Crawl through all sub urls for Domain adding to original url list
// This is where the "real" crawl starts
func subCrawl(urlList UrlList, domain string) UrlList {
	// sub map
	var subMap = map[string][]string{}

	// iterate through unique map of urls
	for url, urlProps := range urlList {
		// check if url is in Domain and if status code has not been set
		if urlProps.status == 1 && strings.Contains(url, domain) {

			// Skip sub-links check if link is pdf, jpg, mp3, or png
			if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".pdf") ||
				strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".mp3") {
				// Then skip to next iteration without looking for links inside the picture, or song, or pdf, etc.
				continue
			}

			// Check for phone numbers on page
			phoneNumbers := findPhoneNumbers(url)

			// If we find any write them to the csv
			if phoneNumbers != nil {
				// clean duplicates from slice
				phoneNumbers = removeDuplicates(phoneNumbers)
				// add to the phoneMap
				subMap[url] = phoneNumbers
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
							link = "http://" + domain + link
							_, found := urlList[link]
							if !found {
								urlList[link] = &UrlInfo{1, domain, url}
							} else {
								continue
							}
						} else if strings.Contains(link, domain) {
							urlList[link] = &UrlInfo{1, domain, url}
						}
					}
				}
			}
			// END getting new links

			// Rest system between single page requests
			time.Sleep(50 * time.Millisecond)

		} else if !strings.Contains(url, domain) {
			// If the url string doesn't have the Domain remove it from url list
			delete(urlList, url)
		}
	}

	writeCSV(subMap)

	// return updated list with new urls
	return urlList
}


// helper func for writing phoneNumbers/urls to csv
func writeCSV(subMap map[string][]string) {
	// Creating csv writer
	wFile, err := os.OpenFile("./PhoneCrawl.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("ERR Opening CSV:", err)
	}
	defer wFile.Close()
	writer := csv.NewWriter(wFile)

	for url, phoneNumbers := range subMap {
		// join slice into one string
		numberString := strings.Join(phoneNumbers, ", ")
		// append the URL to the list of phone numbers and write to csv
		wErr := writer.Write([]string{url, numberString})
		if wErr != nil {
			log.Println(wErr)
			log.Fatal("ERR writing:", url, numberString)
		}
		writer.Flush()
	}
}

// help remove duplicates from phone number slice
func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// Helper func
func findPhoneNumbers(url string) []string {
	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Println("ERR Phone Scrape GET: ", err)
		return nil
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ERR Phone Scrape reading HTTP body: ", err)
		return nil
	}

	// Regex to match Phone Numbers
	re := regexp.MustCompile(`(\d{3}[-.\s]\d{3}[-.\s]\d{4}|\(\d{3}\)\s*\d{3}[-.\s]\d{4}|\d{3}[-.\s]\d{4})`)
	numbers := re.FindAllString(string(body), -1)
	if numbers == nil {
		return nil
	}
	// return the list of found numbers
	return numbers
}

// Takes url string, returns a slice of strings equal to the "href" attributes from anchor links found in the html.
func AllLinks(url string) []string {
	var links []string
	var col []string

	// If no protocol prepend it to url
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: All Links Failed to crawl \"" + url + "\"")
		fmt.Println(err)
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
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					tl := removeHash(attr.Val)
					col = append(col, tl)
					resolve(&links, col)
				}
			}
		}
	}
}

// Removes # from the link
func removeHash(link string) string {
	if strings.Contains(link, "#") {
		var index int
		for n, str := range link {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		return link[:index]
	}
	return link
}

// Checks if a url exits in the slice.
func check(sl []string, s string) bool {
	var check bool
	for _, str := range sl {
		if str == s {
			check = true
			break
		}
	}
	return check
}

// Adds links to the link slice and checks for no repetition in collection
func resolve(sl *[]string, ml []string) {
	for _, str := range ml {
		if check(*sl, str) == false {
			*sl = append(*sl, str)
		}
	}
}
