package bww

import "fmt"

type BwwDocument struct {
	Tunes []*Tune `@@*`
}

type Tune struct {
	BagpipePlayerVersion string      `(BagpipeReader VERSION_SEP @VersionNumber)?`
	Header               *TuneHeader `@@+`
	Body                 *TuneBody   `@@`
}

type TuneHeader struct {
	TuneParameter []*TuneParameter `@@+`
}

type TuneParameter struct {
	Config      *TuneConfig      `@@`
	Tempo       *TuneTempo       `| @@`
	Description *TuneDescription `| @@`
	Comment     string           `| STRING`
}

// TuneConfig like page layout or MIDI note mappings
// these lines start with a defined word e.g. MIDINoteMappings,(...)
type TuneConfig struct {
	Name      string   `@PARAM_DEF PARAM_SEP`
	ParamList []string `PARAM_START @PARAM (PARAM_SEP @PARAM)* PARAM_END`
}

type TuneTempo struct {
	Tempo string `TEMPO_DEF PARAM_SEP @TEMPO_VALUE`
}

// TuneDescription like title, composer, arranger
// they all start with a string "title",(...)
type TuneDescription struct {
	Text      string   `@STRING PARAM_SEP`
	ParamList []string `PARAM_START @PARAM (PARAM_SEP @PARAM)* PARAM_END`
}

type TuneBody struct {
	Staffs []*Staff `@@*`
}

type Staff struct {
	Start   string          `@STAFF_START`
	Symbols []*StaffSymbols `@@*`
	End     string          `@STAFF_END`
}

type StaffSymbols struct {
	PartStart           *string `@PART_START`
	Barline             *string `| @BARLINE`
	TimeSig             *string `| @TIME_SIG`
	Sharp               *string `| @SHARP`
	Natural             *string `| @NATURAL`
	Flat                *string `| @FLAT`
	WholeNote           *string `| @WHOLE_NOTE`
	HalfNote            *string `| @HALF_NOTE`
	QuarterNote         *string `| @QUARTER_NOTE`
	EighthNote          *string `| @EIGHTH_NOTE`
	SixteenthNote       *string `| @SIXTEENTH_NOTE`
	ThirtysecondNote    *string `| @THIRTYSECOND_NOTE`
	Rest                *string `| @REST`
	SingleDots          *string `| @SINGLE_DOT`
	DoubleDots          *string `| @DOUBLE_DOT`
	Fermata             *string `| @FERMATA`
	SingleGrace         *string `| @SINGLE_GRACE`
	Doubling            *string `| @DOUBLING`
	HalfDoubling        *string `| @HALF_DOUBLING`
	ThumbDoubling       *string `| @THUMB_DOUBLING`
	Strike              *string `| @STRIKE`
	GStrike             *string `| @G_STRIKE`
	ThumbStrike         *string `| @THUMB_STRIKE`
	HalfStrike          *string `| @HALF_STRIKE`
	Grip                *string `| @GRIP`
	GGrip               *string `| @G_GRIP`
	ThumbGrip           *string `| @THUMB_GRIP`
	HalfGrip            *string `| @HALF_GRIP`
	Taorluath           *string `| @TAORLUATH`
	Bubbly              *string `| @BUBBLY`
	Birl                *string `| @BIRL`
	ThrowD              *string `| @THROWD`
	HeavyThrowD         *string `| @HEAVY_THROWD`
	HalfThrowD          *string `| @HALF_THROWD`
	HeavyHalfThrowD     *string `| @HEAVY_HALF_THROWD`
	Pele                *string `| @PELE`
	ThumbPele           *string `| @THUMB_PELE`
	HalfPele            *string `| @HALF_PELE`
	DoubleStrike        *string `| @DOUBLE_STRIKE`
	GDoubleStrike       *string `| @G_DOUBLE_STRIKE`
	ThumbDoubleStrike   *string `| @THUMB_DOUBLE_STRIKE`
	HalfDoubleStrike    *string `| @HALF_DOUBLE_STRIKE`
	TripleStrike        *string `| @TRIPLE_STRIKE`
	GTripleStrike       *string `| @G_TRIPLE_STRIKE`
	ThumbTripleStrike   *string `| @THUMB_TRIPLE_STRIKE`
	HalfTripleStrike    *string `| @HALF_TRIPLE_STRIKE`
	DDoubleGrace        *string `| @D_DOUBLE_GRACE`
	EDoubleGrace        *string `| @E_DOUBLE_GRACE`
	FDoubleGrace        *string `| @F_DOUBLE_GRACE`
	GDoubleGrace        *string `| @G_DOUBLE_GRACE`
	ThumbDoubleGrace    *string `| @THUMB_DOUBLE_GRACE`
	TieStart            *string `| @TIE_START`
	TieEnd              *string `| @TIE_END`
	TieOld              *string `| @TIE_OLD`
	IrregularGroupStart *string `| @IRREGULAR_GROUP_START`
	IrregularGroupEnd   *string `| @IRREGULAR_GROUP_END`
	Triplets            *string `| @TRIPLETS`
	TimelineStart       *string `| @TIMELINE_START`
	TimelineEnd         *string `| @TIMELINE_END`
}

func (s StaffSymbols) String() string {
	if s.PartStart != nil {
		return fmt.Sprintf("PartStart(%s)", *s.PartStart)
	}
	if s.Barline != nil {
		return fmt.Sprintf("Barline(%s)", *s.Barline)
	}
	if s.TimeSig != nil {
		return fmt.Sprintf("TimeSig(%s)", *s.TimeSig)
	}
	if s.Sharp != nil {
		return fmt.Sprintf("Sharp(%s)", *s.Sharp)
	}
	if s.Natural != nil {
		return fmt.Sprintf("Natural(%s)", *s.Natural)
	}
	if s.Flat != nil {
		return fmt.Sprintf("Flat(%s)", *s.Flat)
	}
	if s.WholeNote != nil {
		return fmt.Sprintf("WholeNote(%s)", *s.WholeNote)
	}
	if s.HalfNote != nil {
		return fmt.Sprintf("HalfNote(%s)", *s.HalfNote)
	}
	if s.QuarterNote != nil {
		return fmt.Sprintf("QuarterNote(%s)", *s.QuarterNote)
	}
	if s.EighthNote != nil {
		return fmt.Sprintf("EighthNote(%s)", *s.EighthNote)
	}
	if s.SixteenthNote != nil {
		return fmt.Sprintf("SixteenthNote(%s)", *s.SixteenthNote)
	}
	if s.ThirtysecondNote != nil {
		return fmt.Sprintf("ThirtysecondNote(%s)", *s.ThirtysecondNote)
	}
	if s.Rest != nil {
		return fmt.Sprintf("Rest(%s)", *s.Rest)
	}
	if s.SingleDots != nil {
		return fmt.Sprintf("SingleDots(%s)", *s.SingleDots)
	}
	if s.DoubleDots != nil {
		return fmt.Sprintf("DoubleDots(%s)", *s.DoubleDots)
	}
	if s.Fermata != nil {
		return fmt.Sprintf("Fermata(%s)", *s.Fermata)
	}
	if s.SingleGrace != nil {
		return fmt.Sprintf("SingleGrace(%s)", *s.SingleGrace)
	}
	if s.Doubling != nil {
		return fmt.Sprintf("Doubling(%s)", *s.Doubling)
	}
	if s.HalfDoubling != nil {
		return fmt.Sprintf("HalfDoubling(%s)", *s.HalfDoubling)
	}
	if s.ThumbDoubling != nil {
		return fmt.Sprintf("ThumbDoubling(%s)", *s.ThumbDoubling)
	}
	if s.Strike != nil {
		return fmt.Sprintf("Strike(%s)", *s.Strike)
	}
	if s.GStrike != nil {
		return fmt.Sprintf("GStrike(%s)", *s.GStrike)
	}
	if s.ThumbStrike != nil {
		return fmt.Sprintf("ThumbStrike(%s)", *s.ThumbStrike)
	}
	if s.HalfStrike != nil {
		return fmt.Sprintf("HalfStrike(%s)", *s.HalfStrike)
	}
	if s.Grip != nil {
		return fmt.Sprintf("Grip(%s)", *s.Grip)
	}
	if s.GGrip != nil {
		return fmt.Sprintf("GGrip(%s)", *s.GGrip)
	}
	if s.ThumbGrip != nil {
		return fmt.Sprintf("ThumbGrip(%s)", *s.ThumbGrip)
	}
	if s.HalfGrip != nil {
		return fmt.Sprintf("HalfGrip(%s)", *s.HalfGrip)
	}
	if s.Taorluath != nil {
		return fmt.Sprintf("Taorluath(%s)", *s.Taorluath)
	}
	if s.Bubbly != nil {
		return fmt.Sprintf("Bubbly(%s)", *s.Bubbly)
	}
	if s.Birl != nil {
		return fmt.Sprintf("Birl(%s)", *s.Birl)
	}
	if s.ThrowD != nil {
		return fmt.Sprintf("ThrowD(%s)", *s.ThrowD)
	}
	if s.HeavyThrowD != nil {
		return fmt.Sprintf("HeavyThrowD(%s)", *s.HeavyThrowD)
	}
	if s.HalfThrowD != nil {
		return fmt.Sprintf("HalfThrowD(%s)", *s.HalfThrowD)
	}
	if s.HeavyHalfThrowD != nil {
		return fmt.Sprintf("HeavyHalfThrowD(%s)", *s.HeavyHalfThrowD)
	}
	if s.Pele != nil {
		return fmt.Sprintf("Pele(%s)", *s.Pele)
	}
	if s.ThumbPele != nil {
		return fmt.Sprintf("ThumbPele(%s)", *s.ThumbPele)
	}
	if s.HalfPele != nil {
		return fmt.Sprintf("HalfPele(%s)", *s.HalfPele)
	}
	if s.DoubleStrike != nil {
		return fmt.Sprintf("DoubleStrike(%s)", *s.DoubleStrike)
	}
	if s.GDoubleStrike != nil {
		return fmt.Sprintf("GDoubleStrike(%s)", *s.GDoubleStrike)
	}
	if s.ThumbDoubleStrike != nil {
		return fmt.Sprintf("ThumbDoubleStrike(%s)", *s.ThumbDoubleStrike)
	}
	if s.HalfDoubleStrike != nil {
		return fmt.Sprintf("HalfDoubleStrike(%s)", *s.HalfDoubleStrike)
	}
	if s.TripleStrike != nil {
		return fmt.Sprintf("TripleStrike(%s)", *s.TripleStrike)
	}
	if s.GTripleStrike != nil {
		return fmt.Sprintf("GTripleStrike(%s)", *s.GTripleStrike)
	}
	if s.ThumbTripleStrike != nil {
		return fmt.Sprintf("ThumbTripleStrike(%s)", *s.ThumbTripleStrike)
	}
	if s.HalfTripleStrike != nil {
		return fmt.Sprintf("HalfTripleStrike(%s)", *s.HalfTripleStrike)
	}
	if s.DDoubleGrace != nil {
		return fmt.Sprintf("DDoubleGrace(%s)", *s.DDoubleGrace)
	}
	if s.EDoubleGrace != nil {
		return fmt.Sprintf("EDoubleGrace(%s)", *s.EDoubleGrace)
	}
	if s.FDoubleGrace != nil {
		return fmt.Sprintf("FDoubleGrace(%s)", *s.FDoubleGrace)
	}
	if s.GDoubleGrace != nil {
		return fmt.Sprintf("GDoubleGrace(%s)", *s.GDoubleGrace)
	}
	if s.ThumbDoubleGrace != nil {
		return fmt.Sprintf("ThumbDoubleGrace(%s)", *s.ThumbDoubleGrace)
	}
	if s.TieStart != nil {
		return fmt.Sprintf("TieStart(%s)", *s.TieStart)
	}
	if s.TieEnd != nil {
		return fmt.Sprintf("TieEnd(%s)", *s.TieEnd)
	}
	if s.TieOld != nil {
		return fmt.Sprintf("TieOld(%s)", *s.TieOld)
	}
	if s.IrregularGroupStart != nil {
		return fmt.Sprintf("IrregularGroupStart(%s)", *s.IrregularGroupStart)
	}
	if s.IrregularGroupEnd != nil {
		return fmt.Sprintf("IrregularGroupEnd(%s)", *s.IrregularGroupEnd)
	}
	if s.Triplets != nil {
		return fmt.Sprintf("Triplets(%s)", *s.Triplets)
	}
	if s.TimelineStart != nil {
		return fmt.Sprintf("TimelineStart(%s)", *s.TimelineStart)
	}
	if s.TimelineEnd != nil {
		return fmt.Sprintf("TimelineEnd(%s)", *s.TimelineEnd)
	}
	return ""
}
