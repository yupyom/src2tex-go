package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var inputFileName string
var outputFileName string

var inputFile *os.File
var outputFile *os.File

// GetFileNames parses os.Args manually to mimic the old behaviour supporting -<n> options
func GetFileNames() {
	args := os.Args[1:]

	if len(args) == 0 {
		inputFileName = ""
		outputFileName = ""
		return
	}

	// Filter out standard flags
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-latex" || arg == "-tex" {
			continue
		}
		if arg == "-sjis" {
			InputEncoding = "sjis"
			continue
		}
		if arg == "-euc" {
			InputEncoding = "euc"
			continue
		}
		if arg == "-utf8" {
			InputEncoding = "utf8"
			continue
		}
		if arg == "-unicode" {
			IsUnicodeTeX = true
			continue
		}
		if arg == "-legacy" {
			IsUnicodeTeX = false
			continue
		}
		if arg == "-fontdir" {
			if i+1 < len(args) {
				i++
				FontDir = args[i]
			} else {
				fmt.Fprintf(os.Stderr, "Error: -fontdir requires a directory path\n")
				os.Exit(1)
			}
			continue
		}
		if arg == "-proxy" {
			if i+1 < len(args) {
				i++
				ProxyURL = args[i]
			} else {
				fmt.Fprintf(os.Stderr, "Error: -proxy requires a URL\n")
				os.Exit(1)
			}
			continue
		}
		if arg == "-font" {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: -font requires an argument\n")
				os.Exit(1)
			}
			i++
			fontArg := args[i]

			if fontArg == "list" {
				ListFonts()
				os.Exit(0)
			}
			if fontArg == "install" {
				if i+1 >= len(args) {
					fmt.Fprintf(os.Stderr, "Error: -font install requires a font name (or 'all')\n")
					os.Exit(1)
				}
				i++
				SetupProxy()
				if err := InstallFont(args[i]); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
			// Font name specified
			CJKFontName = fontArg
			FontExplicitlySet = true
			continue
		}
		if arg == "-commentfont" {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: -commentfont requires an argument\n")
				os.Exit(1)
			}
			i++
			cfArg := args[i]

			if cfArg == "list" {
				ListCommentFonts()
				os.Exit(0)
			}
			if cfArg == "install" {
				if i+1 >= len(args) {
					fmt.Fprintf(os.Stderr, "Error: -commentfont install requires a font name (or 'all')\n")
					os.Exit(1)
				}
				i++
				SetupProxy()
				if err := InstallCommentFont(args[i]); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
			// If it's a known comment font def name, resolve it
			if cfd := GetCommentFontDef(cfArg); cfd != nil {
				CommentFontName = cfArg
			} else {
				// Raw font name (system font or path)
				CommentFontName = cfArg
			}
			continue
		}
		if arg == "-header" {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: -header requires a file path\n")
				os.Exit(1)
			}
			i++
			HeaderFile = args[i]
			continue
		}
		if arg == "-footer" {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: -footer requires a file path\n")
				os.Exit(1)
			}
			i++
			FooterFile = args[i]
			continue
		}
		filteredArgs = append(filteredArgs, arg)
	}

	if len(filteredArgs) == 1 {
		arg := filteredArgs[0]
		if strings.HasPrefix(arg, "-") && arg != "-" && len(arg) > 1 {
			if arg == "-v" {
				PrintUsageAndExit(true)
			}
			PrintUsageAndExit(false)
		}
		inputFileName = arg
		outputFileName = makeOutputFileName(arg)
	} else if len(filteredArgs) == 2 && strings.HasPrefix(filteredArgs[0], "-") { // e.g. -35 file.c
		arg0 := filteredArgs[0]
		if arg0 == "-v" {
			PrintUsageAndExit(true)
		}

		numStr := arg0[1:]
		if val, err := strconv.Atoi(numStr); err == nil {
			PageLenMax = val
			if PageLenMax > 1024 {
				PageLenMax = 0
			}
		}

		inputFileName = filteredArgs[1]
		outputFileName = makeOutputFileName(filteredArgs[1])
	} else {
		PrintUsageAndExit(false)
	}
}

func makeOutputFileName(fname string) string {
	ext := filepath.Ext(fname)
	base := strings.TrimSuffix(fname, ext)
	if ext != "" {
		base = base + ext
	}
	return base + ".tex"
}

func PrintUsageAndExit(versionOnly bool) {
	if versionOnly {
		fmt.Printf("src2tex-go version %s with %d byte buffer\n", Version, BufferSize)
		os.Exit(1)
	}
	fmt.Printf("src2tex-go version %s — Source code to TeX/LaTeX converter\n\n", Version)
	fmt.Printf("Usage:\n")
	fmt.Printf("  src2latex [options] <file>              Convert source to LaTeX (default)\n")
	fmt.Printf("  src2tex  [options] <file>               Convert source to plain TeX\n")
	fmt.Printf("  src2latex -font list                    List available CJK fonts\n")
	fmt.Printf("  src2latex -font install <name|all>      Download and install a CJK font\n")
	fmt.Printf("\nOptions:\n")
	fmt.Printf("  -latex             Explicitly select LaTeX mode\n")
	fmt.Printf("  -tex               Explicitly select plain TeX mode\n")
	fmt.Printf("  -unicode           Unicode output for XeLaTeX/LuaLaTeX/Tectonic (default)\n")
	fmt.Printf("  -legacy            Legacy pTeX/pLaTeX output (disables Unicode mode)\n")
	fmt.Printf("  -euc               Input encoding: EUC-JP\n")
	fmt.Printf("  -sjis              Input encoding: Shift_JIS\n")
	fmt.Printf("  -utf8              Input encoding: UTF-8\n")
	fmt.Printf("  -font <name>       Select CJK font (default: ipaex)\n")
	fmt.Printf("  -commentfont <name> CJK font for comments (use -commentfont list to see options)\n")
	fmt.Printf("  -fontdir <path>    Custom font directory (default: ~/.src2tex/fonts/)\n")
	fmt.Printf("  -proxy <url>       HTTP proxy for font downloads (or set HTTP_PROXY env var)\n")
	fmt.Printf("  -<n>               Lines per page with line numbers (e.g. -35)\n")
	fmt.Printf("  -0                 Line numbers without page limit\n")
	fmt.Printf("  -header <file>     Custom header preamble (replaces fancyhdr header)\n")
	fmt.Printf("  -footer <file>     Custom footer preamble (replaces fancyhdr footer)\n")
	fmt.Printf("  -v                 Show version information\n")
	fmt.Printf("\nSupported languages:\n")
	fmt.Printf("  *.tex, *.txt, *.doc         =>    TEXT\n")
	fmt.Printf("  *.bas, *.vb                 =>    BASIC\n")
	fmt.Printf("  *.c, *.cpp, *.vc, *.h, *.hpp =>   C, C++, OBJECTIVE-C\n")
	fmt.Printf("  *.cbl, *.cob                =>    COBOL\n")
	fmt.Printf("  *.f, *.for                  =>    FORTRAN\n")
	fmt.Printf("  *.html                      =>    HTML\n")
	fmt.Printf("  *.java                      =>    JAVA\n")
	fmt.Printf("  *.el, *.lsp, *.sc, *.scm    =>    LISP, SCHEME\n")
	fmt.Printf("  makefile                    =>    MAKE\n")
	fmt.Printf("  *.p, *.pas, *.tp            =>    PASCAL\n")
	fmt.Printf("  *.pl, *.prl                 =>    PERL\n")
	fmt.Printf("  *.sh, *.csh, *.ksh          =>    SHELL\n")
	fmt.Printf("  *.tcl, *.tk                 =>    TCL/TK\n")
	fmt.Printf("  *.py                        =>    PYTHON\n")
	fmt.Printf("  *.rb                        =>    RUBY\n")
	fmt.Printf("  *.rs                        =>    RUST\n")
	fmt.Printf("  *.go                        =>    GO\n")
	fmt.Printf("  *.js                        =>    JAVASCRIPT\n")
	fmt.Printf("  *.ts                        =>    TYPESCRIPT\n")
	fmt.Printf("  *.kt                        =>    KOTLIN\n")
	fmt.Printf("  *.swift                     =>    SWIFT\n")
	fmt.Printf("  *.asi, *.asir, *.asr        =>    ASIR\n")
	fmt.Printf("  *.mac, *.max                =>    MACSYMA, MAXIMA\n")
	fmt.Printf("  *.map, *.mpl                =>    MAPLE\n")
	fmt.Printf("  *.mat, *.mma                =>    MATHEMATICA\n")
	fmt.Printf("  *.ml, *.mtlb, *.oct         =>    MATLAB, OCTAVE\n")
	fmt.Printf("  *.mu                        =>    MUPAD\n")
	fmt.Printf("  *.red, *.rdc                =>    REDUCE\n")
	fmt.Printf("\nBuilt-in CJK fonts (use -font <name>):\n")
	fmt.Printf("  ipaex     IPAex Gothic       TeX Live bundled (default, no download needed)\n")
	fmt.Printf("  hackgen   HackGen            Hack + 源柔ゴシック (unified, SIL OFL)\n")
	fmt.Printf("  udev      UDEV Gothic        JetBrains Mono + BIZ UDゴシック (unified, SIL OFL)\n")
	fmt.Printf("  firple    Firple             Fira Code + IBM Plex Sans JP (unified, SIL OFL)\n")
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  src2latex sample.c                     # Basic conversion (Unicode/LaTeX)\n")
	fmt.Printf("  src2latex -euc newton.c                # EUC-JP source\n")
	fmt.Printf("  src2latex -font hackgen sample.rb      # Use HackGen font\n")
	fmt.Printf("  src2latex -font firple -commentfont HaranoAjiMincho-Regular sample.c\n")
	fmt.Printf("                                         # Mincho for comments\n")
	fmt.Printf("  src2latex -legacy sample.c             # Legacy pLaTeX output\n")
	fmt.Printf("  src2latex -0 sample.c                  # With line numbers\n")
	os.Exit(1)
}

func OpenFiles() {
	if inputFileName == "" {
		inputFile = os.Stdin
	} else {
		var err error
		inputFile, err = os.Open(inputFileName)
		if err != nil {
			log.Fatalf("\nError: cannot open %s\n", inputFileName)
		}
	}

	if outputFileName == "" {
		outputFile = os.Stdout
	} else {
		var err error
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			log.Fatalf("\nError: cannot open %s\n", outputFileName)
		}
	}
}

func CloseFiles() {
	if inputFileName != "" && inputFile != nil {
		inputFile.Close()
	}
	if outputFileName != "" && outputFile != nil {
		outputFile.Close()
	}
}
