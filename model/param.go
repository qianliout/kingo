package model

type SearchNameCodeParam struct {
	Code        string
	Name        string
	FilterCrawl bool
}
type SearchCrawlParam struct {
	Code      string
	Year      string
	CrawlType string
}

type SearchProfileParam struct {
}
