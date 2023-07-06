package timeline_normalizer

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/barline"
)

func tuneToParts(tune *music_model.Tune) PartedTune {
	var partedTune PartedTune

	currPart := &MusicPart{}
	//currChunk := MusicChunk{}
	for _, measure := range tune.Measures {
		if measure.LeftBarline != nil && measure.LeftBarline.IsHeavy() {
			if measure.LeftBarline.Timeline == barline.Repeat {
				currPart.WithRepeat = true
			}
		}
		if measure.RightBarline != nil && measure.RightBarline.IsHeavy() {
			partedTune.Parts = append(partedTune.Parts, *currPart)
			currPart = &MusicPart{}
		}
		//if measure.HasTimeline() {
		//
		//} else {
		//	currChunk.Measures = append(currChunk.Measures, measure)
		//}

		// if measure has timeline
		/*		for _, symbol := range measure.Symbols {

				}*/
	}

	return partedTune
}
