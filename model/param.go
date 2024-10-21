package model

type SearchNameCodeParam struct {
	Code        string
	Name        string
	FilterCrawl bool
}
type SearchCrawlParam struct {
	Code         string
	ReportPeriod string
	CrawlType    string
	UniqueID     int64
}

type SearchProfileParam struct {
}
