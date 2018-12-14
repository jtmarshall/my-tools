package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type DateRange struct {
	From time.Time `json:"from,omitempty"`
	To   time.Time `json:"to,omitempty"`
}

type PageView struct {
	URL         string          `json:"url,omitempty"`
	Date        time.Time       `json:"date,omitempty"`
	HasCall     bool            `json:"hasCall,omitempty"`
	TOS         int             `json:"tos,omitempty"`
	Conversions map[string]bool `json:"conversions,omitempty"`
}

// Checks if conversions were made for a given pageview
func (pv *PageView) CheckConversions() {
	// use pageview properties to check against conversion type
	if pv.HasCall {
		pv.Conversions["Call"] = true
	}
	if pv.TOS > 120 {
		pv.Conversions["2min"] = true
	}
	if pv.TOS > 360 {
		pv.Conversions["5min"] = true
	}
}

type SKU struct {
	Number       string `json:"number,omitempty"`
	Network      string `json:"network,omitempty"`
	Targeting    string `json:"targeting,omitempty"`
	Format       string `json:"format,omitempty"`
	Message      string `json:"message,omitempty"`
	AgeRange     string `json:"ageRange,omitempty"`
	Ethnicity    string `json:"ethnicity,omitempty"`
	FamilyRole   string `json:"familyRole,omitempty"`
	Gender       string `json:"gender,omitempty"`
	Income       string `json:"income,omitempty"`
	Interests    string `json:"interests,omitempty"`
	Language     string `json:"language,omitempty"`
	Education    string `json:"education,omitempty"`
	Occupation   string `json:"occupation,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	Religion     string `json:"religion,omitempty"`
}

type Session struct {
	DateRange  *DateRange `json:"dateRange,omitempty"`
	PassportID string     `json:"passportID,omitempty"`
	IP         string     `json:"ip,omitempty"`
	Domain     string     `json:"domain,omitempty"`
	FirstTouch string     `json:"firstTouch,omitempty"`
	LastTouch  string     `json:"lastTouch,omitempty"`
	CallTouch  string     `json:"callTouch,omitempty"`
	DeviceType string     `json:"deviceType,omitempty"`
	Trail      []PageView `json:"trail,omitempty"`
	Geo        string     `json:"geo,omitempty"`
	Browser    string     `json:"browser,omitempty"`
	OS         string     `json:"os,omitempty"`
	SKU        *SKU       `json:"sku,omitempty"`
}

type YakData struct {
	ID        string     `json:"id,omitempty"`
	Facility  string     `json:"facility,omitempty"`
	DateRange *DateRange `json:"dateRange,omitempty"`
	Sessions  []Session  `json:"sessions,omitempty"`
}

var data YakData

func main() {
	// Generate random data on startup
	GenerateYakData()

	// Create router
	router := mux.NewRouter()

	// Endpoints for router
	router.HandleFunc("/retrieve", GetData).Methods("GET")

	// Serve router
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(data)
}

// Create some random data for testing
func GenerateYakData() {
	// Create date frame
	randate := &DateRange{
		To:   time.Now(),
		From: time.Now().Add(-7),
	}

	// Create random YakData
	randYakData := YakData{
		ID:        strconv.Itoa(rand.Intn(300)) + strconv.Itoa(rand.Intn(2000)) + strconv.Itoa(rand.Intn(900)),
		DateRange: randate,
		Sessions:  make([]Session, 0),
	}

	// Loop creating random session data
	for i := 0; i < 7; i++ {
		randata := Session{
			DateRange: randate,
			IP:        strconv.Itoa(rand.Intn(200)) + "." + strconv.Itoa(rand.Intn(200)) + "." + strconv.Itoa(rand.Intn(200)) + "." + strconv.Itoa(rand.Intn(200)),
			PassportID: strconv.Itoa(rand.Intn(300)) + strconv.Itoa(rand.Intn(2000)) + strconv.Itoa(rand.Intn(900)),
			Domain: "",
			FirstTouch: "/",
			LastTouch: "/about/contact",
			CallTouch: "/about/location",
			SKU: &SKU{
				Number: "A"+strconv.Itoa(rand.Intn(6)) + "B"+strconv.Itoa(rand.Intn(20)) + "C"+strconv.Itoa(rand.Intn(20)) +
					"D"+strconv.Itoa(rand.Intn(21)) + "E"+strconv.Itoa(rand.Intn(10)) + "F"+strconv.Itoa(rand.Intn(9)) +
					"G"+strconv.Itoa(rand.Intn(13)) + "H"+strconv.Itoa(rand.Intn(4)) + "I"+strconv.Itoa(rand.Intn(12)) +
					"J"+strconv.Itoa(rand.Intn(21)) + "K"+strconv.Itoa(rand.Intn(5)) + "L"+strconv.Itoa(rand.Intn(11)) +
					"M"+strconv.Itoa(rand.Intn(17)) + "N"+strconv.Itoa(rand.Intn(6)) + "O"+strconv.Itoa(rand.Intn(11)),
			},
		}

		// Push into random yak data sessions
		randYakData.Sessions = append(randYakData.Sessions, randata)
	}

	// Update returned data object with random data
	data = randYakData
}
