// Code generated by "enumer -json -yaml -type=EmbellishmentType"; DO NOT EDIT.

package embellishment

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _EmbellishmentTypeName = "NoEmbellishmentSingleGraceDoublingStrikeGripTaorluathBubblyBirlThrowDPeleDoubleStrikeTripleStrikeGTripleStrikeThumbTripleStrikeHalfTripleStrikeDoubleGrace"

var _EmbellishmentTypeIndex = [...]uint8{0, 15, 26, 34, 40, 44, 53, 59, 63, 69, 73, 85, 97, 110, 127, 143, 154}

const _EmbellishmentTypeLowerName = "noembellishmentsinglegracedoublingstrikegriptaorluathbubblybirlthrowdpeledoublestriketriplestrikegtriplestrikethumbtriplestrikehalftriplestrikedoublegrace"

func (i EmbellishmentType) String() string {
	if i >= EmbellishmentType(len(_EmbellishmentTypeIndex)-1) {
		return fmt.Sprintf("EmbellishmentType(%d)", i)
	}
	return _EmbellishmentTypeName[_EmbellishmentTypeIndex[i]:_EmbellishmentTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _EmbellishmentTypeNoOp() {
	var x [1]struct{}
	_ = x[NoEmbellishment-(0)]
	_ = x[SingleGrace-(1)]
	_ = x[Doubling-(2)]
	_ = x[Strike-(3)]
	_ = x[Grip-(4)]
	_ = x[Taorluath-(5)]
	_ = x[Bubbly-(6)]
	_ = x[Birl-(7)]
	_ = x[ThrowD-(8)]
	_ = x[Pele-(9)]
	_ = x[DoubleStrike-(10)]
	_ = x[TripleStrike-(11)]
	_ = x[GTripleStrike-(12)]
	_ = x[ThumbTripleStrike-(13)]
	_ = x[HalfTripleStrike-(14)]
	_ = x[DoubleGrace-(15)]
}

var _EmbellishmentTypeValues = []EmbellishmentType{NoEmbellishment, SingleGrace, Doubling, Strike, Grip, Taorluath, Bubbly, Birl, ThrowD, Pele, DoubleStrike, TripleStrike, GTripleStrike, ThumbTripleStrike, HalfTripleStrike, DoubleGrace}

var _EmbellishmentTypeNameToValueMap = map[string]EmbellishmentType{
	_EmbellishmentTypeName[0:15]:         NoEmbellishment,
	_EmbellishmentTypeLowerName[0:15]:    NoEmbellishment,
	_EmbellishmentTypeName[15:26]:        SingleGrace,
	_EmbellishmentTypeLowerName[15:26]:   SingleGrace,
	_EmbellishmentTypeName[26:34]:        Doubling,
	_EmbellishmentTypeLowerName[26:34]:   Doubling,
	_EmbellishmentTypeName[34:40]:        Strike,
	_EmbellishmentTypeLowerName[34:40]:   Strike,
	_EmbellishmentTypeName[40:44]:        Grip,
	_EmbellishmentTypeLowerName[40:44]:   Grip,
	_EmbellishmentTypeName[44:53]:        Taorluath,
	_EmbellishmentTypeLowerName[44:53]:   Taorluath,
	_EmbellishmentTypeName[53:59]:        Bubbly,
	_EmbellishmentTypeLowerName[53:59]:   Bubbly,
	_EmbellishmentTypeName[59:63]:        Birl,
	_EmbellishmentTypeLowerName[59:63]:   Birl,
	_EmbellishmentTypeName[63:69]:        ThrowD,
	_EmbellishmentTypeLowerName[63:69]:   ThrowD,
	_EmbellishmentTypeName[69:73]:        Pele,
	_EmbellishmentTypeLowerName[69:73]:   Pele,
	_EmbellishmentTypeName[73:85]:        DoubleStrike,
	_EmbellishmentTypeLowerName[73:85]:   DoubleStrike,
	_EmbellishmentTypeName[85:97]:        TripleStrike,
	_EmbellishmentTypeLowerName[85:97]:   TripleStrike,
	_EmbellishmentTypeName[97:110]:       GTripleStrike,
	_EmbellishmentTypeLowerName[97:110]:  GTripleStrike,
	_EmbellishmentTypeName[110:127]:      ThumbTripleStrike,
	_EmbellishmentTypeLowerName[110:127]: ThumbTripleStrike,
	_EmbellishmentTypeName[127:143]:      HalfTripleStrike,
	_EmbellishmentTypeLowerName[127:143]: HalfTripleStrike,
	_EmbellishmentTypeName[143:154]:      DoubleGrace,
	_EmbellishmentTypeLowerName[143:154]: DoubleGrace,
}

var _EmbellishmentTypeNames = []string{
	_EmbellishmentTypeName[0:15],
	_EmbellishmentTypeName[15:26],
	_EmbellishmentTypeName[26:34],
	_EmbellishmentTypeName[34:40],
	_EmbellishmentTypeName[40:44],
	_EmbellishmentTypeName[44:53],
	_EmbellishmentTypeName[53:59],
	_EmbellishmentTypeName[59:63],
	_EmbellishmentTypeName[63:69],
	_EmbellishmentTypeName[69:73],
	_EmbellishmentTypeName[73:85],
	_EmbellishmentTypeName[85:97],
	_EmbellishmentTypeName[97:110],
	_EmbellishmentTypeName[110:127],
	_EmbellishmentTypeName[127:143],
	_EmbellishmentTypeName[143:154],
}

// EmbellishmentTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func EmbellishmentTypeString(s string) (EmbellishmentType, error) {
	if val, ok := _EmbellishmentTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _EmbellishmentTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to EmbellishmentType values", s)
}

// EmbellishmentTypeValues returns all values of the enum
func EmbellishmentTypeValues() []EmbellishmentType {
	return _EmbellishmentTypeValues
}

// EmbellishmentTypeStrings returns a slice of all String values of the enum
func EmbellishmentTypeStrings() []string {
	strs := make([]string, len(_EmbellishmentTypeNames))
	copy(strs, _EmbellishmentTypeNames)
	return strs
}

// IsAEmbellishmentType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i EmbellishmentType) IsAEmbellishmentType() bool {
	for _, v := range _EmbellishmentTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for EmbellishmentType
func (i EmbellishmentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for EmbellishmentType
func (i *EmbellishmentType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("EmbellishmentType should be a string, got %s", data)
	}

	var err error
	*i, err = EmbellishmentTypeString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for EmbellishmentType
func (i EmbellishmentType) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for EmbellishmentType
func (i *EmbellishmentType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = EmbellishmentTypeString(s)
	return err
}
