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
	"golang.org/x/net/context"
	// Go get...
	"golang.org/x/net/html"
)

var cancelCrawl = false
var skipMap = map[string]bool{
	"www.test.com": true,
}
var domainList []string
var domainStatus = make(map[string]int)
// Page Stats
type PageStats struct {
	domain     string
	url        string
	statusCode int
	respTime   float64
	ttfb       float64
	urlErr     string
	redirects  int
	crawlID    int
}

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

// Crawl only homepages for status updates
func homeCrawl(domains []string) {
	for _, domain := range domains {
		resp, err := http.Get("https://" + domain)
		if err != nil {
			log.Println("Home Crawl HTTPS Request ERR: ", domain, err)
			// retry without https
			resp, err = http.Get("http://" + domain)
			if err != nil {
				log.Println("Home Crawl HTTP Request ERR: ", domain, err)
				continue
			}
		}

		statusCode := resp.StatusCode
		// If we get 5xx AND it has NOT already been logged in domainStatus; send to error handler daemon
		if statusCode >= 500 && domainStatus[domain] != statusCode {
			DaemonAddError(domain)
		}

		// Update Domain status in the map structure
		domainStatus[domain] = resp.StatusCode

		// Explicit close connection for current Domain; (deferring will cause mem leak)
		resp.Body.Close()
	}
}

// RunAllCrawl starts the all crawl
func RunAllCrawl() {
	lockFileName := "crawlLock.txt"

	// Get file stats
	info, err := os.Stat(lockFileName)

	// Check for crawl lock file so we don't duplicate cron job. (If it doesn't exist or file is 4hours old run crawl)
	if err != nil || time.Now().Sub(info.ModTime()) > (4*time.Hour) {
		// Have to double check if statement for mod time on file so we can validate err also
		if os.IsNotExist(err) || time.Now().Sub(info.ModTime()) > (4*time.Hour) {
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

		// Channel buffer for concurrency of 8
		chanBuff := make(chan bool, 8)
		// Prime queue
		for len(chanBuff) < 8 {
			chanBuff <- true
		}

		cancelCrawl = false
		// Create a new context, with cancellation function
		ctx, cancel := context.WithCancel(context.Background())

		// START Crawling!
		for i := 0; i < len(domains); {
			// Runs all domains at once. Fast, but DON'T use this (causes DOS)
			//go func(i int) {
			//	defer wg.Done()
			//	StartCrawl(domains[i], int(crawlID), chan1, db)
			//}(i)

			// If we get cancellation signal
			if cancelCrawl {
				cancel()
				return
			}

			// Check if the Domain is in the skip list
			if skipMap[domains[i]] {
				// If it is: then increment i and continue, so we can skip it
				i++
				continue
			}

			// Check if we have a val in buffer
			if <-chanBuff {
				log.Println(i, domains[i])
				go func(i int) {
					// log.Println("START", i, domains[i])
					// Run
					StartCrawl(ctx, domains[i], int(crawlID), chanBuff)
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
func StartCrawl(ctx context.Context, domain string, crawlID int, c chan bool) {
	// urlList := map[string]int{}  // Declare and Initialize Map for found links
	urlList := make(UrlList)

	defer func() {
		// Reset channel as true once Domain is finished, notify that we're done after this function
		c <- true
	}()

	// check for cancellation
	select {
	case <-ctx.Done():
		return

	// If no cancellation business as usual
	default:
		// Retrieve all links in homepage body
		links := AllLinks(domain)

		// Iterate through grabbed links and insert into url list if not already in there
		for _, linkString := range links {
			if len(linkString) < 2 {
				continue
			}
			// if relative url append to Domain before insertion
			if strings.HasPrefix(linkString, "/") {
				linkString = "https://" + domain + linkString
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

		// Crawl Loop; Keep (re)crawling while there are new links being added
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
}

// Crawl through all sub urls for Domain adding to original url list; where the "real" crawl starts
func subCrawl(urlList UrlList, domain string, crawlID int, dbIn *sql.DB) UrlList {
	// iterate through unique map of urls
	for url, urlProps := range urlList {
		// check if url is in Domain and if status code has not been set
		if urlProps.status == 1 && strings.Contains(url, domain) {
			// struct to hold new page info
			page := PageStats{
				domain:     domain,
				url:        url,
				statusCode: 1,
				respTime:   0,
				ttfb:       0,
				urlErr:     "",
				redirects:  0,
				crawlID:    crawlID,
			}
			ctype := ""

			nextURL := url
			// Check redirects up to 10
			for page.redirects < 10 {
				client := &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					}}

				// GET Request
				resp, err := client.Get(nextURL)
				if err != nil {
					log.Println("Resp Err, Continuing: ", err)
					page.urlErr = "No Response"
					handleOutage(dbIn, domain, url, 0, crawlID, urlProps.referer)
					// Skip to next url if error
					break
				}

				// Handle redirects
				if 299 < resp.StatusCode && resp.StatusCode < 400 {
					reqURL := resp.Request.URL.String()
					nextURL = resp.Header.Get("Location")
					// Check for http to https; by comparing request url with next url sans protocols
					if !(strings.Trim(reqURL, "http") == strings.Trim(nextURL, "https")) {
						// Increment redirect count if NOT a HTTPS promotion
						page.redirects += 1
					}
				} else {
					// Set the status so it doesn't stay one and keep getting crawled
					urlProps.status = resp.StatusCode
					ctype = resp.Header.Get("Content-Type")
					page.url = nextURL

					// Done with redirects; process page
					pageScan(dbIn, resp, &page, urlProps.referer)
					break
				}
			}

			// Only get page links if it is a normal html page
			if strings.HasPrefix(ctype, "text/html") {
				urlList = retrieveLinks(urlList, &page)
			}

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

// Helper func for page stats
func pageScan(dbIn *sql.DB, resp *http.Response, page *PageStats, referer string) {
	// Update Page struct
	page.statusCode = resp.StatusCode

	// If the landing url doesn't contain the domain then we've gone off track; skip
	if !strings.Contains(page.url, page.domain) {
		return
	}

	if page.statusCode >= 400 {
		resp.Body.Close()
		// Log in DB outages, then skip to next iteration
		handleOutage(dbIn, page.domain, page.url, page.statusCode, page.crawlID, referer)
		// Break crawl if we hit 500 error
		if page.statusCode >= 500 {
			cancelCrawl = true
		}
		return
	}

	// read the response body to a variable
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// *!* After we are done with current loop's url: we must manually close connection to avoid mem leak *!*
	resp.Body.Close()

	// Response Time and TTFB, in milliseconds, *!* MAKE SURE we pass in final url from redirects
	ttfb, respTime := timeGet(page.url)

	// Check for closing /html tag, return true if /html is present
	closeTag := checkCloseTag(bodyBytes)
	if !closeTag {
		page.urlErr = "/html"
	}

	// Insert crawl info into DB
	insertPage(dbIn, page.domain, page.url, page.statusCode, respTime, ttfb, page.urlErr, page.redirects, page.crawlID)
}

// Helper for sub page link retrieval
func retrieveLinks(urlList UrlList, page *PageStats) UrlList {
	// Grab all links on current url page
	links := AllLinks(page.url)

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
				if strings.HasPrefix(link, "/") {
					link = "https://" + page.domain + link
					_, found := urlList[link]
					if !found {
						urlList[link] = &UrlInfo{1, page.domain, page.url}
					} else {
						continue
					}
				} else if strings.Contains(link, page.domain) {
					urlList[link] = &UrlInfo{1, page.domain, page.url}
				}
			}
		}
	}
	// END getting new links
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
func checkCloseTag(body []byte) bool {
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
