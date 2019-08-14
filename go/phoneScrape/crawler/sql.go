package crawler

import (
	"database/sql"
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

type FacilityLabel struct {
	FacilityName string
	FacilityType string
}

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