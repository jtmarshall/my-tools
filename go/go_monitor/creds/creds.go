package creds

const (
	// DB connection creds
	DbHost          = "dbhost"
	DbName          = "dbname"
	DbUser          = "user"
	DbPass          = "%" + "password"
	DbConnectString = DbUser + ":" + DbPass + "@tcp(" + DbHost + ")/" + DbName

	// AWS Email Config
	EmailHost     = "emailhost"
	EmailHostUser = "emailuser"
	EmailHostPass = "emailpass"
	EmailPort     = ":port"
)

var (
	// Emails
	SoloEmail = []string{"tester@test.com"}
	FromAddr  = "from@test.com"
	ToAddr    = []string{"toAddr@test.com"}
)
