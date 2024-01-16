package part_iterator

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/barline"
	"banduslib/internal/common/music_model/symbols/time_line"
	"banduslib/internal/interfaces"
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

func tuneToParts(tune *music_model.Tune) (parts []*music_model.MusicPart) {
	currPart := &music_model.MusicPart{}
	var partMeasures []*music_model.Measure
	for i, measure := range tune.Measures {
		partMeasures = append(partMeasures, measure)

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

func splitMeasureWithTimeline(measure *music_model.Measure) (measures []*music_model.Measure) {
	if !measure.HasTimeline() {
		measures = append(measures, measure)
		return measures
	}

	currentMeasure := &music_model.Measure{}
	for _, symbol := range measure.Symbols {
		if symbol.IsTimeline() &&
			symbol.TimeLine.BoundaryType == time_line.End {

		}
 	}

}

func New(tune *music_model.Tune) interfaces.MusicPartIterator {
	pi := &partIterator{
		currentIdx: -1,
	}
	pi.parts = tuneToParts(tune)
	return pi
}
