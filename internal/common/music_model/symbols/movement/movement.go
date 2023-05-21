package movement

import "banduslib/internal/common"

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=MovementType

type MovementType uint

const (
	NoMovement MovementType = iota
	Cadence
)

type Movement struct {
	Type    MovementType   `yaml:"type"`
	Pitches []common.Pitch `yaml:"pitches,omitempty"` // Only used by cadence
	Fermata bool           `yaml:"fermata,omitempty"` // Only used by cadence
}
