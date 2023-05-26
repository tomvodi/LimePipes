package import_message

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Type
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Fix

type Type uint

const (
	NoType Type = iota
	Info
	Warning
	Error
)

type Fix uint

const (
	NoFix Fix = iota
	SkipSymbol
)

type ImportMessage struct {
	Symbol string
	Type   Type
	Text   string
	Fix    Fix
}
