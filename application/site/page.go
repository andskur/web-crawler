package site

// Page represent web-site page structure with own URL
// and slice of the links - pointers to other pages
type Page struct {
	Url        *Url    `json:"url" xml:"url"`
	TotalLinks int     `json:"total,omitempty" xml:"total,omitempty"`
	Links      []*Page `json:"links,omitempty" xml:"links>page,omitempty"`
}
