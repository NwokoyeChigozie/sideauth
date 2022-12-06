package config

type ServerConfiguration struct {
	Port                       string
	Secret                     string
	AccessTokenExpireDuration  int
	RefreshTokenExpireDuration int
}

type Microservices struct {
	Admin        string
	Auth         string
	Boilerplate  string
	Cron         string
	Feedback     string
	Internaldocs string
	Notification string
	Payment      string
	Productlink  string
	Referral     string
	Reminders    string
	Roles        string
	Subscription string
	Transactions string
	Upload       string
	Verification string
	Widget       string
}
