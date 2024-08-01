package helper

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/movement"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/timeline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/tuplet"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
)

var CompareOpts = cmpopts.IgnoreUnexported(
	tune.Tune{},
	measure.Measure{},
	measure.TimeSignature{},
	measure.ImportMessage{},
	symbols.Symbol{},
	symbols.Note{},
	tune.Tune{},
	symbols.Rest{},
	embellishment.Embellishment{},
	tuplet.Tuplet{},
	movement.Movement{},
	timeline.TimeLine{},
	barline.Barline{},
)
