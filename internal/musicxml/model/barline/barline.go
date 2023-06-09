package barline

import (
	"banduslib/internal/common/music_model/barline"
	"encoding/xml"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
)

type Barline struct {
	XMLName  xml.Name `xml:"barline"`
	Location string   `xml:"location,attr"`
	Style    BarStyle `xml:"bar-style"`
	Repeat   *Repeat  `xml:"repeat,omitempty"`
}

func FromMusicModel(muMoBar *barline.Barline, loc Location) Barline {
	barL := Barline{
		XMLName: xml.Name{
			Local: "barline",
		},
		Location: loc.String(),
		Style:    NewBarStyle(convertBarlineType(muMoBar.Type)),
	}

	if muMoBar.Timeline == barline.Repeat {
		dir := Forward
		if loc == Right {
			dir = Backward
		}
		barL.Repeat = NewRepeat(dir)
	}

	return barL
}

func convertBarlineType(barlineType barline.BarlineType) Style {
	kebapC := strcase.KebabCase(barlineType.String())
	style, err := StyleString(kebapC)
	if err != nil {
		log.Error().Err(err).Msg("failed converting barline type (MusicModel) " +
			"to barline tyle (musicxml)")
		return None
	}
	return style
}
