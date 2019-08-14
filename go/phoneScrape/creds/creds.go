package creds

const (
	// DB connection creds
	DbHost          = "monitoring-dev1.clkd4diaelqh.us-east-1.rds.amazonaws.com"
	DbName          = "monitor"
	DbUser          = "acadia"
	DbPass          = "Z5NeeAZa49XeMfrm" //OLD: G803$h#dOAu
	DbConnectString = DbUser + ":" + DbPass + "@tcp(" + DbHost + ")/" + DbName

	// AWS Email Config
	EmailHost     = "email-smtp.us-east-1.amazonaws.com"
	EmailHostUser = "AKIAICWC7VIKR4B2WI2A"
	EmailHostPass = "AnkURD2UOvdpVQbOVetcA+CEHSpIvy79z2x8nOeRi2/Q"
	EmailPort     = ":587"

	// Basic Auth
	AuthUser = "AcadiaMarketing"
	AuthPass = "mtVi39pXaz"
)

var (
	// Emails
	SoloEmail    = []string{"jordan.marshall@acadiahealthcare.com"}
	FromAddr     = "Monitor@acadiahealthcare.com"
	ToAddr       = []string{"jordan.marshall@acadiahealthcare.com", "eric.austin@acadiahealthcare.com", "edwin.orjales@acadiahealthcare.com"}
	EmailList404 = []string{
		"jordan.marshall@acadiahealthcare.com",
		"eric.austin@acadiahealthcare.com",
		"edwin.orjales@acadiahealthcare.com",
		"Allison.Isaacs@acadiahealthcare.com",
		"Claire.Baldwin@acadiahealthcare.com",
		"Matthew.Fung-A-Fat@acadiahealthcare.com",
		"nolan.omalley@acadiahealthcare.com",
		"Ryan.Beagan@acadiahealthcare.com",
		"ryan.milyard@acadiahealthcare.com",
		"Courtney.Wainner@acadiahealthcare.com",
		"Allison.Baioni@acadiahealthcare.com",
		"taylor.wood@acadiahealthcare.com"}
)
