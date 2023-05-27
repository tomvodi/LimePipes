package bww

import "github.com/alecthomas/participle/v2/lexer"

var BwwLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{
			Name:    "BagpipeReader",
			Pattern: `Bagpipe Reader|Bagpipe Music Writer Gold|Bagpipe Musicworks Gold`,
			Action:  lexer.Push("BagpipeReader"),
		},
		{
			Name:    "TEMPO_DEF",
			Pattern: `TuneTempo`,
			Action:  lexer.Push("TuneTempo"),
		},
		{
			Name:    "PARAM_START",
			Pattern: `\(`,
			Action:  lexer.Push("ParamList"),
		},
		{
			Name:    "TIMELINE_END",
			Pattern: `_'|bis_'`,
		},
		{
			Name:    "DALSEGNO",
			Pattern: `dalsegno`,
		},
		{
			Name:    "DACAPOALFINE",
			Pattern: `dacapoalfine`,
		},
		{
			Name:    "PARAM_DEF",
			Pattern: `MIDINoteMappings|FrequencyMappings|InstrumentMappings|GracenoteDurations|FontSizes|TuneFormat`,
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "STAFF_START",
			Pattern: `&`,
			Action:  lexer.Push("Staff"),
		},
		{
			Name:    "STRING",
			Pattern: `"[^"]*"`,
		},
		{
			Name:    "WHITESPACE",
			Pattern: `\s+`,
		},
	},
	"BagpipeReader": {
		{
			Name:    "VERSION_SEP",
			Pattern: `:`,
		},
		{
			Name:    "VersionNumber",
			Pattern: `\d+\.\d+`,
			Action:  lexer.Pop(),
		},
	},
	"TuneTempo": {
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "TEMPO_VALUE",
			Pattern: `\d+`,
			Action:  lexer.Pop(),
		},
	},
	"ParamList": {
		{
			Name:    "PARAM_END",
			Pattern: `\)`,
			Action:  lexer.Pop(),
		},
		{
			Name:    "PARAM",
			Pattern: `[^,\)]+`,
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
	},
	"Staff": {
		{
			Name:    "STAFF_END",
			Pattern: `''!I|!t|!I`,
			Action:  lexer.Pop(),
		},
		{
			Name:    "PART_START",
			Pattern: `I!''|I!`,
		},
		{
			Name:    "NEXT_STAFF_START",
			Pattern: `&`,
		},
		{
			Name:    "BARLINE",
			Pattern: `!\s`,
		},
		{
			Name:    "SEGNO",
			Pattern: `segno`,
		},
		{
			Name:    "FINE",
			Pattern: `fine`,
		},
		{
			Name:    "SPACE",
			Pattern: `space`,
		},
		{
			Name:    "SHARP",
			Pattern: `sharplg|sharpla|sharpb|sharpc|sharpd|sharpe|sharpf|sharphg|sharpha`,
		},
		{
			Name:    "NATURAL",
			Pattern: `naturallg|naturalla|naturalb|naturalc|naturald|naturale|naturalf|naturalhg|naturalha`,
		},
		{
			Name:    "FLAT",
			Pattern: `flatlg|flatla|flatb|flatc|flatd|flate|flatf|flathg|flatha`,
		},
		{
			Name:    "HALF_NOTE",
			Pattern: `LG_2|LA_2|B_2|C_2|D_2|E_2|F_2|HG_2|HA_2`,
		},
		{
			Name:    "QUARTER_NOTE",
			Pattern: `LG_4|LA_4|B_4|C_4|D_4|E_4|F_4|HG_4|HA_4`,
		},
		{
			Name:    "EIGHTH_NOTE",
			Pattern: `LG_8|LA_8|B_8|C_8|D_8|E_8|F_8|HG_8|HA_8|LGl_8|LAl_8|Bl_8|Cl_8|Dl_8|El_8|Fl_8|HGl_8|HAl_8|LGr_8|LAr_8|Br_8|Cr_8|Dr_8|Er_8|Fr_8|HGr_8|HAr_8`,
		},
		{
			Name:    "SIXTEENTH_NOTE",
			Pattern: `LG_16|LA_16|B_16|C_16|D_16|E_16|F_16|HG_16|HA_16|LGl_16|LAl_16|Bl_16|Cl_16|Dl_16|El_16|Fl_16|HGl_16|HAl_16|LGr_16|LAr_16|Br_16|Cr_16|Dr_16|Er_16|Fr_16|HGr_16|HAr_16`,
		},
		{
			Name:    "THIRTYSECOND_NOTE",
			Pattern: `LG_32|LA_32|B_32|C_32|D_32|E_32|F_32|HG_32|HA_32|LGl_32|LAl_32|Bl_32|Cl_32|Dl_32|El_32|Fl_32|HGl_32|HAl_32|LGr_32|LAr_32|Br_32|Cr_32|Dr_32|Er_32|Fr_32|HGr_32|HAr_32`,
		},
		{
			Name:    "WHOLE_NOTE",
			Pattern: `LG_1|LA_1|B_1|C_1|D_1|E_1|F_1|HG_1|HA_1`,
		},
		{
			Name:    "TIME_SIG",
			Pattern: `2_2|3_2|2_4|3_4|4_4|5_4|5_4|6_4|7_4|C_|C|2_8|3_8|4_8|5_8|6_8|7_8|8_8|9_8|10_8|11_8|12_8|15_8|18_8|21_8|2_16|3_16|4_16|5_16|6_16|7_16|8_16|9_16|10_16|11_16|12_16`,
		},
		{
			Name:    "REST",
			Pattern: `REST_16|REST_1|REST_2|REST_4|REST_8|REST_32`,
		},
		{
			Name:    "FERMATA",
			Pattern: `fermatlg|fermatla|fermatb|fermatc|fermatd|fermate|fermatf|fermathg|fermatha`,
		},
		{
			Name:    "DOUBLING",
			Pattern: `dblg|dbla|dbb|dbc|dbd|dbe|dbf|dbhg|dbha`,
		},
		{
			Name:    "HALF_DOUBLING",
			Pattern: `hdblg|hdbla|hdbb|hdbc|hdbd|hdbe|hdbf`,
		},
		{
			Name:    "THUMB_DOUBLING",
			Pattern: `tdblg|tdbla|tdbb|tdbc|tdbd|tdbe|tdbf`,
		},
		{
			Name:    "STRIKE",
			Pattern: `strlg|strla|strb|strc|strd|stre|strf|strhg`,
		},
		{
			Name:    "LIGHT_G_STRIKE",
			Pattern: `lgstd`,
		},
		{
			Name:    "G_STRIKE",
			Pattern: `gstla|gstb|gstc|gstd|gste|gstf`,
		},
		{
			Name:    "LIGHT_THUMB_STRIKE",
			Pattern: `ltstd`,
		},
		{
			Name:    "THUMB_STRIKE",
			Pattern: `tstla|tstb|tstc|tstd|ltstd|tste|tstf|tsthg`,
		},
		{
			Name:    "LIGHT_HALF_STRIKE",
			Pattern: `lhstd`,
		},
		{
			Name:    "HALF_STRIKE",
			Pattern: `hstla|hstb|hstc|hstd|lhstd|hste|hstf|hsthg`,
		},
		{
			Name:    "G_GRIP",
			Pattern: `ggrpla|ggrpb|ggrpc|ggrpdb|ggrpd|ggrpe|ggrpf`,
		},
		{
			Name:    "THUMB_GRIP",
			Pattern: `tgrpla|tgrpb|tgrpc|tgrpdb|tgrpd|tgrpe|tgrpf|tgrphg`,
		},
		{
			Name:    "HALF_GRIP",
			Pattern: `hgrpla|hgrpb|hgrpc|hgrpdb|hgrpd|hgrpe|hgrpf|hgrphg|hgrpha`,
		},
		{
			Name:    "GRIP",
			Pattern: `grpb|grp|hgrp|grpb`,
		},
		{
			Name:    "BUBBLY",
			Pattern: `bubly|hbubly`,
		},
		{
			Name:    "G_BIRL",
			Pattern: `gbr`,
		},
		{
			Name:    "THUMB_BIRL",
			Pattern: `tbr`,
		},
		{
			Name:    "A_BIRL",
			Pattern: `abr`,
		},
		{
			Name:    "BIRL",
			Pattern: `brl`,
		},
		{
			Name:    "THROWD",
			Pattern: `hthrd|thrd`,
		},
		{
			Name:    "HEAVY_THROWD",
			Pattern: `hhvthrd|hvthrd`,
		},
		{
			Name:    "LIGHT_PELE",
			Pattern: `lpeld`,
		},
		{
			Name:    "PELE",
			Pattern: `pella|pelb|pelc|peld|pele|pelf`,
		},
		{
			Name:    "LIGHT_THUMB_PELE",
			Pattern: `ltpeld`,
		},
		{
			Name:    "THUMB_PELE",
			Pattern: `tpella|tpelb|tpelc|tpeld|tpele|tpelf|tpelhg`,
		},
		{
			Name:    "LIGHT_HALF_PELE",
			Pattern: `lhpeld`,
		},
		{
			Name:    "HALF_PELE",
			Pattern: `hpella|hpelb|hpelc|hpeld|hpele|hpelf|hpelhg`,
		},
		{
			Name:    "LIGHT_DOUBLE_STRIKE",
			Pattern: `lst2d`,
		},
		{
			Name:    "DOUBLE_STRIKE",
			Pattern: `st2la|st2b|st2c|st2d|st2e|st2f|st2hg|st2ha`,
		},
		{
			Name:    "LIGHT_G_DOUBLE_STRIKE",
			Pattern: `lgst2d`,
		},
		{
			Name:    "G_DOUBLE_STRIKE",
			Pattern: `gst2la|gst2b|gst2c|gst2d|gst2e|gst2f`,
		},
		{
			Name:    "LIGHT_THUMB_DOUBLE_STRIKE",
			Pattern: `ltst2d`,
		},
		{
			Name:    "THUMB_DOUBLE_STRIKE",
			Pattern: `tst2la|tst2b|tst2c|tst2d|tst2e|tst2f|tst2hg`,
		},
		{
			Name:    "LIGHT_HALF_DOUBLE_STRIKE",
			Pattern: `lhst2d`,
		},
		{
			Name:    "HALF_DOUBLE_STRIKE",
			Pattern: `hst2la|hst2b|hst2c|hst2d|hst2e|hst2f|hst2hg|hst2ha`,
		},
		{
			Name:    "LIGHT_TRIPLE_STRIKE",
			Pattern: `lst3d`,
		},
		{
			Name:    "TRIPLE_STRIKE",
			Pattern: `st3la|st3b|st3c|st3d|st3e|st3f|st3hg|st3ha`,
		},
		{
			Name:    "LIGHT_G_TRIPLE_STRIKE",
			Pattern: `lgst3d`,
		},
		{
			Name:    "G_TRIPLE_STRIKE",
			Pattern: `gst3la|gst3b|gst3c|gst3d|gst3e|gst3f`,
		},
		{
			Name:    "LIGHT_THUMB_TRIPLE_STRIKE",
			Pattern: `ltst3d`,
		},
		{
			Name:    "THUMB_TRIPLE_STRIKE",
			Pattern: `tst3la|tst3b|tst3c|tst3d|tst3e|tst3f|tst3hg`,
		},
		{
			Name:    "LIGHT_HALF_TRIPLE_STRIKE",
			Pattern: `lhst3d`,
		},
		{
			Name:    "HALF_TRIPLE_STRIKE",
			Pattern: `hst3la|hst3b|hst3c|hst3d|hst3e|hst3f|hst3hg|hst3ha`,
		},
		{
			Name:    "CADENCE",
			Pattern: `cadged|cadge|caded|cade|cadaed|cadae|cadgf|cadaf`,
		},
		{
			Name:    "FERMAT_CADENCE",
			Pattern: `fcadged|fcadge|fcaded|fcade|fcadaed|fcadae|fcadgf|fcadaf`,
		},
		{
			Name:    "EMBARI",
			Pattern: `embari|pembari`,
		},
		{
			Name:    "ENDARI",
			Pattern: `endari|pendari`,
		},
		{
			Name:    "CHEDARI",
			Pattern: `chedari|pchedari`,
		},
		{
			Name:    "HEDARI",
			Pattern: `hedari|phedari`,
		},
		{
			Name:    "DILI",
			Pattern: `dili|pdili`,
		},
		{
			Name:    "TRA",
			Pattern: `tra8|tra|htra|ptra8|ptra|phtra`,
		},
		{
			Name:    "EDRE",
			Pattern: `edreb|edrec|edred|edre|pedreb|pedrec|pedred|pedre`,
		},
		{
			Name:    "G_EDRE",
			Pattern: `gedre`,
		},
		{
			Name:    "THUMB_EDRE",
			Pattern: `tedre`,
		},
		{
			Name:    "HALF_EDRE",
			Pattern: `dre`,
		},
		{
			Name:    "DARE",
			Pattern: `chedare|dare|pdare`,
		},
		{
			Name:    "G_DARE",
			Pattern: `gdare`,
		},
		{
			Name:    "THUMB_DARE",
			Pattern: `tdare`,
		},
		{
			Name:    "HALF_DARE",
			Pattern: `hedale`,
		},
		{
			Name:    "CHECHERE",
			Pattern: `chechere|pchechere`,
		},
		{
			Name:    "THUMB_CHECHERE",
			Pattern: `tchechere`,
		},
		{
			Name:    "HALF_CHECHERE",
			Pattern: `hchechere`,
		},
		{
			Name:    "GRIP_ABBREV",
			Pattern: `pgrp`,
		},
		{
			Name:    "DEDA",
			Pattern: `deda`,
		},
		{
			Name:    "ENBAIN",
			Pattern: `enbain|penbain`,
		},
		{
			Name:    "G_ENBAIN",
			Pattern: `genbain`,
		},
		{
			Name:    "THUMB_ENBAIN",
			Pattern: `tenbain`,
		},
		{
			Name:    "OTRO",
			Pattern: `otro|potro`,
		},
		{
			Name:    "G_OTRO",
			Pattern: `gotro`,
		},
		{
			Name:    "THUMB_OTRO",
			Pattern: `totro`,
		},
		{
			Name:    "ODRO",
			Pattern: `odro|podro`,
		},
		{
			Name:    "G_ODRO",
			Pattern: `godro`,
		},
		{
			Name:    "THUMB_ODRO",
			Pattern: `todro`,
		},
		{
			Name:    "ADEDA",
			Pattern: `adeda|padeda`,
		},
		{
			Name:    "G_ADEDA",
			Pattern: `gadeda`,
		},
		{
			Name:    "THUMB_ADEDA",
			Pattern: `tadeda`,
		},
		{
			Name:    "ECHO_BEATS",
			Pattern: `echolg|echola|echohg|echoha|echob|echoc|echod|echoe|echof`,
		},
		{
			Name:    "DARODO",
			Pattern: `phdarodo|pdarodo16|pdarodo|hdarodo|darodo16|darodo`,
		},
		{
			Name:    "HIHARIN",
			Pattern: `hiharin|phiharin`,
		},
		{
			Name:    "RODIN",
			Pattern: `rodin`,
		},
		{
			Name:    "CHELALHO",
			Pattern: `chelalho`,
		},
		{
			Name:    "DIN",
			Pattern: `din`,
		},
		{
			Name:    "LEMLUATH",
			Pattern: `hlemlabrea|lembbrea|lembrea|hlemlg|hlemla|lemb|lem`,
		},
		{
			Name:    "LEMLUATH_ABBREV",
			Pattern: `plbrea|plbbrea|phllabrea|phlla|plb|pl`,
		},
		{
			Name:    "TAORLUATH_PIO",
			Pattern: `htarlabrea|tarbbrea|tarbrea|htarlg|htarla`,
		},
		{
			Name:    "TAORLUATH",
			Pattern: `tarb|tar|htar`,
		},
		{
			Name:    "TRIPLINGS",
			Pattern: `phtriplg|phtripla|phtripb|phtripc|pttriplg|pttripla|pttripb|pttripc|ptriplg|ptripla|ptripb|ptripc`,
		},
		{
			Name:    "TAORLUATH_AMACH",
			Pattern: `ptmb|ptmc|ptmd`,
		},
		{
			Name:    "TAORLUATH_ABBREV",
			Pattern: `ptbrea|ptbbrea|phtlabrea|phtla|ptb|pt`,
		},
		{
			Name:    "CRUNLUATH",
			Pattern: `hcrunllabrea|crunlbbrea|crunlbrea|hcrunllgla|hcrunllg|hcrunlla|crunlb|crunl`,
		},
		{
			Name:    "CRUNLUATH_AMACH",
			Pattern: `pcmb|pcmc|pcmd`,
		},
		{
			Name:    "CRUNLUATH_ABBREV",
			Pattern: `pcbrea|pcbbrea|phclabrea|phcla|pcb|pc`,
		},
		{
			Name:    "D_DOUBLE_GRACE",
			Pattern: `dlg|dla|db|dc`,
		},
		{
			Name:    "E_DOUBLE_GRACE",
			Pattern: `elg|ela|eb|ec|ed`,
		},
		{
			Name:    "F_DOUBLE_GRACE",
			Pattern: `flg|fla|fb|fc|fd|fe`,
		},
		{
			Name:    "G_DOUBLE_GRACE",
			Pattern: `glg|gla|gb|gc|gd|ge|gf`,
		},
		{
			Name:    "THUMB_DOUBLE_GRACE",
			Pattern: `tlg|tla|tb|tc|td|te|tf|thg`,
		},
		{
			Name:    "SINGLE_GRACE",
			Pattern: `ag|bg|cg|dg|eg|fg|gg|tg`,
		},
		{
			Name:    "TIE_START",
			Pattern: `\^ts`,
		},
		{
			Name:    "TIE_END",
			Pattern: `\^te`,
		},
		{
			Name:    "TIE_OLD",
			Pattern: `\^tlg|\^tla|\^tb|\^tc|\^td|\^tf|\^thg|\^tha`,
		},
		{
			Name:    "IRREGULAR_GROUP_START",
			Pattern: `\^2s|\^3s|\^43s|\^46s|\^53s|\^54s|\^64s|\^74s|\^76s`,
		},
		{
			Name:    "IRREGULAR_GROUP_END",
			Pattern: `\^2e|\^3e|\^43e|\^46e|\^53e|\^54e|\^64e|\^74e|\^76e`,
		},
		{
			Name:    "TRIPLETS",
			Pattern: `\^3lg|\^3la|\^3b|\^3c|\^3d|\^3f|\^3hg|\^3ha`,
		},
		{
			Name:    "TIMELINE_START",
			Pattern: `'224|'22|'23|'24|'intro|'25|'26|'27|'28|'1|'2|'si|'do|'bis`,
		},
		{
			Name:    "TIMELINE_END",
			Pattern: `_'|bis_'`,
		},
		{
			Name:    "SINGLE_DOT",
			Pattern: `'lg|'la|'b|'c|'d|'e|'f|'hg|'ha`,
		},
		{
			Name:    "DOUBLE_DOT",
			Pattern: `''lg|''la|''b|''c|''d|''e|''f|''hg|''ha`,
		},
		{
			Name:    "STRING",
			Pattern: `"[^"]*"`,
		},
		{
			Name:    "TEMPO_DEF",
			Pattern: `TuneTempo`,
			Action:  lexer.Push("TuneTempo"),
		},
		{
			Name:    "PARAM_START",
			Pattern: `\(`,
			Action:  lexer.Push("ParamList"),
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "WHITESPACE",
			Pattern: `\s+`,
		},
	},
})
