package main

import (
	"./crawler"
	"./router"

	"time"

	"gopkg.in/robfig/cron.v2"
)

func main() {
	/* Cron Setup */
	c := cron.New()
	// Run 'All Crawl'
	c.AddFunc("0 0 4 * * *", func() {
		go crawler.RunAllCrawl()
	})
	// Update (weekly/monthly) cache for front end every 2 hours
	c.AddFunc("0 0 */2 * * *", func() {
		go crawler.StatusReportUpdateWeekly()
		go crawler.StatusReportUpdateMonthly()
	})
	// Send Weekly Email
	c.AddFunc("@weekly", func() {
		go crawler.WeeklyEmailReport()

		// Run 404 email list once a week
		go crawler.NewEmail404()
	})
	// Every month delete old data from DB
	c.AddFunc("@monthly", func() {
		go crawler.DeleteOldData()
	})
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
	// Create channel to receive value from crawler
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
		// 10sec between homepage crawls
		time.Sleep(10 * time.Second)
	}
}
