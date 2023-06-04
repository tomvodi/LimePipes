package model

import "encoding/xml"

type ScorePartList struct {
	XMLName xml.Name    `xml:"part-list"`
	Parts   []ScorePart `xml:"score-part"`
}
