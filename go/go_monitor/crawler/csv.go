package crawler

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
)


// Create 404 csv for after all crawl
func Generate404CSV() int {
	file, err := os.OpenFile("404.csv", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	// Set header names
	heads := []string{"404 Link", "Domain", "Referred From", "Timestamp"}
	// Create string matrix for csv with headers input first
	strToCSV := [][]string{heads}
	// Get 404 map in list format
	pageMap := Get404List()

	// To store the keys in slice in sorted order
	var keys []string
	for k := range pageMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate through sorted 404 keys and append values to string matrix for csv consumption
	for _, item := range keys {
		strToCSV = append(strToCSV, pageMap[item])
	}

	// Create CSV writer, then write to CSV, then flush
	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strToCSV)
	csvWriter.Flush()

	return len(strToCSV)
}

// Create 404 csv for after all crawl
func GenerateCSV(filename string, headers []string, crawlList SearchUrlList) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	// Create string matrix for csv with headers input first
	strToCSV := [][]string{headers}

	// Iterate through crawlList appending url, status, comment
	for k, url := range crawlList {
		strToCSV = append(strToCSV, []string{k, strconv.Itoa(url.status), url.comment})
	}

	// Create CSV writer, then write to CSV, then flush
	csvWriter := csv.NewWriter(file)
	err = csvWriter.WriteAll(strToCSV)
	if err != nil {
		log.Println(err)
	}
	csvWriter.Flush()
}

// Generate csv for solo manual crawl
func Solo404CSV(input map[string][]string, domainName string) {

	fileName := domainName + "-404.csv"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	// Set header names
	heads := []string{"404 Link", "Domain", "Referred From", "Timestamp"}
	// Create string matrix for csv with headers input first
	strToCSV := [][]string{heads}

	// Iterate through 404 input and append values to string matrix for csv consumption
	for _, item := range input {
		strToCSV = append(strToCSV, item)
	}

	// Create CSV writer, then write to CSV, then flush
	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strToCSV)
	csvWriter.Flush()
}

// Generate csv for sitemap manual crawl
func SitemapCSV(input map[string][]string, domainName string) {

	fileName := domainName + "-sitemap.csv"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	// Set header names
	heads := []string{"Page", "Domain", "Status", "Timestamp", "Error On"}
	// Create string matrix for csv with headers input first
	strToCSV := [][]string{heads}

	// Iterate through sitemap input and append values to string matrix for csv consumption
	for _, item := range input {
		strToCSV = append(strToCSV, item)
	}

	// Create CSV writer, then write to CSV, then flush
	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strToCSV)
	csvWriter.Flush()
}

// Delete a given file
func deleteFile(file string) {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		os.Exit(1)
	}

	log.Println("Deleted File: " + file)
}