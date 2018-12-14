package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"../crawler"
	"../creds"
	"crypto/subtle"
	"log"
)

// TemplateData struct to hold on the values passed to template
type TemplateData struct {
	Domains map[string]int
	Time    string
}

// DomainInfo struct for Graph View info we want for each individual domain
type DomainInfo struct {
	FacilityName    string
	Status          int
	AvgResponse     float64
	Outages         int
	Errors          int
	LastUpdate      string
	GraphDataOutage []float64
	GraphData404    []float64
}

var ManualCrawlCache = map[string]*crawler.DomainStats{}

// Run start up the router
func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}
	mux := http.NewServeMux()

	// Display React View
	mux.Handle("/", HandlerBasicAuth(http.FileServer(http.Dir("templates/build"))))
	// mux.Handle("/404", HandlerBasicAuth(http.FileServer(http.Dir("templates/build"))))
	// mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "templates/build/index.html")
	//})

	// Starts manual crawl on single domain, Returns stats to page
	mux.HandleFunc("/api/runCrawl", BasicAuth(manualCrawl))

	// Returns weekly data
	mux.HandleFunc("/api/monitorstatus", BasicAuth(monitorStatusPage))

	// Returns monthly data
	mux.HandleFunc("/api/monthlymonitorstatus", BasicAuth(monitorMonthlyStatusPage))

	// Returns list of 404
	mux.HandleFunc("/api/404list", get404List)
	// Returns list of 404
	mux.HandleFunc("/api/getFacilities", getFacilityList)

	// Heartbeat
	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Ok")
		})

	// Start go routine for server to listen on port
	http.ListenAndServe(":"+port, mux)
}

// Middleware that requires builder7 user/pass to access route; For Handler
func HandlerBasicAuth(handler http.Handler) http.HandlerFunc {
	// Creds needed to access page
	username := creds.AuthUser
	password := creds.AuthPass
	realm := "Enter username and password"

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler.ServeHTTP(w, r)
	}
}

// For HandlerFunc
func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	// Creds needed to access page
	username := creds.AuthUser
	password := creds.AuthPass
	realm := "Enter username and password"

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

// Manually start a full crawl for a domain
func manualCrawl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		//accessToken := "9=D2QWf%oP#o:4^V<D+$rZx:n~F*@?"
		//
		//// Check for token header
		//token := r.Header.Get("x-acadia-token")
		//// Reject if not valid
		//if token != accessToken {
		//	log.Println("ERR: Invalid Token")
		//	fmt.Fprint(w, "ERR: Invalid Token")
		//	return
		//}

		// Get form values
		domain := r.FormValue("domain")
		userEmail := r.FormValue("email")
		crawlType := r.FormValue("crawlType")

		// Make sure values for domain and email
		if domain == "" || userEmail == "" {
			log.Println("Form Data 'domain' or 'email' is missing")
			fmt.Fprint(w, "ERR: 'domain' or 'email' is missing")
			return
		}

		// Check if cached in map; return if it is
		if _, ok := ManualCrawlCache[domain]; ok {
			// right out response
			fmt.Fprint(w, "Crawl Started: "+domain)

			// If difference between time is less than an hour return here without updating cache
			if ManualCrawlCache[domain].TimeStamp.Sub(time.Now()) > (-1 * time.Hour) {

				// Split input email into array of email strings
				sepEmail := strings.Split(userEmail, ",")

				// Grab cached response
				list404 := ManualCrawlCache[domain].List404
				count404 := len(ManualCrawlCache[domain].List404)

				// Generate solo csv
				crawler.Solo404CSV(list404, domain)
				// Email solo csv
				crawler.SoloEmail404(sepEmail, domain, count404)
				return
			}
		}

		// Crawl/Email routine, so we don't hold up response
		go func() {
			// Split input email into array of email strings
			sepEmail := strings.Split(userEmail, ",")

			// Start new crawl request; receives domain stats object
			crawlResp, crawlID := crawler.SoloCrawl(domain)

			// Handle whether user ordered sitemap or just 404's
			if crawlType == "sitemap" {
				sitemapList := crawlResp[domain].Sitemap
				// CSV for sitemap
				crawler.SitemapCSV(sitemapList, domain)
				// Email sitemap csv
				crawler.SitemapEmail(sepEmail, domain)
			} else {
				list404 := crawlResp[domain].List404
				count404 := len(crawlResp[domain].List404)
				// Generate solo 404 csv
				crawler.Solo404CSV(list404, domain)
				// Email solo 404 csv
				crawler.SoloEmail404(sepEmail, domain, count404)
			}

			// Store it in the cache
			ManualCrawlCache[domain] = crawlResp[domain]

			statusJSON, _ := json.Marshal(crawlResp)

			log.Println("Cleaning DB", crawlID, string(statusJSON))
			// Clean up db after solo crawl finishes; delete the rows after we handle data temporarily stored in DB by crawlID
			crawler.DeleteSoloData(crawlID)
		}()

		// right out response
		fmt.Fprint(w, "Crawl Started: "+domain)
	}
}

// Return domain status card info for front-end
func monitorStatusPage(w http.ResponseWriter, r *http.Request) {
	// set CORS header so we can test
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get status codes map(dict) from crawler
	statusCodes := crawler.GetStatus()

	// access cached status object; (or new one if expired)
	cacheStatus, lastUpdate := crawler.GetWeeklyMonitorReport()

	// create new response map for combine
	resp := map[string]DomainInfo{}

	// combine status codes with cached info into
	for key := range cacheStatus {
		resp[key] = DomainInfo{
			cacheStatus[key].FacilityName,
			statusCodes[key],
			cacheStatus[key].AvgResponse,
			cacheStatus[key].Outages,
			cacheStatus[key].Errors,
			lastUpdate.Format(time.Stamp),
			cacheStatus[key].GraphDataOutage,
			cacheStatus[key].GraphData404,
		}
	}

	statusJSON, _ := json.Marshal(resp)
	// right out status response
	fmt.Fprint(w, string(statusJSON))
}

// Return domain status card info for front-end
func monitorMonthlyStatusPage(w http.ResponseWriter, r *http.Request) {
	// set CORS header so we can test
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get status codes map(dict) from crawler
	statusCodes := crawler.GetStatus()

	// access cached status object; (or new one if expired)
	cacheStatus, lastUpdate := crawler.GetMonthlyMonitorReport()

	// create new response map for combine
	resp := map[string]DomainInfo{}

	// combine status codes with cached info into
	for key := range cacheStatus {
		resp[key] = DomainInfo{
			cacheStatus[key].FacilityName,
			statusCodes[key],
			cacheStatus[key].AvgResponse,
			cacheStatus[key].Outages,
			cacheStatus[key].Errors,
			lastUpdate.Format(time.Stamp),
			cacheStatus[key].GraphDataOutage,
			cacheStatus[key].GraphData404,
		}
	}

	statusJSON, _ := json.Marshal(resp)
	// right out status response
	fmt.Fprint(w, string(statusJSON))
}

// Return domain status card info for front-end
func get404List(w http.ResponseWriter, r *http.Request) {
	// set CORS header so we can test
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fofSet := crawler.Get404DomainList()

	statusJSON, _ := json.Marshal(fofSet)
	// right out status response
	fmt.Fprint(w, string(statusJSON))
}

// Return domain/facility name pairs
func getFacilityList(w http.ResponseWriter, r *http.Request) {
	// set CORS header so we can test
	w.Header().Set("Access-Control-Allow-Origin", "*")

	facilitySet := crawler.GenerateFacilityList()

	statusJSON, _ := json.Marshal(facilitySet)
	// right out status response
	fmt.Fprint(w, string(statusJSON))
}
