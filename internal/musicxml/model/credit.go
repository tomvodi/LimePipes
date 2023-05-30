package model

import "encoding/xml"

const (
	CreditTypePageNumber string = "page number"
	CreditTypeTitle      string = "title"
	CreditTypeSubtitle   string = "subtitle"
	CreditTypeComposer   string = "composer"
	CreditTypeArranger   string = "arranger"
	CreditTypeLyricist   string = "lyricist"
	CreditTypeRights     string = "rights"
	CreditTypePartName   string = "part name"
)

type Credit struct {
	XMLName xml.Name `xml:"credit"`
	Page    uint     `xml:"page,attr"`
	Type    string   `xml:"credit-type"`
	Words   CreditWords
}
