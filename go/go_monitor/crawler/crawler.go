package crawler

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
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
var domainStatus = make(map[string]int) //1st index is status code, 2nd is offline count
// struct for properties we wanna keep on each url
type UrlInfo struct {
	status  int
	domain  string
	referer string
}

// Map for found links
type UrlList map[string]*UrlInfo

// Populate Domain list on init
func init() {
	getDomainList()
	// Sort list alphabetically
	sort.Strings(domainList)
}

// GetStatus return status codes of Domain list
func GetStatus() map[string]int {
	return domainStatus
}

// Run db crawl and return true to channel
func RunHomeCrawl(c chan bool) {
	homeCrawl(domainList)
	c <- true
}

// RunAllCrawl starts the all crawl
func RunAllCrawl() {
	lockFileName := "crawlLock.txt"

	// Get file stats
	info, err := os.Stat(lockFileName)

	// Check for crawl lock file so we don't duplicate cron job. (If it doesn't exist or file is 48hours old run crawl)
	if err != nil || time.Now().Sub(info.ModTime()) > (48*time.Hour) {
		// Have to double check if statement for mod time on file so we can validate err also
		if os.IsNotExist(err) || time.Now().Sub(info.ModTime()) > (48*time.Hour) {
			// Lock file doesn't exist; create the lock file
			f, fErr := os.Create(lockFileName)
			defer f.Close()
			if fErr != nil {
				os.Exit(1)
			}
			// Make sure we delete the file on crawl stop
			defer deleteFile(lockFileName)

			// Safe to start crawl
			log.Println("All Crawl Started.")
			// Check what domains to skip
			UpdateSkipMap()
			// Reset Domain list
			getDomainList()
			// Start all crawl
			allCrawl(domainList)
			log.Println("All Crawl End.")
		}
	}
}

// Crawl only homepages for status updates
func homeCrawl(domains []string) {
	for _, domain := range domains {
		resp, err := http.Get("https://" + domain)
		if err != nil {
			log.Println("GET HTTPS Request ERR: ", err)
			// retry without https
			resp, err = http.Get("http://" + domain)
			if err != nil {
				log.Println("GET HTTP Request ERR: ", err)
				continue
			}
		}

		statusCode := resp.StatusCode
		// If we get 5xx send to error handler daemon
		if statusCode >= 500 {
			DaemonAddError(domain)
		}

		// Update Domain status in the map structure
		domainStatus[domain] = resp.StatusCode

		// Explicit close connection for current Domain; (deferring will cause mem leak)
		resp.Body.Close()
	}
}

// Gets all url links in html returned from response
func SoloCrawl(domain string) (map[string]*DomainStats, int) {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}

	crawlID := 0

	// Insert new crawl into table to create crawlID
	res, err := db.Exec("INSERT INTO crawl (start_time) VALUES (?)", time.Now())
	if err != nil {
		// Close DB
		db.Close()
		fmt.Println("Exec err:", err.Error())
	} else {
		fmt.Println("Crawling")
		// If no error inserting new crawl, we grab the last insert as crawlID
		tempCrawlID, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Crawl ID Error:", err.Error())
		}
		crawlID = int(tempCrawlID)

		urlList := make(UrlList)

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

		for len(urlList) != len(subCrawl(urlList, domain, crawlID, db)) {

			// If we go too high, or if list length doesn't change from last iteration break out
			if len(urlList) > 499 || listLen == len(urlList) {
				break
			}
			// Save length for this iteration
			listLen = len(urlList)
		}

		// Send email report after manual crawl finishes
		return SoloDomainReport(domain, crawlID)
	}

	// Return empty values if error
	return map[string]*DomainStats{}, 0
}

// allCrawl Everything
func allCrawl(domains []string) {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	startTime := time.Now()

	// Insert new crawl into table to create crawlID
	res, err := db.Exec("INSERT INTO crawl (start_time) VALUES (?)", startTime)
	if err != nil {
		fmt.Println("Exec err:", err.Error())
	} else {
		// If no error inserting new crawl, we grab the last insert as crawlID
		crawlID, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Crawl ID Error:", err.Error())
		}

		// Channel buffer for concurrency of 12
		chanBuff := make(chan bool, 12)
		// Prime queue
		for len(chanBuff) < 12 {
			chanBuff <- true
		}

		// START Crawling!
		for i := 0; i < len(domains); {
			// Runs all domains at once. Fast, but DON'T use this (causes DOS)
			//go func(i int) {
			//	defer wg.Done()
			//	StartCrawl(domains[i], int(crawlID), chan1, db)
			//}(i)

			// Check if the Domain is in the skip list
			if skipMap[domains[i]] {
				// If it is: then increment i and continue, so we can skip it
				i++
				continue
			}

			// Check if we have a val in buffer
			if <-chanBuff {
				go func(i int) {
					// log.Println("START", i, domains[i])
					// Run
					StartCrawl(domains[i], int(crawlID), chanBuff)
				}(i)

				// Increment outside go routine, or else we can't continue until that Domain finishes
				i++
			}
			// Rest the system
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Gets all url links returned from html response
func StartCrawl(domain string, crawlID int, c chan bool) {
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

	// DB connect; separate for each domain
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Save length for urlList
	listLen := -1

	// Crawl Loop
	// Keep (re)crawling while there are new links being added
	for len(urlList) != len(subCrawl(urlList, domain, crawlID, db)) {

		// If we go too high, or if list length doesn't change from last iteration break out
		if len(urlList) > 499 || listLen == len(urlList) {
			break
		}
		// Save length for this iteration
		listLen = len(urlList)
	}

	// Start updating status table after crawl iteration finished.
	go UpdateStatusTable(domain, crawlID)
}

// Crawl through all sub urls for Domain adding to original url list
// This is where the "real" crawl starts
func subCrawl(urlList UrlList, domain string, crawlID int, dbIn *sql.DB) UrlList {
	// iterate through unique map of urls
	for url, urlProps := range urlList {
		// check if url is in Domain and if status code has not been set
		if urlProps.status == 1 && strings.Contains(url, domain) {
			urlErr := ""
			redirects := 0

			// Start request
			preResp, respErr := http.NewRequest("GET", url, nil)
			if respErr != nil {
				fmt.Println("PreResp Err, Continuing: ", respErr)
				urlErr = "No Response"
				handleOutage(dbIn, domain, url, 0, crawlID, urlProps.referer)
				// Skip to next url if error
				continue
			}

			// Need to separate out into transport request since we want redirects
			linkResp, linkErr := http.DefaultTransport.RoundTrip(preResp)
			if linkErr != nil {
				fmt.Println("LinkResp Err, Continuing: ", linkErr)
				urlErr = "No Response"
				handleOutage(dbIn, domain, url, 0, crawlID, urlProps.referer)
				continue
			} else {
				// Get requested url string
				reqURL := linkResp.Request.URL.String()
				// Get landing url string
				landingURL := linkResp.Header.Get("Location")

				// If the landing url doesn't contain the domain then we've gone off track; skip iteration
				if !strings.Contains(landingURL, domain) {
					continue
				}

				// If Connection Successful, grab status code
				code := linkResp.StatusCode

				// HTTP to HTTPS handling, status code would be 301
				if code == 301 {
					// Check for http to https redirects; by comparing request url with landing url sans protocols
					if strings.Trim(reqURL, "http") == strings.Trim(landingURL, "https") {
						http.Get(landingURL)
						resp, err := http.Get(landingURL)
						if err != nil {
							log.Println("HTTP to HTTPS Error: ", err)
						}

						// Update the url to have HTTPS
						url = landingURL

						// Set the code as the status code from the HTTPS GET request so we can handle properly
						code = resp.StatusCode
						resp.Body.Close()
					}
				}

				// Set the status so it doesn't stay one and keep getting crawled
				urlProps.status = code

				// *!* After we are done with current loop's url: we must manually close connection to avoid mem leak *!*
				linkResp.Body.Close()

				if code >= 400 {
					linkResp.Body.Close()
					// Log in DB outages, then skip to next iteration
					handleOutage(dbIn, domain, url, code, crawlID, urlProps.referer)
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
							handleOutage(dbIn, domain, url, code, crawlID, urlProps.referer)
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
							// Ignore http to https redirects, by retracting the increment
							if strings.Trim(resp.Request.URL.String(), "http") == strings.Trim(nextURL, "https") {
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
							handleOutage(dbIn, domain, url, code, crawlID, urlProps.referer)
							break
						}
					}
				}
				// END Check redirects

				// Make sure we are only dealing with good requests from here down
				if code > 200 {
					continue
				}

				// Response Time and TTFB, in milliseconds, *!* MAKE SURE we pass in final url from redirects
				ttfb, respTime := timeGet(nextURL)

				// Check if link is pdf, jpg, mp3, or png
				if strings.HasSuffix(nextURL, ".jpg") || strings.HasSuffix(nextURL, ".pdf") ||
					strings.HasSuffix(nextURL, ".png") || strings.HasSuffix(nextURL, ".mp3") {

					// Insert into DB early if not a normal web page
					insertPage(dbIn, domain, nextURL, code, respTime, ttfb, urlErr, redirects, crawlID)
					// Then skip to next iteration without looking for links inside the picture, or song, or pdf, etc.
					continue
				}

				// Check for closing /html tag, return true if /html is present
				closeTag := checkCloseTag(nextURL)
				if !closeTag {
					urlErr = "/html"
				}

				// Grab all links on current url page
				links := AllLinks(nextURL)

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

				// Insert crawl info into DB
				insertPage(dbIn, domain, nextURL, code, respTime, ttfb, urlErr, redirects, crawlID)

				// Rest system between single page requests
				time.Sleep(50 * time.Millisecond)
			}
			// END if linkResp Successful
		} else if !strings.Contains(url, domain) {
			// If the url string doesn't have the Domain remove it from url list
			//fmt.Println("deleting: ", url)
			delete(urlList, url)
		}
		// END check if url is in Domain and status code
	}
	// return updated list with new urls
	return urlList
}

// Get time to first byte and response time, returns two floats
func timeGet(url string) (float64, float64) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Time Get Err: ", err)
		return 0, 0
	}

	var start, dns time.Time
	ttfb := time.Since(start)

	// Trace setup
	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
			//fmt.Printf("Time to first byte: %v\n", ttfb)
		},
	}

	start = time.Now()
	// Make the request
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	timeReq, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println("TimeReq Err: "+url, err)
		return 0, 0
	}
	defer timeReq.Body.Close() // *!* CLOSE THE FLIPPIN REQUEST!!!

	respTime := time.Since(start)

	// Convert time values(nanoseconds) to milliseconds
	msTTFB := float64(ttfb) / float64(time.Millisecond)
	msRespTime := float64(respTime) / float64(time.Millisecond)

	return msTTFB, msRespTime
}

// Check response body for closing /html tags
func checkCloseTag(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Check Closing /HTML Tag ERR: ", url)
		log.Println(err)
		return true
	}
	defer resp.Body.Close() // close Body when the function returns

	body, err := ioutil.ReadAll(resp.Body)
	lowerHTML := strings.Contains(string(body), "</html>")
	upperHTML := strings.Contains(string(body), "</HTML>")
	return lowerHTML || upperHTML
}

// Takes url string, returns a slice of strings equal to the "href" attributes from anchor links found in the html.
func AllLinks(url string) []string {
	links := []string{}
	col := []string{}

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

// Finds optimizely tags on domains
func OptlyLinks(url string) []string {
	links := []string{}
	col := []string{}

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
		if tokenType == html.StartTagToken && token.DataAtom.String() == "script" {
			for _, attr := range token.Attr {

				if attr.Key == "src" {
					if strings.Contains(attr.Val, "optimizely") {
						tl := removeHash(attr.Val)
						col = append(col, tl)
						resolve(&links, col)
					}
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
