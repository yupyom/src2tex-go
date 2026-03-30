package main

// Constants mapped from src2tex.h options
const (
	Version        = "2.23"
	BufferSize     = 16384 // Equivalent to Unix BUFFER_SIZE
	MaxFileNameLen = 256
	Space          = 0.525
	HTabSkip       = 8
	VTabSkip       = 4
	FormulaLenMax  = 512

	LatexOption = ""
	LatexStyle  = "article"
)

// Font options corresponding to NTT/TeX default settings
var (
	Bold       = "\\bf "
	Italic     = "\\it "
	Roman      = "\\rm\\mc "
	SmallCaps  = "\\sc "
	Slant      = "\\sl "
	Typewriter = "\\tt\\mc "

	TextModeFont = Typewriter
	TeXModeFont  = Roman
)

// Global Configuration
var (
	HtabSize = HTabSkip
	VtabSize = VTabSkip
)

// FlagChar represents a flagged character in the input stream.
//
//	flag: 0 = text mode, 1 = quasi-TeX mode, 2 = TeX mode
//	character: input character (or EOF as -1)
//	bufferIdx: index in the circular buffer
type FlagChar struct {
	Flag      int
	Character int
	BufferIdx int
}

// Global mode flags (from langflag.c)
var (
	TXTFlag  int
	BASFlag  int
	CBLFlag  int
	CFlag    int
	F77Flag  int
	HTMLFlag int
	JAVAFlag int
	LISPFlag int
	MAKEFlag int
	PASFlag  int
	PERLFlag int
	SHFlag   int
	TCLFlag  int

	PythonFlag int
	RubyFlag   int
	RustFlag   int
	GoFlag     int
	JSFlag     int
	TSFlag     int
	KotlinFlag int
	SwiftFlag  int

	ASRFlag  int
	MACFlag  int
	MAPFlag  int
	MATFlag  int
	MLABFlag int
	MPADFlag int
	REDFlag  int

	MFlag       int  // MATLAB alias
	IsLatexMode bool // Set to true when running as src2latex

	IsUnicodeTeX  = true // default: generate Unicode-aware preamble (xeCJK/fontspec)
	InputEncoding string

	PageLenMax = -1 // Options related to output pagination

	HeaderFile string // Custom header preamble file (replaces fancyheader section)
	FooterFile string // Custom footer preamble file (replaces fancyfooter section)
)

// Circular buffer corresponding to getdata.c
var Buffer [BufferSize]int
var BufCounter = 0
var BufPtr = 0
var LineCounter = 0
