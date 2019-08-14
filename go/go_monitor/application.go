package main

import (
	"./crawler"
	"./router"
	"log"
	"time"

	"gopkg.in/robfig/cron.v2"
)

func main() {
	/* Cron Setup */
	c := cron.New()
	// Run 'All Crawl'
	_, err := c.AddFunc("0 0 4 * * *", func() {
		go crawler.RunAllCrawl()
	})
	if err != nil {
		log.Println(err)
	}
	// Update (weekly/monthly) cache for front end every 2 hours
	_, err = c.AddFunc("0 0 */2 * * *", func() {
		go crawler.StatusReportUpdateWeekly()
		go crawler.StatusReportUpdateMonthly()
	})
	if err != nil {
		log.Println(err)
	}
	// Send Weekly Email 6am
	_, err = c.AddFunc("0 0 6 * * 0", func() {
		go crawler.NewEmail404()
		go crawler.WeeklyEmailReport()
	})
	if err != nil {
		log.Println(err)
	}
	// Every month delete old data from DB
	_, err = c.AddFunc("@monthly", func() {
		go crawler.DeleteOldData()
	})
	if err != nil {
		log.Println(err)
	}
	c.Start()
	/* END Cron Setup */

	// Fill monitor status request objects (weekly/monthly) on startup
	go crawler.StatusReportUpdateWeekly()
	go crawler.StatusReportUpdateMonthly()

	// Spin up server
	go router.Run()

	// Start error handler daemon
	go crawler.ErrorDaemon()

	/* Main "Daemon Process" */
	//Create channel to receive value from crawler
	ch := make(chan bool)
	check := true
	for {
		// If true then we start the crawl
		if check {
			// Set flag false
			check = false
			go crawler.RunHomeCrawl(ch)
			// After crawl finishes reset flag back to true so we can spin up the crawl again
			check = <-ch
		}
		// 30sec between homepage crawls
		time.Sleep(30 * time.Second)
	}
}