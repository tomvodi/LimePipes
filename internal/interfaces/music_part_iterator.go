package interfaces

import "banduslib/internal/common/music_model"

// MusicPartIterator returns a tune split up into parts.
// If the part is a repeated one it may contain alternative endings (time lines).
// If the alternative ending is in the middle of a measure,
// the iterator returns a measure for the first half and for the second half
// to be compatible with MusicXML as in bww files, the start of a time line
// can be in the middle of a measure. This is often the case, when the first and
// the second time line have different up-beats.
// Also, if there is an advanced time line in the underlying tune (e.g. 2-of-4),
// this will also be normalized, so that the 2-of-4 time line is copied to part 4.
type MusicPartIterator interface {
	// Count returns the amount of parts
	Count() int

	// Nr returns the part with number. First part has number 1 and not 0
	// Returns nil if there is no part with that number
	Nr(nr int) *music_model.MusicPart

	HasNext() bool
	GetNext() *music_model.MusicPart
}
