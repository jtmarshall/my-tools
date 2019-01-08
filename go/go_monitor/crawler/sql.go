package crawler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"../creds"
)

// import credentials for DB connection
const connectString = creds.DbConnectString

// Struct to create format for Domain metrics
type DomainStats struct {
	AvgRespTime float64
	MaxRespTime float64
	AvgTTFB     float64
	MaxTTFB     float64
	CountURL    int
	TotalErr    int
	List404     map[string][]string
	Sitemap     map[string][]string
	TimeStamp   time.Time
}

// Properties for 404 email
type Four04Props struct {
	Page         sql.NullString
	Domain       sql.NullString
	Referer      sql.NullString
	TimeStamp    time.Time
	FacilityName sql.NullString
}

// Properties for sitemap email
type SitemapProps struct {
	Page      sql.NullString
	Domain    sql.NullString
	TimeStamp time.Time
	Status    sql.NullString
	Referer   sql.NullString
}

// DomainInfo struct for all info we want for each individual domain
type DomainInfo struct {
	FacilityName    string
	Status          int
	AvgResponse     float64
	Outages         int
	Errors          int
	GraphDataOutage []float64
	GraphData404    []float64
}

// Global cache map to return to frontend
var weeklyCacheStatus = make(map[string]DomainInfo)
var monthlyCacheStatus = make(map[string]DomainInfo)
var cacheFacilities = make(map[string]string)
var cached404Domains = map[string][]Four04Props{}
var cache404List = map[string][]string{}
var lastStatusUpdate = time.Now()

// GetDomainList Query DB for Domain list
func getDomainList() {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var domain string

	rows, err := db.Query("SELECT domain FROM status")
	if err != nil {
		log.Fatal(err)
	}

	// clear list before we start adding again
	domainList = domainList[:0]

	for rows.Next() {
		err := rows.Scan(&domain)
		if err != nil {
			log.Fatal(err)
		}
		domainList = append(domainList, domain)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	// Manually Close
	rows.Close()
}

// updates map of domains to skip in all crawl, pulled from DB
func UpdateSkipMap() {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var domain string

	rows, err := db.Query("SELECT domain FROM status WHERE skip_crawl = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	// use temporary skip map to collect values
	tempMap := map[string]bool{}

	for rows.Next() {
		err := rows.Scan(&domain)
		if err != nil {
			log.Fatal(err)
		}
		tempMap[domain] = true
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	// Manually Close
	rows.Close()

	// update saved skip map
	skipMap = tempMap
}

// Insert standard info into page table
// *!* Pass in DB pointer to avoid connection overflow!!! *!*
func insertPage(dbIn *sql.DB, domain string, url string, code int, respTime float64,
	ttfb float64, urlErr string, redirects int, crawlID int) {

	// Insert crawl info into DB
	_, err2 := dbIn.Exec(
		`INSERT INTO page (Domain, url, datetime, status_code, response_time, ttfb, error, redirects, crawl_id) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		domain, url, time.Now(), code, respTime, ttfb, urlErr, redirects, crawlID)
	if err2 != nil {
		log.Fatal("INSERT page ERR: ", err2)
	}
	// fmt.Println("Finish InsertDB Func")
}

// Called if status code is 5xx, insert info into outages table
func handleOutage(dbIn *sql.DB, domain string, url string, code int, crawlID int, refererURL string) {
	// Check if we need to send an email alert
	if code == 504 {
		// Send error to daemon to handle email
		DaemonAddError(url)
	}

	// Insert outage into DB outages table
	_, err2 := dbIn.Exec(
		"INSERT INTO outages (Domain, page, datetime, status_code, crawl_id, Referer) VALUES (?, ?, ?, ?, ?, ?)",
		domain, url, time.Now(), code, crawlID, refererURL)
	if err2 != nil {
		log.Println("INSERT outages ERR: ", err2)
	}
}

// After Domain crawl update status table
func UpdateStatusTable(domain string, crawlID int) {
	var (
		avgRespTime  float64 = 0
		maxRespTime  float64 = 0
		avgTTFB      float64 = 0
		maxTTFB      float64 = 0
		countURL             = 0
		totalErr             = 0
		totalCrawled         = 0
	)

	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Println("UPDATE STATUS DB CONNECT ERR: ", err)
	}
	defer db.Close()

	// Query for page STATS
	page, queryErr := db.Query(
		`SELECT AVG(response_time), MAX(response_time), AVG(ttfb), MAX(ttfb), COUNT(id)
				FROM page WHERE crawl_id = ? AND Domain = ?`,
		crawlID, domain)
	if queryErr != nil {
		log.Println("QUERY PAGE ERR: ", queryErr)
		return
	}
	defer page.Close()

	for page.Next() {
		scanErr := page.Scan(&avgRespTime, &maxRespTime, &avgTTFB, &maxTTFB, &countURL)
		if scanErr != nil {
			log.Println("PAGE.NEXT ERR: ", scanErr)
			return
		}
	}

	// Query for outages/errors
	outage, outageErr := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE crawl_id = ? AND Domain = ? AND status_code >= ?",
		crawlID, domain, 500)
	if outageErr != nil {
		log.Println("OUTAGE ERR: ", outageErr)
		return
	}
	defer outage.Close()

	for outage.Next() {
		outErr := outage.Scan(&totalErr)
		if outErr != nil {
			log.Println("OUTAGE.NEXT ERR: ", outErr)
			return
		}
	}

	// Get total_crawled urls for whole crawlID, so far
	crawlCount, crawlCountErr := db.Query("SELECT COUNT(id) FROM page WHERE crawl_id = ?", crawlID)
	if crawlCountErr != nil {
		log.Println("CRAWL COUNT ERR: ", crawlCountErr)
		return
	}
	defer crawlCount.Close()

	for crawlCount.Next() {
		countErr := crawlCount.Scan(&totalCrawled)
		if countErr != nil {
			log.Println("CRAWL COUNT.NEXT ERR: ", countErr)
			return
		}
	}

	// Update status table for Domain
	_, updateErr := db.Exec(
		`UPDATE status SET datetime = ?, avg_ttfb = ?, avg_response = ?, total_errors = ?, total_urls = ?,
					max_response = ?, max_ttfb = ? WHERE Domain = ?`,
		time.Now(), avgTTFB, avgRespTime, totalErr, countURL, maxRespTime, maxTTFB, domain)
	if updateErr != nil {
		log.Println("UPDATE ERR: ", updateErr)
		return
	}

	// Finally update crawl end time and "total_crawled" after all stuff for this Domain is done  for crawl table
	_, updateTotalCrawedlErr := db.Exec(
		"UPDATE crawl SET end_time = ?, total_crawled = ? WHERE id = ?", time.Now(), totalCrawled, crawlID)
	if updateTotalCrawedlErr != nil {
		log.Println("UPDATE TOTAL_CRAWLED ERR: ", updateTotalCrawedlErr)
		return
	}
}

// Creates and returns a domain/facility name map from DB
func GenerateFacilityList() map[string]string {

	if len(cached404Domains) < 1 || lastStatusUpdate.Sub(time.Now()) < (-1*time.Hour) {
		// Generate list if empty or cache is old
		facilityList := map[string]string{}

		// DB connect
		db, err := sql.Open("mysql", connectString)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		// Query domain/facility name pairs
		status, queryErr := db.Query(`SELECT domain, facility_name FROM status`)
		if queryErr != nil {
			log.Println("QUERY Status Table ERR: ", queryErr)
			return facilityList
		}
		defer status.Close()

		for status.Next() {
			var (
				domain       sql.NullString
				facilityName sql.NullString
			)
			scanErr := status.Scan(&domain, &facilityName)
			if scanErr != nil {
				log.Println("STATUS.NEXT ERR: ", scanErr)
				return facilityList
			}

			facilityList[domain.String] = facilityName.String
		}

		cacheFacilities = facilityList
	}

	return cacheFacilities
}

// Returns 404 list separated by domain
func Get404DomainList() map[string][]Four04Props {
	if len(cached404Domains) < 1 {
		// Generate list if empty
		Generate404List()
	}

	if lastStatusUpdate.Sub(time.Now()) < (-1 * time.Hour) {
		// If expired start update in separate routine so we can still return stale values
		go Generate404List()
		lastStatusUpdate = time.Now()
	}

	return cached404Domains
}

// Returns 404list with no separation, only alphabetical
func Get404List() map[string][]string {
	if len(cache404List) < 1 {
		// Generate list if empty
		Generate404List()
	}

	if lastStatusUpdate.Sub(time.Now()) < (-1 * time.Hour) {
		// If expired start update in separate routine so we can still return stale values
		go Generate404List()
		lastStatusUpdate = time.Now()
	}

	return cache404List
}

// Create list of found 404's
func Generate404List() {

	tempFacilities := GenerateFacilityList()

	// DB connect
	db, err := sql.Open("mysql", connectString+"?parseTime=true")
	if err != nil {
		log.Println("DB Connect ERR: ", err)
	}
	defer db.Close()

	toTime := time.Now()
	// get previous day
	fromTime := toTime.Add(-24 * time.Hour)

	// Query for 404
	query404, queryErr := db.Query(
		`SELECT DISTINCT domain, page, datetime, Referer FROM outages WHERE status_code = ? AND datetime BETWEEN ? AND ? GROUP BY page, Referer`,
		404, fromTime, toTime)
	if queryErr != nil {
		log.Println("QUERY404 PAGE ERR: ", queryErr)
	}
	defer query404.Close()

	// Clear the list
	cache404List = map[string][]string{}

	// Temporary map for list gathering
	temp404List := map[string][]string{}

	for query404.Next() {
		// var page string
		props := Four04Props{}

		scan404Err := query404.Scan(&props.Domain, &props.Page, &props.TimeStamp, &props.Referer)
		if scan404Err != nil {
			log.Println("query404.NEXT ERR: ", scan404Err)
			break
		}

		temp404List[props.Page.String] = []string{props.Page.String, props.Domain.String, props.Referer.String, props.TimeStamp.String()}

		// Set facility name value from temp cached facilities
		props.FacilityName.String = tempFacilities[props.Domain.String]

		// grouping cache by domain
		currentDomain := props.Domain.String
		cached404Domains[currentDomain] = append(cached404Domains[currentDomain], props)
	}

	// Update 404 cached list
	cache404List = temp404List
}

// StatusReport object for front-end
func StatusReportUpdateWeekly() {
	// update status cache object in go routine
	go func() {
		// get time
		toTime := time.Now()

		// DB connect
		db, err := sql.Open("mysql", connectString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		tempCacheStatus := make(map[string]DomainInfo)

		// update time
		lastStatusUpdate = time.Now()
		// previous week for query
		last7Time := toTime.AddDate(0, 0, -7)
		numDays := int(toTime.Sub(last7Time).Hours() / 24)

		// Loop through all domains getting metrics
		for _, domain := range domainList {
			// Reset vars for every Domain
			var (
				// Use nullable types for temp values so the scan doesn't throw a fit
				facilityName = ""
				avgResponse  sql.NullFloat64
				outages      = 0
				err404       = 0
				status       = 0
				weeklyOutage = make([]float64, numDays)
				weekly404    = make([]float64, numDays)
			)

			// Query for status table for STATS
			page, queryErr := db.Query(
				`SELECT facility_name, avg_response FROM status WHERE domain = ?`,
				domain)
			if queryErr != nil {
				fmt.Println("QUERY PAGE ERR: ", queryErr)
				break
			}
			defer page.Close()

			for page.Next() {
				scanErr := page.Scan(&facilityName, &avgResponse)
				if scanErr != nil {
					fmt.Println("PAGE.NEXT ERR: ", scanErr)
					break
				}
			}

			// Query for 404's
			find404, find404Err := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code = 404 AND datetime BETWEEN ? AND ?",
				domain, toTime.AddDate(0, 0, -1), toTime)
			if find404Err != nil {
				fmt.Println("Err404 ERR: ", find404Err)
				break
			}
			defer find404.Close()

			for find404.Next() {
				out404Err := find404.Scan(&err404)
				if out404Err != nil {
					fmt.Println("Err404.NEXT ERR: ", out404Err)
					break
				}
			}

			// Query for outages(500)
			outage, outageErr := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code >= 500 AND datetime BETWEEN ? AND ?",
				domain, last7Time, toTime)
			if outageErr != nil {
				fmt.Println("OUTAGE ERR: ", outageErr)
				break
			}
			defer outage.Close()

			for outage.Next() {
				outErr := outage.Scan(&outages)
				if outErr != nil {
					fmt.Println("OUTAGE.NEXT ERR: ", outErr)
					break
				}
			}

			// If outages/404 exist in the previous value, query outages by day
			if outages > 0 || err404 > 0 {
				tempGraph := make([]float64, numDays)
				tempWeekly404 := make([]float64, numDays)

				// iterate over starting with latest date
				for i := int(numDays) - 1; i >= 0; i-- {
					// subtract a day for each number of days in query
					tempTo := toTime.AddDate(0, 0, -i)
					tempFrom := toTime.AddDate(0, 0, -(i + 1))

					// Query for outages
					tempOut, tempOutErr := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code >= 500 AND datetime BETWEEN ? AND ?",
						domain, tempFrom, tempTo)
					if tempOutErr != nil {
						fmt.Println("SUB OUTAGE ERR: ", tempOutErr)
						break
					}

					for tempOut.Next() {
						// insert into tempGraph array index by increasing date
						outErr := tempOut.Scan(&tempGraph[len(tempGraph)-1-i])
						if outErr != nil {
							fmt.Println("SUB OUTAGE.NEXT ERR: ", outErr)
							break
						}
					}
					// manually close db connection because sub loop
					tempOut.Close()

					// Query for daily 404
					temp404, temp404Err := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code = 404 AND datetime BETWEEN ? AND ?",
						domain, tempFrom, tempTo)
					if temp404Err != nil {
						fmt.Println("SUB OUTAGE ERR: ", temp404Err)
						break
					}

					for temp404.Next() {
						// insert into tempGraph array index by increasing date
						tempErr := temp404.Scan(&tempWeekly404[len(tempWeekly404)-1-i])
						if tempErr != nil {
							fmt.Println("SUB TEMP404.NEXT ERR: ", tempErr)
							break
						}
					}
					// manually close db connection because sub loop
					temp404.Close()
				}

				weeklyOutage = tempGraph
				weekly404 = tempWeekly404
			}

			// update domain's info in local cache
			tempCacheStatus[domain] = DomainInfo{facilityName, status,
				avgResponse.Float64, outages, err404, weeklyOutage, weekly404}
		}
		// END Domain loop

		// update saved cache with tempCache obj
		weeklyCacheStatus = tempCacheStatus
	}()
}

// Return cached object with time
func GetWeeklyMonitorReport() (map[string]DomainInfo, time.Time) {

	return weeklyCacheStatus, lastStatusUpdate
}

// Updates monthly graph values in cache object for front-end
func StatusReportUpdateMonthly() {
	// update status cache object in go routine
	go func() {
		// get time
		toTime := time.Now()

		// DB connect
		db, err := sql.Open("mysql", connectString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		tempCacheStatus := make(map[string]DomainInfo)

		// update time
		lastStatusUpdate = time.Now()
		// previous week for query
		last30Time := toTime.AddDate(0, 0, -30)
		numDays := int(toTime.Sub(last30Time).Hours() / 24)

		// Loop through all domains getting metrics
		for _, domain := range domainList {
			// Reset vars for every Domain
			var (
				// Use nullable types for temp values so the scan doesn't throw a fit
				facilityName = ""
				avgResponse  sql.NullFloat64
				outages      = 0
				err404       = 0
				status       = 0
				monthlyGraph = make([]float64, numDays)
				monthly404   = make([]float64, numDays)
			)

			// Query for status table for STATS
			page, queryErr := db.Query(
				`SELECT facility_name, avg_response FROM status WHERE domain = ?`,
				domain)
			if queryErr != nil {
				fmt.Println("QUERY PAGE ERR: ", queryErr)
				break
			}
			defer page.Close()

			for page.Next() {
				scanErr := page.Scan(&facilityName, &avgResponse)
				if scanErr != nil {
					fmt.Println("PAGE.NEXT ERR: ", scanErr)
					break
				}
			}

			// Query for recent 404 (last day)
			find404, find404Err := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code = 404 AND datetime BETWEEN ? AND ?",
				domain, toTime.AddDate(0, 0, -1), toTime)
			if find404Err != nil {
				fmt.Println("Err404 ERR: ", find404Err)
				break
			}
			defer find404.Close()

			for find404.Next() {
				out404Err := find404.Scan(&err404)
				if out404Err != nil {
					fmt.Println("Err404.NEXT ERR: ", out404Err)
					break
				}
			}

			// Query for monthly outages
			outage, outageErr := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code >= 500 AND datetime BETWEEN ? AND ?",
				domain, last30Time, toTime)
			if outageErr != nil {
				fmt.Println("OUTAGE ERR: ", outageErr)
				break
			}
			defer outage.Close()

			for outage.Next() {
				outErr := outage.Scan(&outages)
				if outErr != nil {
					fmt.Println("OUTAGE.NEXT ERR: ", outErr)
					break
				}
			}

			// If outages exist in the previous value, query outages by day
			if outages > 0 || err404 > 0 {
				tempGraph := make([]float64, numDays)
				tempMonthly404 := make([]float64, numDays)

				// iterate over starting with latest date
				for i := int(numDays) - 1; i >= 0; i-- {
					// subtract a day for each number of days in query
					tempTo := toTime.AddDate(0, 0, -i)
					tempFrom := toTime.AddDate(0, 0, -(i + 1))

					// Query for outages/errors
					tempOut, tempOutErr := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND datetime BETWEEN ? AND ?",
						domain, tempFrom, tempTo)
					if outageErr != nil {
						fmt.Println("SUB OUTAGE ERR: ", tempOutErr)
						break
					}

					for tempOut.Next() {
						// insert into tempGraph array index by increasing date
						outErr := tempOut.Scan(&tempGraph[len(tempGraph)-1-i])
						if outErr != nil {
							fmt.Println("SUB OUTAGE.NEXT ERR: ", outErr)
							break
						}
					}
					// manually close db connection because sub loop
					tempOut.Close()

					// Query for daily 404
					temp404, temp404Err := db.Query("SELECT COUNT(DISTINCT page) FROM outages WHERE domain = ? AND status_code = 404 AND datetime BETWEEN ? AND ?",
						domain, tempFrom, tempTo)
					if temp404Err != nil {
						fmt.Println("SUB OUTAGE ERR: ", temp404Err)
						break
					}

					for temp404.Next() {
						// insert into tempGraph array index by increasing date
						tempErr := temp404.Scan(&tempMonthly404[len(tempMonthly404)-1-i])
						if tempErr != nil {
							fmt.Println("SUB TEMP404.NEXT ERR: ", tempErr)
							break
						}
					}
					// manually close db connection because sub loop
					temp404.Close()
				}

				monthlyGraph = tempGraph
				monthly404 = tempMonthly404
			}

			// update domain's info in local cache
			tempCacheStatus[domain] = DomainInfo{facilityName, status,
				avgResponse.Float64, outages, err404, monthlyGraph, monthly404}
		}
		// END Domain loop

		// update saved cache with tempCache obj
		monthlyCacheStatus = tempCacheStatus
	}()
}

// Return cached object with time
func GetMonthlyMonitorReport() (map[string]DomainInfo, time.Time) {

	return monthlyCacheStatus, lastStatusUpdate
}

// WeeklyEmailReport to query stats for all domains and send email once a week
func WeeklyEmailReport() {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create map for Domain keys and DomainsStats values
	statusMap := map[string]*DomainStats{}

	// 7 day time frame previous to now
	toTime := time.Now()
	fromTime := toTime.AddDate(0, 0, -7)

	// Loop through all domains getting metrics
	for _, domain := range domainList {
		// Reset vars for every Domain
		var (
			// Use nullable types for temp values so the scan doesn't throw a fit
			avgRespTime sql.NullFloat64
			maxRespTime sql.NullFloat64
			avgTTFB     sql.NullFloat64
			maxTTFB     sql.NullFloat64
			countURL    = 0
			totalErr    = 0
		)

		// Query for page STATS
		page, queryErr := db.Query(
			`SELECT AVG(response_time), MAX(response_time), AVG(ttfb), MAX(ttfb), COUNT(id)
					FROM page WHERE Domain = ? AND datetime BETWEEN ? AND ?`,
			domain, fromTime, toTime)
		if queryErr != nil {
			fmt.Println("QUERY PAGE ERR: ", queryErr)
			break
		}
		defer page.Close()

		for page.Next() {
			scanErr := page.Scan(&avgRespTime, &maxRespTime, &avgTTFB, &maxTTFB, &countURL)
			if scanErr != nil {
				fmt.Println("PAGE.NEXT ERR: ", scanErr)
				break
			}
		}

		// TODO: Query for distinct urls that are missing /html tag

		// Query for outages/errors
		outage, outageErr := db.Query("SELECT COUNT(DISTINCT id) FROM outages WHERE Domain = ? AND datetime BETWEEN ? AND ?",
			domain, fromTime, toTime)
		if outageErr != nil {
			fmt.Println("OUTAGE ERR: ", outageErr)
			break
		}
		defer outage.Close()

		for outage.Next() {
			outErr := outage.Scan(&totalErr)
			if outErr != nil {
				fmt.Println("OUTAGE.NEXT ERR: ", outErr)
				break
			}
		}

		// Create map for 404url keys and properties values
		four04Map := map[string][]string{}

		// Query for 404
		query404, queryErr := db.Query(
			`SELECT DISTINCT domain, page, datetime, Referer FROM outages WHERE status_code = ? AND datetime BETWEEN ? AND ?`,
			404, fromTime, toTime)
		if queryErr != nil {
			fmt.Println("QUERY404 PAGE ERR: ", queryErr)
		}
		defer query404.Close()

		for query404.Next() {
			props := Four04Props{}

			scan404Err := query404.Scan(&props.Domain, &props.Page, &props.TimeStamp, &props.Referer)
			if scan404Err != nil {
				fmt.Println("query404.NEXT ERR: ", scan404Err)
				break
			}
			// add 404 page to map with properties from query
			four04Map[props.Page.String] = []string{props.Page.String, props.Domain.String, props.Referer.String, props.TimeStamp.String()}
		}

		// Assign struct values to Domain
		statusMap[domain] = &DomainStats{
			// Get "real" float values from nullable types, and round them to 2 decimals before setting Domain values
			AvgRespTime: avgRespTime.Float64,
			MaxRespTime: maxRespTime.Float64,
			AvgTTFB:     avgTTFB.Float64,
			MaxTTFB:     maxTTFB.Float64,
			CountURL:    countURL,
			TotalErr:    totalErr,
			List404:     four04Map,
		}
	}
	// END Domain loop

	// Send email with passed in status map
	WeeklyEmail(statusMap, domainList)
}

// DeleteOldData removes rows from DB that are more than 30 days old
func DeleteOldData() {
	// 30 day time frame from now
	oldTime := time.Now().AddDate(0, 0, -30)

	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete rows from page older than 30 days
	_, deletePageErr := db.Exec(
		"DELETE FROM page WHERE datetime < ?", oldTime)
	if deletePageErr != nil {
		log.Println("DELETE PAGE ERR: ", deletePageErr)
	}

	// Delete rows from crawl older than 30 days
	_, deleteCrawlErr := db.Exec(
		"DELETE FROM crawl WHERE end_time < ?", oldTime)
	if deleteCrawlErr != nil {
		log.Println("DELETE CRAWL ERR: ", deleteCrawlErr)
	}

	// Delete rows from outage older than 30 days
	_, deleteOutagesErr := db.Exec(
		"DELETE FROM outages WHERE datetime < ?", oldTime)
	if deleteOutagesErr != nil {
		log.Println("DELETE OUTAGES ERR: ", deleteOutagesErr)
	}
}

// Similar report as Weekly Report, but runs only for given domain with given crawlID and doesn't send email
func SoloDomainReport(domain string, crawlID int) (map[string]*DomainStats, int) {
	// DB connect
	db, err := sql.Open("mysql", connectString+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create map for Domain keys and DomainsStats values
	statusMap := map[string]*DomainStats{}

	// Reset vars for every Domain
	var (
		// Use nullable types for temp values so the scan doesn't throw a fit
		avgRespTime sql.NullFloat64
		maxRespTime sql.NullFloat64
		avgTTFB     sql.NullFloat64
		maxTTFB     sql.NullFloat64
		countURL    = 0
		totalErr    = 0
	)

	// Query for page STATS
	page, queryErr := db.Query(
		`SELECT DISTINCT AVG(response_time), MAX(response_time), AVG(ttfb), MAX(ttfb), COUNT(id)
					FROM page WHERE Domain = ? AND crawl_id = ?`,
		domain, crawlID)
	if queryErr != nil {
		log.Println("QUERY PAGE ERR: ", queryErr)
	}
	defer page.Close()

	for page.Next() {
		scanErr := page.Scan(&avgRespTime, &maxRespTime, &avgTTFB, &maxTTFB, &countURL)
		if scanErr != nil {
			log.Println("PAGE.NEXT ERR: ", scanErr)
		}
	}

	// Query for outages/errors
	outage, outageErr := db.Query("SELECT COUNT(id) FROM outages WHERE Domain = ? AND crawl_id = ?",
		domain, crawlID)
	if outageErr != nil {
		log.Println("OUTAGE ERR: ", outageErr)
	}
	defer outage.Close()

	for outage.Next() {
		outErr := outage.Scan(&totalErr)
		if outErr != nil {
			log.Println("OUTAGE.NEXT ERR: ", outErr)
		}
	}

	// Create map for 404url keys and properties values
	four04Map := map[string][]string{}

	// Query for 404
	query404, queryErr := db.Query(
		`SELECT domain, page, datetime, referer FROM outages WHERE status_code = ? AND crawl_id = ?`,
		404, crawlID)
	if queryErr != nil {
		fmt.Println("QUERY404 PAGE ERR: ", queryErr)
	}
	defer query404.Close()

	for query404.Next() {
		props := Four04Props{}

		scan404Err := query404.Scan(&props.Domain, &props.Page, &props.TimeStamp, &props.Referer)
		if scan404Err != nil {
			fmt.Println("query404.NEXT ERR: ", scan404Err)
			break
		}
		// add 404 page to map with properties from query
		four04Map[props.Page.String] = []string{props.Page.String, props.Domain.String, props.Referer.String, props.TimeStamp.String()}
	}

	// Create map for all pages
	sitemap := map[string][]string{}

	// Query for pages, piece together sitemap
	querySitemap, queryErr := db.Query(
		`SELECT domain, url, datetime, status_code FROM page WHERE crawl_id = ?`, crawlID)
	if queryErr != nil {
		fmt.Println("Sitemap Query PAGE ERR: ", queryErr)
	}
	defer querySitemap.Close()

	for querySitemap.Next() {
		sitemapProps := SitemapProps{}

		scanSitemapErr := querySitemap.Scan(&sitemapProps.Domain, &sitemapProps.Page, &sitemapProps.TimeStamp, &sitemapProps.Status)
		if scanSitemapErr != nil {
			fmt.Println("Sitemap.NEXT ERR: ", scanSitemapErr)
			break
		}
		// add page to map with properties from query
		sitemap[sitemapProps.Page.String] = []string{sitemapProps.Page.String, sitemapProps.Domain.String, sitemapProps.Status.String, sitemapProps.TimeStamp.String(), ""}
	}

	// Query for all outage pages, for sitemap
	querySitemapOutage, queryErr := db.Query(
		`SELECT domain, page, datetime, status_code, referer FROM outages WHERE crawl_id = ?`, crawlID)
	if queryErr != nil {
		fmt.Println("Sitemap Outage Query PAGE ERR: ", queryErr)
	}
	defer querySitemapOutage.Close()

	for querySitemapOutage.Next() {
		sitemapOutageProps := SitemapProps{}

		scanSitemapErr := querySitemapOutage.Scan(&sitemapOutageProps.Domain, &sitemapOutageProps.Page, &sitemapOutageProps.TimeStamp,
			&sitemapOutageProps.Status, &sitemapOutageProps.Referer)
		if scanSitemapErr != nil {
			fmt.Println("Sitemap Outage.NEXT ERR: ", scanSitemapErr)
			break
		}
		// add outage page to map with properties from query
		sitemap[sitemapOutageProps.Page.String] = []string{sitemapOutageProps.Page.String, sitemapOutageProps.Domain.String,
			sitemapOutageProps.Status.String, sitemapOutageProps.TimeStamp.String(), sitemapOutageProps.Referer.String}
	}

	// Assign struct values to Domain
	statusMap[domain] = &DomainStats{
		// Get "real" float values from nullable types, and round them to 2 decimals before setting Domain values
		AvgRespTime: avgRespTime.Float64,
		MaxRespTime: maxRespTime.Float64,
		AvgTTFB:     avgTTFB.Float64,
		MaxTTFB:     maxTTFB.Float64,
		CountURL:    countURL,
		TotalErr:    totalErr,
		List404:     four04Map,
		Sitemap:     sitemap,
		TimeStamp:   time.Now(),
	}

	// Send email with solo status map
	return statusMap, crawlID
}

// DeleteSoloData removes rows from DB that were accumulated with a manual crawl
func DeleteSoloData(crawlID int) {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete rows from page table with crawl id
	_, deletePageErr := db.Exec(
		"DELETE FROM page WHERE crawl_id = ?", crawlID)
	if deletePageErr != nil {
		log.Println("Solo DELETE PAGE ERR: ", deletePageErr)
	}

	// Delete rows from crawl table
	_, deleteCrawlErr := db.Exec(
		"DELETE FROM crawl WHERE id = ?", crawlID)
	if deleteCrawlErr != nil {
		log.Println("Solo DELETE CRAWL ERR: ", deleteCrawlErr)
	}

	// Delete rows from outage table
	_, deleteOutagesErr := db.Exec(
		"DELETE FROM outages WHERE crawl_id = ?", crawlID)
	if deleteOutagesErr != nil {
		log.Println("Solo DELETE OUTAGES ERR: ", deleteOutagesErr)
	}
}
