package part_iterator

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/barline"
	"banduslib/internal/interfaces"
	"github.com/jinzhu/copier"
)

type partIterator struct {
	parts      []*music_model.MusicPart
	currentIdx int
}

func (p *partIterator) Nr(nr int) *music_model.MusicPart {
	if len(p.parts) < nr || nr <= 0 {
		return nil
	}

	return p.parts[nr-1]
}

func (p *partIterator) Count() int {
	return len(p.parts)
}

func (p *partIterator) HasNext() bool {
	if len(p.parts) == 0 {
		return false
	}
	if p.currentIdx+1 >= len(p.parts) {
		return false
	}

	return true
}

func (p *partIterator) GetNext() *music_model.MusicPart {
	if !p.HasNext() {
		return nil
	}
	p.currentIdx++
	return p.parts[p.currentIdx]
}

func splitTuneIntoParts(tune *music_model.Tune) (parts []*music_model.MusicPart) {
	currPart := &music_model.MusicPart{}
	var partMeasures []*music_model.Measure
	for i, measure := range tune.Measures {

		timelineMeasures := splitMeasureWithTimeline(measure)
		for _, timelineMeasure := range timelineMeasures {
			partMeasures = append(partMeasures, timelineMeasure)
		}

		if measure.LeftBarline != nil && measure.LeftBarline.IsHeavy() {
			if measure.LeftBarline.Timeline == barline.Repeat {
				currPart.WithRepeat = true
			}
		}

		if measure.RightBarline != nil && measure.RightBarline.IsHeavy() {
			// Part finished
			currPart.Measures = partMeasures
			parts = append(parts, currPart)
			if i+1 < len(tune.Measures) {
				currPart = &music_model.MusicPart{}
				partMeasures = []*music_model.Measure{}
			}
		}
	}

	return parts
}

// splitMeasureWithTimeline splits a given measure into multiple measures when there is more
// than one timeline in it or the timeline is surrounded by other symbols.
func splitMeasureWithTimeline(
	measure *music_model.Measure,
) (measures []*music_model.Measure) {
	if !measure.ContainsTimelineSymbols() ||
		measure.IsOneTimeline() {
		measures = append(measures, measure)
		return measures
	}

	currentMeasure := &music_model.Measure{}

	// first measure is like the original one
	// but without symbols and no right barline
	copier.Copy(currentMeasure, measure)
	currentMeasure.Symbols = []*music_model.Symbol{}
	currentMeasure.RightBarline = nil

	for _, symbol := range measure.Symbols {

		// When time line start comes after some other symbols
		if symbol.IsTimelineStart() {
			if currentMeasure.HasSymbols() {
				measures = append(measures, currentMeasure)
				currentMeasure = &music_model.Measure{}
			}

			currentMeasure.Symbols = append(currentMeasure.Symbols, symbol)
			continue
		}

		// When other symbols come after time line end
		if symbol.IsTimelineEnd() {
			currentMeasure.Symbols = append(currentMeasure.Symbols, symbol)
			if currentMeasure.HasSymbols() {
				measures = append(measures, currentMeasure)
				currentMeasure = &music_model.Measure{}
			}
			continue
		}

		currentMeasure.Symbols = append(currentMeasure.Symbols, symbol)
	}
	measures = append(measures, currentMeasure)

	// last measure has the same right barline as the original one
	lastMeasure := measures[len(measures)-1]
	lastMeasure.RightBarline = measure.RightBarline

	return measures
}

// distributeTimelines takes timelines that are noted in one part but affect another one
// e.g. second of part 2 which may be located in part 1 and copies this measure time line into
// the right part
func (p *partIterator) distributeTimelines() {

}

func New(tune *music_model.Tune) interfaces.MusicPartIterator {
	pi := &partIterator{
		currentIdx: -1,
	}
	pi.parts = splitTuneIntoParts(tune)
	pi.distributeTimelines()

	return pi
}
