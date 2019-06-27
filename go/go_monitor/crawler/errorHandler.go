package crawler

import (
	"net/http"
	"time"
	"log"
)

var errorMap = make(map[string]int)

// ErrorDaemon Starts the error handler service
// where we double check errors that come up during crawl loops before sending out emails
func ErrorDaemon() {
	// error loop
	for {
		// Pause 30 sec start of each loop to allow error queue to fill
		time.Sleep(30 * time.Second)

		// Check if errors need to be handled
		if len(errorMap) > 0 {
			for errURL := range errorMap {
				// get status of current error url
				urlStatus := checkURLStatus(errURL)

				// If 500 error again, then send error email
				if urlStatus >= 500 && errorMap[errURL] != urlStatus {
					// update error map to have latest status code
					errorMap[errURL] = urlStatus
					// then send out the error email
					EmailAlert(errURL, urlStatus)
				}

				// If status recovers
				if urlStatus == 200 {
					// Remove errUrl from the map
					delete(errorMap, errURL)
					// Only send recovery email if last status was 500
					if errorMap[errURL] >= 500 {
						EmailRecoveryAlert(errURL, urlStatus)
					}
				}
			}
		} else {
			// No errors continue to sleep
			continue
		}
	}
}

// AddError allow files to insert errors into the map
func DaemonAddError(err string) {
	// add into map, set starting code to 0 so we know it's new
	errorMap[err] = 0
	log.Println("Adding URL to ErrorDaemon:", err)
}

// re-parse url
func checkURLStatus(url string) int {
	// url passed in already has "https://" or http whatever
	resp, err := http.Get("https://" + url)
	if err != nil {
		log.Println(err)
		return 0
	}
	if err != nil {
		log.Println("Daemon Check GET Request ERR:", err)
		// retry without https
		resp, err = http.Get("http://" + url)
		if err != nil {
			log.Println("Daemon ERR2:", err)
		}
	}
	defer resp.Body.Close()

	return resp.StatusCode
}