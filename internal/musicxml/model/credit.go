package model

import "encoding/xml"

const (
	CreditTypePageNumber string = "page number"
	CreditTypeTitle             = "title"
	CreditTypeSubtitle          = "subtitle"
	CreditTypeComposer          = "composer"
	CreditTypeArranger          = "arranger"
	CreditTypeLyricist          = "lyricist"
	CreditTypeRights            = "rights"
	CreditTypePartName          = "part name"
)

type Credit struct {
	XMLName xml.Name `xml:"credit"`
	Page    uint     `xml:"page,attr"`
	Type    string   `xml:"credit-type"`
	Words   CreditWords
}
