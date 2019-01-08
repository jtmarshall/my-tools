package crawler

import (
	"net/http"
	"time"
	"log"
)

var errorMap = make(map[string]bool)

// ErrorDaemon Starts the error handler service
// where we double check errors that come up during crawl loops before sending out emails
func ErrorDaemon() {
	// error loop
	for {
		// Pause 30 sec start of each loop to allow error queue to fill
		time.Sleep(30 * time.Second)

		// Check if errors need to be handled
		if len(errorMap) > 0 {
			log.Println(errorMap)
			for errURL := range errorMap {
				// get status of current error url
				urlStatus := checkURLStatus(errURL)

				// If error again, then send error email
				if urlStatus >= 500 {
					EmailAlert(errURL, urlStatus)
				}

				// Then remove errUrl from the map
				delete(errorMap, errURL)
			}
		} else {
			// No errors continue to sleep
			continue
		}
		// fmt.Println(errorMap)
	}
}

// AddError allow files to insert errors into the map
func DaemonAddError(err string) {
	// add into map
	errorMap[err] = true
	log.Println("Add error url:", err)
}

// re-parse url
func checkURLStatus(url string) int {
	// url passed in already has "https://" or http whatever
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}