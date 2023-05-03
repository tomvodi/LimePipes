package common

type BwwDocument struct {
	Tunes []*Tune `@@*`
}

type Tune struct {
	BagpipePlayerVersion string      `BagpipeReader VERSION_SEP @VersionNumber`
	Header               *TuneHeader `@@`
	Body                 *TuneBody   `@@`
}

type TuneHeader struct {
	TuneParameter []*TuneParameter `@@*`
}

type TuneParameter struct {
	Config      *TuneConfig      `@@`
	Tempo       *TuneTempo       `| @@`
	Description *TuneDescription `| @@`
}

type TuneConfig struct {
	Name      string   `@PARAM_DEF PARAM_SEP`
	ParamList []string `PARAM_START @PARAM (PARAM_SEP @PARAM)* PARAM_END`
}

type TuneTempo struct {
	Tempo string `TEMPO_DEF PARAM_SEP @TEMPO_VALUE`
}

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
	TimeSig             *string `| @TIME_SIG`
	Sharp               *string `| @SHARP`
	Natural             *string `| @NATURAL`
	FLAT                *string `| @FLAT`
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
	Bubly               *string `| @BUBLY`
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
