# src2tex-go

`src2tex-go` is a fast and simple utility that converts source code from various programming languages into cleanly formatted TeX / LaTeX documents. This is a Go port and extension of the long-standing C utility `src2tex` (v2.12), developed with a focus on maintaining as much backward compatibility as possible.

## Features

- **Multi-language support**: Automatically detects and formats comments and source code structures for 30+ languages including C, C++, Java, Pascal, Lisp, Scheme, BASIC, Fortran, Perl, Tcl/Tk, and more.
- **Modern language additions**: Go port adds support for Python, Ruby, Rust, Go, JavaScript, TypeScript, Kotlin, and Swift.
- **TeX and LaTeX modes**: Built-in support for both standard plain TeX mode and LaTeX mode output.
- **Unicode mode (default)**: Produces XeLaTeX / LuaLaTeX / Tectonic compatible output by default, with automatic CJK font configuration. Generated `.tex` files use the `iftex` package for automatic engine detection, so the same file compiles with any supported engine.
- **CJK font management**: The `-font` option allows selecting from multiple Japanese monospace fonts including IPAex Gothic (default), HackGen, UDEV Gothic, and Firple. Built-in font download and installation from GitHub.
- **Automatic EPS conversion**: Converts legacy `\special{epsfile=...}` commands to `\includegraphics` and auto-converts EPS files to PDF using Ghostscript.
- **Backward compatible**: Identical behavior and output mapping to the original `src2tex-212` C application (including implicit Japanese pTeX style mappings).

## Naming Convention

This project uses the following naming scheme:

| Name | Purpose |
|---|---|
| `src2tex-go` | Project name (Go port of src2tex) |
| `src2latex` | Binary name (automatically operates in LaTeX mode) |
| `src2tex` | Rename/symlink the binary to this name for plain TeX mode |

> **Note**: The legacy name `src2latexg` is still recognized for backward compatibility.

## Prerequisites

| Tool | Purpose | Required |
|---|---|---|
| [Go](https://go.dev/dl/) 1.21+ | Building from source | Required |
| [Task (go-task)](https://taskfile.dev/installation/) | Task runner | Optional (manual build also works) |
| [XeLaTeX](https://tug.org/xetex/), [LuaLaTeX](https://www.luatex.org/), or [Tectonic](https://tectonic-typesetting.github.io/) | PDF generation | Any one required for PDF output |
| [Ghostscript](https://www.ghostscript.com/) (`gs`) | EPS→PDF conversion | Required for samples with figures |

The only Go dependency beyond the standard library is `golang.org/x/text` for encoding conversion.

## How to Compile

### Using go-task (Recommended)

If you have [go-task](https://taskfile.dev/) installed, you can use the included `Taskfile.yml` for building, converting samples, and generating PDFs. go-task is a cross-platform task runner — the same commands work on Windows, macOS, and Linux.

```bash
# Build for the current platform
task build

# Cross-compile for all platforms
task build:all

# Build for a specific platform
task build:darwin-arm64    # macOS (Apple Silicon)
task build:darwin-amd64    # macOS (Intel)
task build:linux-amd64     # Linux (amd64)
task build:windows-amd64   # Windows (amd64)

# Convert all samples to TeX
task samples

# Generate PDFs from samples (converts + compiles)
task pdf

# Use Tectonic instead of XeLaTeX
task pdf:compile TEX_ENGINE=tectonic

# Cleanup
task clean          # Remove binaries only
task clean:samples  # Remove generated TeX/PDF files
task clean:all      # Remove everything
```

### Manual Build

You can compile this project using standard Go build tools:

```bash
go build -o src2latex
```

### Cross-compilation (Manual)

The Go toolchain easily produces native, standalone cross-platform binaries:
```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o src2latex-darwin-arm64
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o src2latex-darwin-amd64
# Windows
GOOS=windows GOARCH=amd64 go build -o src2latex-win-amd64.exe
# Linux
GOOS=linux GOARCH=amd64 go build -o src2latex-linux-amd64
```

## Usage

### Basic Usage

The default binary name is `src2latex`. It automatically operates in **LaTeX mode** with **Unicode mode** enabled.

```bash
./src2latex <input file>
```

The converted file will have a `.tex` extension appended to its name.
- Example: `hanoi.c` → `hanoi.c.tex`

### Mode Switching

The output mode is determined automatically by the binary name:

| Binary name | Default mode |
|---|---|
| `src2latex` | **LaTeX mode** (default build name) |
| `src2tex` | **plain TeX mode** (use via symlink, etc.) |

```bash
# To use plain TeX mode
ln -s src2latex src2tex
./src2tex <input file>
```

You can also override the mode explicitly with the `-latex` / `-tex` flags, regardless of the binary name.

### Options

| Option | Description |
|---|---|
| `-latex` | Explicitly select LaTeX mode |
| `-tex` | Explicitly select plain TeX mode |
| `-unicode` | Unicode output mode for XeLaTeX / LuaLaTeX / Tectonic (**default**) |
| `-legacy` | Legacy pTeX / pLaTeX output mode (disables Unicode mode) |
| `-euc` | Specify input file encoding as EUC-JP |
| `-sjis` | Specify input file encoding as Shift_JIS |
| `-utf8` | Specify input file encoding as UTF-8 |
| `-font <name>` | Select CJK font (default: `ipaex`) |
| `-commentfont <name>` | Select CJK font for comments (e.g., mincho/serif) |
| `-fontdir <path>` | Custom font installation directory (default: `~/.src2tex/fonts/`) |
| `-proxy <url>` | HTTP proxy for font downloads (or set HTTP_PROXY env var) |
| `-<n>` | Limit output to `n` lines per page (with line numbers) |
| `-0` | Output with line numbers (no page limit) |
| `-header <file>` | Custom header preamble file (replaces fancyhdr header section) |
| `-footer <file>` | Custom footer preamble file (replaces fancyhdr footer section) |
| `-v` | Display version information |

### Examples

```bash
# Basic conversion (Unicode/LaTeX mode, default)
./src2latex samples/hanoi.c

# EUC-JP source
./src2latex -euc samples/newton.c

# Use HackGen font
./src2latex -font hackgen samples/hanoi.rb

# Use mincho font for Japanese in comments
./src2latex -font firple -commentfont haranoaji samples/hanoi.c

# With line numbers
./src2latex -0 samples/hanoi.c

# Custom header/footer
./src2latex -header my_header.tex -footer my_footer.tex samples/hanoi.c

# Legacy pLaTeX output
./src2latex -legacy samples/hanoi.c
```

### Generating PDFs

Compile the generated TeX files with your preferred TeX engine.
Files generated in the default Unicode mode have been verified to work with all three engines below:

```bash
# Using XeLaTeX (included in TeX Live / MacTeX)
xelatex hanoi.c.tex

# Using LuaLaTeX (included in TeX Live / MacTeX)
lualatex hanoi.c.tex

# Using Tectonic (lightweight, easy to install)
tectonic hanoi.c.tex
```
While the original src2tex was designed for older versions of TeX and LaTeX from the 1990s, this version has been updated to generate preambles suitable for modern environments. By using the -font and -commentfont options, you can assign your preferred Japanese typefaces to both code and comment sections.

Please note that the -font option for code sections assumes the use of fonts where full-width characters are exactly twice the width of half-width characters. In this mode, instead of using boxes for indentation, the output uses literal spaces or tabs.

> **Note**: Generated `.tex` files use the `iftex` package for automatic engine detection. The same `.tex` file can be compiled by any of the above engines without modification.
> - **XeLaTeX**: Uses `fontspec` + `xeCJK` + `zxjafont`
> - **LuaLaTeX**: Uses `luatexja-preset`
> - **Tectonic**: Uses the XeTeX engine internally, so the XeLaTeX packages apply

```bash
# Using traditional pLaTeX (for files converted with -legacy)
platex hanoi.c.tex
dvipdfmx hanoi.c.dvi
```

## CJK Font Management

The `-font` option lets you switch the CJK font used for comments and Japanese text.

### Built-in Fonts

| Name | Font | Type | License | Description |
|---|---|---|---|---|
| `ipaex` | IPAex Gothic | Non-unified | IPA | Bundled with TeX Live. No download needed (**default**) |
| `hackgen` | HackGen | Unified | SIL OFL | Hack + Gen Jyuu Gothic. Popular programming font |
| `udev` | UDEV Gothic | Unified | SIL OFL | JetBrains Mono + BIZ UD Gothic. Universal Design |
| `firple` | Firple | Unified | SIL OFL | Fira Code + IBM Plex Sans JP. Ligature support. Half:Full = 1:2 |

**Unified fonts** (HackGen, UDEV Gothic, Firple) have consistent Latin/CJK character widths, resulting in accurate text selection and copy-paste in the PDF output.

### Installing and Using Fonts

```bash
# List available fonts
./src2latex -font list

# Download and install fonts
./src2latex -font install hackgen       # Install HackGen
./src2latex -font install all           # Install all fonts

# Convert using installed font
./src2latex -font hackgen samples/hanoi.rb

# Custom font directory
./src2latex -fontdir /path/to/fonts -font hackgen samples/hanoi.rb
```

Fonts are installed to `~/.src2tex/fonts/` by default. Override with `-fontdir`.

### Custom Fonts

Add custom font definitions in `~/.src2tex/fonts.json`:

```json
{
  "fonts": [
    {
      "Name": "myFont",
      "DisplayName": "My Custom Font",
      "License": "MIT",
      "Unified": true,
      "RegularFile": "MyFont-Regular.ttf",
      "BoldFile": "MyFont-Bold.ttf",
      "Description": "Custom font description"
    }
  ]
}
```

## Comment Font Management

The `-commentfont` option lets you specify a separate CJK font for comment text. When using unified fonts (HackGen, Firple, etc.) for code, comments appear in gothic/sans-serif by default. By specifying a mincho/serif font for comments, you get elegant typesetting that harmonizes with Computer Modern.

### Available Comment Fonts

| Name | Font | License | Description |
|---|---|---|---|
| `haranoaji` | Harano Aji Mincho | SIL OFL | Bundled with TeX Live. No download needed |
| `ipaexm` | IPAex Mincho | IPA | Bundled with TeX Live. No download needed |
| `noto-serif` | Noto Serif JP | SIL OFL | Google Noto serif font. Requires download |

### Installing and Using Comment Fonts

```bash
# List available comment fonts
./src2latex -commentfont list

# Download and install comment fonts
./src2latex -commentfont install noto-serif    # Install Noto Serif JP
./src2latex -commentfont install all           # Install all comment fonts

# Use TeX Live bundled fonts (no download needed)
./src2latex -font firple -commentfont haranoaji samples/hanoi.c
./src2latex -font hackgen -commentfont ipaexm samples/hanoi.rb

# Use downloaded font
./src2latex -font firple -commentfont noto-serif samples/hanoi.c
```

> **How it works**: The `-commentfont` font is set as `\setCJKmainfont` (XeLaTeX) / `\setmainjfont` (LuaLaTeX). When comments switch to `\rm` (roman) mode, Japanese text uses this serif font. Code sections continue using the gothic font via `\setCJKmonofont`.

### Downloading Fonts Through a Proxy

If you access the internet through a proxy server, use the `-proxy` option or environment variables:

```bash
# Using -proxy option (applies to both -font install and -commentfont install)
./src2latex -proxy http://proxy.example.com:8080 -font install hackgen
./src2latex -proxy http://proxy.example.com:8080 -commentfont install noto-serif

# Authenticated proxy
./src2latex -proxy http://user:pass@proxy.example.com:8080 -font install hackgen

# Using environment variables (Go's http package auto-detects these)
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
./src2latex -font install all
./src2latex -commentfont install all
```

## Supported Languages

Language detection is based on file extension. When the extension is ambiguous, keyword-based detection is attempted.

### Languages from original src2tex-212

| Extension | Language |
|---|---|
| `.tex`, `.txt`, `.doc` | TEXT |
| `.bas`, `.vb` | BASIC |
| `.c`, `.cpp`, `.vc`, `.h`, `.hpp` | C, C++, Objective-C |
| `.cbl`, `.cob` | COBOL |
| `.f`, `.for` | FORTRAN |
| `.html` | HTML |
| `.java` | Java |
| `.el`, `.lsp`, `.sc`, `.scm` | Lisp, Scheme |
| `makefile` | Make |
| `.p`, `.pas`, `.tp` | Pascal |
| `.pl`, `.prl` | Perl |
| `.sh`, `.csh`, `.ksh` | Shell |
| `.tcl`, `.tk` | Tcl/Tk |
| `.asi`, `.asir`, `.asr` | Asir |
| `.mac`, `.max` | Macsyma, Maxima |
| `.map`, `.mpl` | Maple |
| `.mat`, `.mma` | Mathematica |
| `.ml`, `.mtlb`, `.oct` | MATLAB, Octave |
| `.mu` | MuPAD |
| `.red`, `.rdc` | REDUCE |

### Languages added in Go port

| Extension | Language |
|---|---|
| `.py` | Python |
| `.rb` | Ruby |
| `.rs` | Rust |
| `.go` | Go |
| `.js` | JavaScript |
| `.ts` | TypeScript |
| `.kt` | Kotlin |
| `.swift` | Swift |

## FAQ

### Q1: How do I write comments in src2tex format?

You can write TeX math and commands directly within source code comments. For example, in C:

```c
r = sqrt(x*x + y*y);   /* radius $r=\sqrt{x^2+y^2}$ */
```

In REDUCE:
```reduce
int(x/sqrt(1-x^2), x);  % integration $\int{x\over\sqrt{1-x^2}}\,dx$
```

### Q2: How do I include EPS/PS figures?

Use the `\special{epsfile=...}` command within TeX-mode comment areas:

```c
/* See the following numerical simulation.
                    {\special{epsfile=simulation.eps}} */
```

In Unicode mode (default), `\special{epsfile=...}` is automatically converted to `\includegraphics`, and EPS files are auto-converted to PDF via Ghostscript.

### Q3: Can I change comment area fonts?

Yes, use the `\src2tex{...}` escape sequence:

```c
/* {\src2tex{texfont=tt}} */     /* Switch to typewriter font */
/* {\src2tex{texfont=rm}} */     /* Switch back to roman font */
/* {\src2tex{texfont=bf}} */     /* Bold font */
/* {\src2tex{texfont=it}} */     /* Italic font */
```

Program area fonts can also be changed:
```c
/* {\src2tex{textfont=bf}} */    /* Bold program text */
```

### Q4: Can I change tab/indent width?

```c
/* {\src2tex{htab=4}} */         /* Set horizontal tab size to 4 */
/* {\src2tex{vtab=2}} */         /* Set vertical tab size to 2 */
```

### Q5: Can I use a different LaTeX style file?

Write the desired documentstyle in a comment at the beginning of your source file:

```pascal
(* {\documentstyle[twocolumn,12pt]{article}} *)
```

Note: In Unicode mode (default), a XeLaTeX-compatible preamble is always used, so in-source documentstyle directives are ignored. This feature is effective when converting with `-legacy` mode.

### Q6: How do I output line numbers?

Use the `-0` or `-<n>` options:

```bash
./src2latex -0 samples/hanoi.c       # Line numbers (no limit)
./src2latex -35 samples/hanoi.c      # 35 lines per page with numbers
```

### Q7: Can I customize the page header and footer?

Yes, use the `-header` / `-footer` options to replace the default fancyhdr header/footer configuration with your own custom files.

First, create a custom header file (e.g., `my_header.tex`):

```tex
\renewcommand{\headrulewidth}{0.4pt}
\fancyhf{}
\fancyhead[L]{\rm My Project Name}
\fancyhead[R]{\rm \today}
```

Create a custom footer file (e.g., `my_footer.tex`) similarly:

```tex
\fancyfoot[C]{\thepage}
\fancyfoot[L]{\rm Confidential}
\fancyfoot[R]{\rm Draft}
```

Specify them during conversion:

```bash
# Customize header only
./src2latex -header my_header.tex samples/hanoi.c

# Customize footer only
./src2latex -footer my_footer.tex samples/hanoi.c

# Customize both
./src2latex -header my_header.tex -footer my_footer.tex samples/hanoi.c
```

> **Note**: `\usepackage{fancyhdr}` and `\pagestyle{fancy}` are emitted automatically, so your custom file should only contain the settings that follow (e.g., `\fancyhead`, `\fancyfoot`, `\renewcommand{\headrulewidth}`).

> **Tip**: To completely hide headers or footers, specify a file containing only `\fancyhf{}`.

### Q8: Can I make specific keywords bold?

src2tex does not have this feature built-in, but you can post-process the TeX output with `sed` or similar tools:

```bash
./src2latex sample.c
sed -e 's/}main(){/}{\\bf main()}{/g' sample.c.tex > sample_modified.tex
```

src2tex translates each keyword in program area to the form `}keyword{`, which makes post-processing straightforward.

## Notes

### Encoding

- Original src2tex-212 samples (`newton.c`, `simpson.c`, `farmer+hen.scm`) are EUC-JP encoded. Use `-euc` when converting them.
- Samples created for the Go port (`hanoi.c`, `hanoi.go`, `hanoi.py`, etc.) are UTF-8. Use `-utf8` or leave as default.
- With `-legacy` mode, output is formatted for pTeX / pLaTeX.

### Unicode Mode

- Unicode mode is enabled by default. Use `-legacy` to disable it.
- The Unicode mode generates a preamble compatible with XeLaTeX, LuaLaTeX, and Tectonic.
- The `iftex` package is used for automatic engine detection: XeLaTeX uses `xeCJK` + `zxjafont` (or custom font settings), while LuaLaTeX uses `luatexja-preset`. This allows the same `.tex` file to compile with any supported engine.
- Compatibility macros for plain TeX commands (`\eqalign`, `\pmatrix`, etc.) are automatically included.
- `\special{epsfile=...}` directives are auto-converted to `\includegraphics`.
- EPS files are auto-converted to PDF via Ghostscript (`gs` command required).

#### Supported TeX Engines

| Engine | Description | Installation |
|---|---|---|
| [XeLaTeX](https://tug.org/xetex/) | Included in TeX Live / MacTeX. Unicode + OpenType font support | `brew install --cask mactex` (macOS) |
| [LuaLaTeX](https://www.luatex.org/) | Included in TeX Live / MacTeX. Advanced font control via Lua | Bundled with TeX Live |
| [Tectonic](https://tectonic-typesetting.github.io/) | Lightweight XeTeX-based engine with automatic package download | `brew install tectonic` (macOS) |

### Sample Files

The `samples/` directory contains:

| File | Language | Content |
|---|---|---|
| `hanoi.c` | C | Towers of Hanoi (English/Japanese comments) |
| `newton.c` | C | Newton-Raphson method (with math & figures) |
| `simpson.c` | C | Simpson's rule (with math & figures) |
| `hanoi.go` | Go | Towers of Hanoi in Go |
| `hanoi.py` | Python | Towers of Hanoi in Python |
| `hanoi.rb` | Ruby | Towers of Hanoi in Ruby |
| `hanoi.rs` | Rust | Towers of Hanoi in Rust |
| `hanoi.js` | JavaScript | Towers of Hanoi in JavaScript |
| `hanoi.ts` | TypeScript | Towers of Hanoi in TypeScript |
| `hanoi.kt` | Kotlin | Towers of Hanoi in Kotlin |
| `hanoi.swift` | Swift | Towers of Hanoi in Swift |
| `farmer+hen.scm` | Scheme | Farmer and Hen puzzle |
| `popgen.red` | REDUCE | Population genetics PDE |
| `sqrt_mat.red` | REDUCE | Square root of matrix |

### Quick Start with Samples

```bash
# UTF-8 samples (created for Go port) — Unicode mode (default)
./src2latex samples/hanoi.c
./src2latex samples/hanoi.go
./src2latex samples/hanoi.py
./src2latex -font hackgen samples/hanoi.rb    # With HackGen font
./src2latex -font firple -commentfont haranoaji samples/hanoi.c  # Mincho comments

# EUC-JP samples (from original src2tex-212)
./src2latex -euc samples/newton.c
./src2latex -euc samples/simpson.c
./src2latex -euc samples/farmer+hen.scm

# ASCII samples
./src2latex samples/popgen.red
./src2latex samples/sqrt_mat.red

# Generate PDFs (using any supported engine)
tectonic samples/hanoi.c.tex     # Tectonic
xelatex samples/hanoi.c.tex      # XeLaTeX
lualatex samples/hanoi.c.tex     # LuaLaTeX
```

## About the Versioning

As a play on words between the Go language and the Japanese word for the number five ("go"), I've used the decimal expansion of $ \sqrt 5 $
 for the version numbers. The current version is 2.23.

## Credits

The original `src2tex` version 2.12 was developed by Kazuo AMANO and Shinichi NOMOTO.

The original sample files (and related resources) included here are the copyrighted property of the original author.

- samples/farmer+hen.scm
- samples/hanoi.c
- samples/newton.c
- samples/popgen.red
- samples/simpson.c
- samples/sqrt_mat.red

> ** NOTE **: The file hanoi.c was originally included as hanoi89.c, but it has been adjusted to be compatible with C11. This version served as the basis for the hanoi implementation in other languages.

## License

Subject to the license terms of the original src2tex-212.
