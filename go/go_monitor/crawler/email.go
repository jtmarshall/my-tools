package crawler

import (
	"../creds"
	"../email"
	"os"

	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/mail"
	"net/smtp"
	"sort"
	"strings"
	"time"
)

const (
	// AWS Config
	emailHost     = creds.EmailHost
	emailHostUser = creds.EmailHostUser
	emailHostPass = creds.EmailHostPass
	emailPort     = creds.EmailPort
)

var (
	// Emails
	soloEmail    = creds.SoloEmail
	fromaddr     = creds.FromAddr
	toaddr       = creds.ToAddr
	emailList404 = creds.EmailList404
)

// NewEmail404 allows for csv attachment
func NewEmail404() {
	lockFile := "lock404.txt"

	// Get lock file stats
	info, err := os.Stat(lockFile)
	// Check for lock file before we send emails. (If it doesn't exist or file is 24hours run email func)
	if err != nil || time.Now().Sub(info.ModTime()) > (48 * time.Hour) {
		// Have to double check if statement for mod time on file so we can validate err also
		if os.IsNotExist(err) || time.Now().Sub(info.ModTime()) > (48 * time.Hour) {
			// Lock file doesn't exist; create the lock file
			f, fErr := os.Create(lockFile)
			defer f.Close()
			if fErr != nil {
				os.Exit(1)
			}

			// START: 404 email creation
			fileName := "404.csv"

			// Generate/overwrite csv
			errLength := Generate404CSV()
			fmt.Println("404 Count: ", errLength)

			subject := "Monitor: 404 List"

			m := email.NewMessage(subject, "See 404's Attached. \r\n You can view live daily updated 404's here: \r\n http://monitor.acadiadevelopment.com/ \r\n" +
				"Time: " + time.Now().Format(time.Stamp))
			// compose the message
			m.From = mail.Address{Name: "Acadia Monitoring", Address: fromaddr}
			m.To = emailList404

			// if there is more than just the headers in the csv
			if errLength > 1 {
				// add attachments
				if err := m.Attach(fileName); err != nil {
					log.Fatal(err)
				}
			} else {
				// if no 404's edit body txt and don't attach empty csv
				m.Body = "No new 404's yay! \r\n You can view live daily updated 404's here: \r\n http://monitor.acadiadevelopment.com/ \r\n" +
					"Time: " + time.Now().Format(time.Stamp)
			}

			// send it
			auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)
			if err := email.Send(emailHost+emailPort, auth, m); err != nil {
				log.Fatal(err)
			}

			// Delete the 404 file after were finished
			deleteFile(fileName)

			// Remove lock file after email send
			deleteFile(lockFile)
		}
	}
}

// Email 404 csv for manual crawl
func SoloEmail404(inEmail []string, domain string, errCount int) {
	subject := domain + " 404 List"

	body := ""

	if errCount > 0 {
		body = "404 CSV Attached \r\n" + "Time: " + time.Now().Format(time.Stamp)
	} else {
		body = "No 404 Found \r\n" + "Time: " + time.Now().Format(time.Stamp)
	}

	// compose the message
	m := email.NewMessage(subject, body)
	m.From = mail.Address{Name: "Acadia Monitoring", Address: fromaddr}
	m.To = inEmail

	fileName := domain + "-404.csv"
	// If there are 404's attach the csv
	if errCount > 0 {
		// add attachments
		if err := m.Attach(fileName); err != nil {
			log.Fatal(err)
		}
	}

	// send it
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)
	if err := email.Send(emailHost+emailPort, auth, m); err != nil {
		log.Fatal(err)
	}

	// Delete the file after were finished
	deleteFile(fileName)
}

// Email sitemap csv for manual crawl
func SitemapEmail(inEmail []string, domain string) {
	subject := domain + " Sitemap"

	body := "Sitemap attached \r\n" + "Time: " + time.Now().Format(time.Stamp)

	// compose the message
	m := email.NewMessage(subject, body)
	m.From = mail.Address{Name: "Acadia Monitoring", Address: fromaddr}
	m.To = inEmail

	fileName := domain + "-sitemap.csv"

	// Attach csv
	if err := m.Attach(fileName); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)
	if err := email.Send(emailHost+emailPort, auth, m); err != nil {
		log.Fatal(err)
	}

	// Delete the file after were finished
	deleteFile(fileName)
}

// *Deprecated* EmailDaily404 sends out 404 list everyday
func EmailDaily404() {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)

	currentTime := time.Now().Format(time.Stamp)

	fofSet := Get404DomainList()

	count404 := 0
	// To store the keys in slice in sorted order
	var keys []string
	for k := range fofSet {
		keys = append(keys, k)
		count404 += len(fofSet[k])
	}
	sort.Strings(keys)

	// templateData struct to pass sorted keys and four04Set to html template
	templateData := struct {
		SortedKeys []string
		ErrSet     map[string][]Four04Props
		ErrCount   int
	}{keys, fofSet, count404}

	// Parse template that will be execute with passed in data struct
	tmpl, _ := template.ParseFiles("templates/email404.html")

	// byte buffer to read in template with data
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, templateData)

	toHead := "To: " + strings.Join(emailList404, ",") + "\r\n"
	subject := "Subject: Monitor: 404 List " + currentTime + " \n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Setup message to send
	msg := []byte(toHead + subject + mime + "\n")
	// append byte message to byte buffer
	msg = append(msg, buf.Bytes()...)

	// Connect to the server, authenticate, set the sender, recipient, and pass in the msg
	emailErr := smtp.SendMail(emailHost+emailPort, auth, fromaddr, emailList404, msg)
	if emailErr != nil {
		log.Fatal(emailErr)
	}
}

// EmailAlert notifies email list of error; without Domain needed
func EmailAlert(url string, code int) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)

	// Setup message to send
	msg := []byte(
		"To: " + strings.Join(soloEmail, ",") + "\r\n" +
			"Subject: Monitor: " + fmt.Sprint(code) + " Alert\r\n" + "\r\n" +
			"Status: " + fmt.Sprint(code) + "\r\n" +
			"On page: " + url + "\r\n" +
			"Time: " + time.Now().Format(time.Stamp))

	// Connect to the server, authenticate, set the sender, recipient, and pass in the msg
	err := smtp.SendMail(emailHost+emailPort, auth, fromaddr, soloEmail, msg)
	if err != nil {
		log.Fatal(err)
	}
}

// BackOnline notifies email list of a Domain recovering from 5xx status
func BackOnline(domain string, url string, code int) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)

	// Setup message to send
	msg := []byte(
		"To: " + strings.Join(toaddr, ",") + "\r\n" +
			"Subject: Monitor: Domain Recovered\r\n" + "\r\n" +
			domain + " is back online. \r\n" +
			"Status: " + fmt.Sprint(code) + "\r\n" +
			"On page: " + url + "\r\n" +
			"Time: " + time.Now().Format(time.Stamp))

	// Connect to the server, authenticate, set the sender, recipient, and pass in the msg
	err := smtp.SendMail(emailHost+emailPort, auth, fromaddr, toaddr, msg)
	if err != nil {
		log.Fatal(err)
	}
}

// WeeklyEmail sends out automatic email every week
func WeeklyEmail(summSet map[string]*DomainStats, domainList []string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)

	toHead := "To: " + strings.Join(toaddr, ",") + "\r\n"
	subject := "Subject: Monitor: Weekly Report \n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Setup message to send
	msg := []byte(toHead + subject + mime + "\n")
	formatTxt := []byte(`<table border="0" cellpadding="1" cellspacing="0" height="100%" id="bodyTable">`)
	msg = append(msg, formatTxt...)

	// We must iterate through the sorted list of domains for the report to be alphabetical;
	// 	rather than just through the given map because map iteration order is intentionally undefined.
	for _, domain := range domainList {

		// Create text byte chunk for Domain
		domainTxt := []byte(
			`
			<tr>
				<td align="left" valign="top">
					<table border="0" cellpadding="4" cellspacing="0" width="400" id="emailContainer">
						<tr>
							<h3> ` + domain + ` </h3>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Avg Response: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].AvgRespTime) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Max Response: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].MaxRespTime) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Avg TTFB: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].AvgTTFB) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Max TTFB: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].MaxTTFB) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total Crawled: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].CountURL) + `
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total Errors: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].TotalErr) + `
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total 404: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].List404) + `
							</td>
						</tr>
						<tr>
							<td> </td> <td> </td>
						</tr>
					</table>
				</td>
			</tr>
			`)

		// Append Domain text chunk to msg
		msg = append(msg, domainTxt...)
	}

	endFormatTxt := []byte(`</table>`)
	msg = append(msg, endFormatTxt...)

	// Connect to the server, authenticate, set the sender, recipient, and pass in the msg
	err := smtp.SendMail(emailHost+emailPort, auth, fromaddr, toaddr, msg)
	if err != nil {
		log.Fatal(err)
	}
}

// SoloCrawlEmail sends out automatic email every week
func SoloCrawlEmail(summSet map[string]*DomainStats, domain string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailHostUser, emailHostPass, emailHost)

	toHead := "To: " + strings.Join(toaddr, ",") + "\r\n"
	subject := "Subject: Monitor: " + domain + " Report \n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Setup message to send
	msg := []byte(toHead + subject + mime + "\n")
	formatTxt := []byte(`<table border="0" cellpadding="1" cellspacing="0" height="100%" id="bodyTable">`)
	msg = append(msg, formatTxt...)

	// Create text byte chunk for Domain
	domainTxt := []byte(
		`
			<tr>
				<td align="left" valign="top">
					<table border="0" cellpadding="4" cellspacing="0" width="400" id="emailContainer">
						<tr>
							<h3> ` + domain + ` </h3>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Avg Response: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].AvgRespTime) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Max Response: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].MaxRespTime) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Avg TTFB: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].AvgTTFB) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Max TTFB: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprintf("%.2f", summSet[domain].MaxTTFB) + `<i>ms</i>
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total Crawled: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].CountURL) + `
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total Errors: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].TotalErr) + `
							</td>
						</tr>
						<tr>
							<td align="left" valign="top">
								<b>Total 404: </b>
							</td>
							<td align="left" valign="top">
								` + fmt.Sprint(summSet[domain].List404) + `
							</td>
						</tr>
						<tr>
							<td> </td> <td> </td>
						</tr>
					</table>
				</td>
			</tr>
			`)

	// Append Domain text chunk to msg
	msg = append(msg, domainTxt...)

	endFormatTxt := []byte(`</table>`)
	msg = append(msg, endFormatTxt...)

	// Connect to the server, authenticate, set the sender, recipient, and pass in the msg
	err := smtp.SendMail(emailHost+emailPort, auth, fromaddr, toaddr, msg)
	if err != nil {
		log.Fatal(err)
	}
}
