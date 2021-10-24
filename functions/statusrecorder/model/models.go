package model

type Site struct {
	Id           string
	Name         string
	SiteProvider string
	Address      string
	Records      []Record
	Data         string
}

type Record struct {
	SiteId    string
	SiteName  string
	Available int
	Total     int
}

type CCSCharger struct {
	Identifier string
	Name       string
	Capacity   int64
}
