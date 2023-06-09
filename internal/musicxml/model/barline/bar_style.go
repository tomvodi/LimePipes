package barline

import "encoding/xml"

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=kebab -type=Style

type Style uint8

const (
	None Style = iota
	Regular
	Dashed
	Dotted
	Heavy
	HeavyHeavy
	HeavyLight
	LightHeavy
	LightLight
	Short
	Tick
)

type BarStyle struct {
	XMLName xml.Name `xml:"bar-style"`
	Value   string   `xml:",chardata"`
}

func NewBarStyle(style Style) BarStyle {
	return BarStyle{
		XMLName: xml.Name{Local: "bar-style"},
		Value:   style.String(),
	}
}
