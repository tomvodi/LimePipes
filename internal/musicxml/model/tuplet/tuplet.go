package tuplet

import (
	"banduslib/internal/common/music_model/symbols/tuplet"
	"encoding/xml"
)

//go:generate go run github.com/dmarkham/enumer -transform=lower -type=Type
//go:generate go run github.com/dmarkham/enumer -transform=lower -type=NumberStyle
//go:generate go run github.com/dmarkham/enumer -transform=lower -type=Brackets

type Type uint8

const (
	Start Type = iota
	Stop
)

type NumberStyle uint8

const (
	Invisible NumberStyle = iota
	None
	Both
	Actual
)

type Brackets uint8

const (
	NoBrackets Brackets = iota
	No
	Yes
)

type Tuplet struct {
	XMLName    xml.Name `xml:"tuplet"`
	Type       string   `xml:"type,attr"`
	Bracket    string   `xml:"bracket,attr,omitempty"`
	ShowNumber string   `xml:"show-number,attr,omitempty"`
}

func NewTuplet(tType Type, brackets Brackets, showNmbr NumberStyle) *Tuplet {
	var bracketsVal string
	if brackets != NoBrackets {
		bracketsVal = brackets.String()
	}
	var numberVal string
	if showNmbr != Invisible {
		numberVal = showNmbr.String()
	}
	return &Tuplet{
		XMLName: xml.Name{
			Local: "tuplet",
		},
		Type:       tType.String(),
		Bracket:    bracketsVal,
		ShowNumber: numberVal,
	}
}

func FromMusicModel(tpl *tuplet.Tuplet) *Tuplet {
	tType := Start
	if tpl.BoundaryType == tuplet.End {
		tType = Stop
	}
	var nrStyle NumberStyle
	var brackets Brackets
	if tType == Start {
		nrStyle = Both
		if tpl.VisibleNotes == 3 && tpl.PlayedNotes == 2 {
			nrStyle = Actual
		}
		brackets = Yes
	}

	return NewTuplet(tType, brackets, nrStyle)
}
