package helper

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

type tuneFix struct {
}

func (t *tuneFix) Fix(muMo music_model.MusicModel) {
	for _, tune := range muMo {
		fixComposerArranger(tune)
		fixComposerTrad(tune)
		removeTimeSigFromTuneType(tune)
		removeSpecialCharsFromTuneType(tune)
		trimSpaces(tune)
		fixTitle(tune)
		capitalizeTuneType(tune)
	}
}

func fixComposerArranger(tune *music_model.Tune) {
	if tune.Composer == "" {
		return
	}

	//parts := strings.Split(tune.Composer, ",-([")
	regX := regexp.MustCompile(`(?i)arranged|arrangement|arr[.:/ ]+`)
	parts := regX.Split(tune.Composer, -1)
	if len(parts) == 2 {
		arranger := fixComposerArrangerField(parts[1])
		sep := regexp.MustCompile(`,|-`)
		arrangerSplit := sep.Split(arranger, -1)

		// tuneTitle after arranger
		if len(arrangerSplit) > 1 {
			tune.Arranger = fixComposerArrangerField(arrangerSplit[0])
			tune.Composer = fixComposerArrangerField(arrangerSplit[1])
			return
		}

		// tuneTitle before arranger
		if strings.TrimSpace(parts[0]) != "" {
			tune.Composer = fixComposerArrangerField(parts[0])
			tune.Arranger = arranger
		} else {
			// only arranger in tuneTitle field
			tune.Composer = ""
			tune.Arranger = arranger
		}
	}
}

func fixComposerTrad(tune *music_model.Tune) {
	regX := regexp.MustCompile(`(?i)^trad\.?$`)
	trimmedComposer := strings.TrimSpace(tune.Composer)
	if regX.MatchString(trimmedComposer) {
		tune.Composer = "Traditional"
	}
}

func removeTimeSigFromTuneType(tune *music_model.Tune) {
	trimmedType := strings.TrimSpace(tune.Type)
	regX := regexp.MustCompile(`\d+/\d+`)
	typeWithoutTimesig := regX.ReplaceAllString(trimmedType, "")
	typeWithoutTimesig = strings.TrimSpace(typeWithoutTimesig)
	tune.Type = typeWithoutTimesig
}

func removeSpecialCharsFromTuneType(tune *music_model.Tune) {
	trimmedType := strings.TrimSpace(tune.Type)
	trimmedType = strings.Trim(trimmedType, ".:/-|")
	tune.Type = trimmedType
}

func trimSpaces(tune *music_model.Tune) {
	tune.Title = strings.TrimSpace(tune.Title)
	tune.Composer = strings.TrimSpace(tune.Composer)
	tune.Arranger = strings.TrimSpace(tune.Arranger)
	tune.Type = strings.TrimSpace(tune.Type)
}

func capitalizeTuneType(tune *music_model.Tune) {
	trimmedType := strings.TrimSpace(tune.Type)
	caser := cases.Title(language.English)
	tune.Type = caser.String(trimmedType)
}

func fixComposerArrangerField(arr string) string {
	arranger := strings.TrimSpace(arr)
	arranger = strings.Trim(arranger, ".:/-[](),")
	arranger = strings.Replace(arranger, "by", "", -1)
	arranger = strings.TrimSpace(arranger)

	return arranger
}

func fixTitle(tune *music_model.Tune) {
	tune.Title = strings.Replace(tune.Title, "_", " ", -1)

	caser := cases.Title(language.English)
	tune.Title = caser.String(tune.Title)
}

func NewTuneFixer() interfaces.TuneFixer {
	return &tuneFix{}
}
