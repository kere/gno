package util

const (
	// SEmptyString ''
	SEmptyString = ""
	// SQuot ''
	SQuot = "'"
	// SDoubleQuot ''
	SDoubleQuot = "\""

	//SComma ,
	SComma      = ","
	SBraceLeft  = "{"
	SBraceRight = "{"
	// SLineBreak = "\n"
	SLineBreak = "\n"
	// STab \t
	STab = "\t"
	// SDoller $
	SDoller = "$"
	// SDot.
	SDot = "."
)

var (
	// BDoller $
	BDoller = []byte("$")
	// BNull null
	BNull = []byte("NULL")
	// BEqual =
	BEqual = []byte("=")

	//BComma ,
	BComma = []byte(",")

	// BEmptyString ''
	BEmptyString = []byte("")
	//BStarKey *
	BStarKey = []byte("*")

	//BSpace ' '
	BSpace        = []byte(" ")
	BBraceLeft    = []byte("{")
	BBraceRight   = []byte("}")
	BBracketLeft  = []byte("[")
	BBracketRight = []byte("]")

	BQuestionMark = []byte("?")
	BQuote        = []byte("'")
	BDoubleQuote  = []byte("\"")
	BDot          = []byte(".")
	BCharE        = []byte("e")
	BZero         = []byte("0")
	BNaN          = []byte("NaN")
	BTab          = []byte("\t")
	BLineBreak    = []byte("\n")
)
