package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// 1. Initialize IsLatexMode based on executable name
	execName := filepath.Base(os.Args[0])
	if execName == "src2latex" || execName == "src2latex.exe" || execName == "src2latexg" {
		IsLatexMode = true
	} else {
		IsLatexMode = false
	}

	// 2. Also check for explicit -latex / -tex flag
	for _, arg := range os.Args[1:] {
		if arg == "-latex" || arg == "-tex" {
			IsLatexMode = true
		}
	}

	// 3. Parse arguments and configure input/output files
	GetFileNames()

	// 4. Initialize Language Flags based on extension and contents
	InitLangFlag()

	// 5. Open associated files (sets inputFile and outputFile)
	OpenFiles()

	// 6. Execute core translation sequence
	cptr := []string{inputFileName, outputFileName}
	fptr := []*os.File{inputFile, outputFile}

	Text2Tex(cptr, fptr)

	// 7. Cleanup
	CloseFiles()

	// 8. Post-process: convert legacy \special{epsfile=...} to \includegraphics for Unicode TeX
	if IsUnicodeTeX && outputFileName != "" {
		postProcessEpsSpecials(outputFileName)
		postProcessCommentQuote(outputFileName)
	}

	os.Exit(0)
}

// postProcessEpsSpecials reads the output .tex file and converts dvips-style
// \special{epsfile=<file> hscale=<h> vscale=<v> hoffset=<x>} to
// \includegraphics[scale=<s>]{<file>} for XeLaTeX/Tectonic compatibility.
func postProcessEpsSpecials(texFile string) {
	data, err := os.ReadFile(texFile)
	if err != nil {
		return
	}
	content := string(data)

	// Match \special{epsfile=<filename> ...optional params...}
	re := regexp.MustCompile(`\\special\{epsfile=([^\s}]+)([^}]*)\}`)

	if !re.MatchString(content) {
		return // nothing to do
	}

	// Convert referenced EPS files to PDF with correct BoundingBox
	texDir := filepath.Dir(texFile)
	for _, submatch := range re.FindAllStringSubmatch(content, -1) {
		if len(submatch) >= 2 {
			epsName := submatch[1]
			epsPath := filepath.Join(texDir, epsName)
			pdfPath := filepath.Join(texDir, strings.TrimSuffix(epsName, ".eps")+".pdf")
			convertEpsToPdf(epsPath, pdfPath)
		}

	}

	content = re.ReplaceAllStringFunc(content, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		epsFileName := submatches[1]
		// Strip .eps extension so graphicx auto-discovers the .pdf version
		epsFileName = strings.TrimSuffix(epsFileName, ".eps")
		params := ""
		if len(submatches) > 2 {
			params = submatches[2]
		}

		// Parse hscale, vscale, hoffset from params
		hscaleRe := regexp.MustCompile(`hscale=([0-9.]+)`)
		vscaleRe := regexp.MustCompile(`vscale=([0-9.]+)`)
		hoffsetRe := regexp.MustCompile(`hoffset=([0-9.]+)`)

		hscaleMatch := hscaleRe.FindStringSubmatch(params)
		vscaleMatch := vscaleRe.FindStringSubmatch(params)
		hoffsetMatch := hoffsetRe.FindStringSubmatch(params)

		// dvips hscale/vscale are direct scaling factors for the EPS image
		// (1.0 = natural size). Use scale= for graphicx.
		var opts []string
		if len(hscaleMatch) > 1 && len(vscaleMatch) > 1 {
			opts = append(opts, fmt.Sprintf("scale=%s", vscaleMatch[1]))
		} else if len(hscaleMatch) > 1 {
			opts = append(opts, fmt.Sprintf("scale=%s", hscaleMatch[1]))
		} else if len(vscaleMatch) > 1 {
			opts = append(opts, fmt.Sprintf("scale=%s", vscaleMatch[1]))
		}

		optStr := ""
		if len(opts) > 0 {
			optStr = "[" + strings.Join(opts, ",") + "]"
		}

		result := fmt.Sprintf("\\includegraphics%s{%s}", optStr, epsFileName)

		// If there's a hoffset, use \hfill to space from previous image
		// (dvips hoffset is absolute positioning; in LaTeX flow we approximate with \hfill)
		if len(hoffsetMatch) > 1 {
			result = fmt.Sprintf("\\hfill %s", result)
		}

		return result
	})

	// When dvips used hoffset to place a second image beside the first,
	// we get: \includegraphics...\n\hfill \includegraphics...
	// Wrap them in \begin{center}...\quad...\end{center} for centered
	// side-by-side layout instead of spreading to page edges.
	sideBySideRe := regexp.MustCompile(`(\\includegraphics[^\n]+)\n\\hfill\s+(\\includegraphics[^\n]+)`)
	content = sideBySideRe.ReplaceAllString(content, "\\begin{center}\n${1}\\quad ${2}\n\\end{center}")

	// For standalone \hfill \includegraphics (not preceded by another image),
	// center the image. Also handle lines starting with \mbox{} (from newline handling).
	hfillImgRe := regexp.MustCompile(`(?m)^(\s*(?:\\mbox\{\}\s*)?)\\hfill\s+(\\includegraphics[^\n]+)`)
	content = hfillImgRe.ReplaceAllString(content, "\\begin{center}\n${2}\n\\end{center}")

	// Post-process \vskip adjustments for dvips→includegraphics conversion.
	// In dvips, \special overlays images without consuming vertical space,
	// so the original TeX uses \vskip to reserve space. With \includegraphics,
	// the image height is automatically consumed, so large \vskip values must
	// be replaced with small spacing appropriate for inline images.
	//
	// Pattern: \vskip <large>cm after last includegraphics → remove
	vskipAfterRe := regexp.MustCompile(`(\\includegraphics[^\n]*\n)\\vskip\s+[0-9.]+\s*cm`)
	content = vskipAfterRe.ReplaceAllString(content, "${1}")
	// Also handle \vskip after \end{center}
	vskipAfterCenterRe := regexp.MustCompile(`(\\end\{center\}\n)\\vskip\s+[0-9.]+\s*cm`)
	content = vskipAfterCenterRe.ReplaceAllString(content, "${1}")

	// Pattern: \vskip <N>cm right before includegraphics or \begin{center} → remove
	vskipBeforeRe := regexp.MustCompile(`\\vskip\s+[0-9.]+\s*cm\n(\\includegraphics|\\begin\{center\})`)
	content = vskipBeforeRe.ReplaceAllString(content, "${1}")

	// When /* comment marker and figure block are in the same paragraph,
	// use \par to end the paragraph so /* appears above the figure.
	inlineFigRe := regexp.MustCompile(`(\\kern[0-9.]+em\s+)(\\begin\{center\})`)
	content = inlineFigRe.ReplaceAllString(content, "${1}\\par\n${2}")
	// Also handle bare \includegraphics (not wrapped in center)
	inlineFigRe2 := regexp.MustCompile(`(\\kern[0-9.]+em\s+)(\\includegraphics[^\n]+)`)
	content = inlineFigRe2.ReplaceAllString(content, "${1}\\par\n\\begin{center}\n${2}\n\\end{center}")

	// Remove the blank \mbox{}\hfill line immediately before /* + figure lines
	// to reduce vertical gap (the blank line is from source but creates excess space with figures)
	blankBeforeFigRe := regexp.MustCompile(`\\noindent\n\\mbox\{\}\\hfill\n\n\\noindent\n(\\mbox\{\}\\rm\\mc\s+\{\\tt /\}\{\\tt \*\})`)
	content = blankBeforeFigRe.ReplaceAllString(content, "\\noindent\n${1}")

	os.WriteFile(texFile, []byte(content), 0644)
}

// postProcessCommentQuote fixes the issue where \begin{quote} and \end{quote}
// inside comment blocks (lines with {\tt\#}) shift the # marker to the right.
// Instead of using LaTeX's quote environment (which would indent # too), we
// remove the quote environment and manually add indentation to the text content
// inside the quote block. This keeps # at the left margin while the text gets
// the expected indentation.
func postProcessCommentQuote(texFile string) {
	data, err := os.ReadFile(texFile)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	inQuote := false
	var result []string

	quoteStartRe := regexp.MustCompile(`(\{\\tt\\#\}[^\n]*)\{\}\\begin\{quote\}`)
	quoteEndRe := regexp.MustCompile(`(\{\\tt\\#\}[^\n]*)\{\}\\end\{quote\}`)
	// Match comment text lines: {\\tt\\#} followed by spacing ({} or \\kern1\\charwd or \\ )
	commentTextRe := regexp.MustCompile(`(\{\\tt\\#\}(?:\\kern1\\charwd|\\ )\s+)\{\}`)

	for _, line := range lines {
		if quoteStartRe.MatchString(line) {
			// Remove \begin{quote} from the line, enter quote mode
			line = quoteStartRe.ReplaceAllString(line, "${1}{}")
			result = append(result, line)
			inQuote = true
			continue
		}
		if quoteEndRe.MatchString(line) {
			// Remove \end{quote} from the line, exit quote mode
			line = quoteEndRe.ReplaceAllString(line, "${1}{}")
			inQuote = false
			result = append(result, line)
			continue
		}
		if inQuote && commentTextRe.MatchString(line) {
			// Add manual indentation after {} for quote block content
			line = commentTextRe.ReplaceAllString(line, "${1}{}\\hspace{\\leftmargini}")
		}
		result = append(result, line)
	}

	os.WriteFile(texFile, []byte(strings.Join(result, "\n")), 0644)
}

// convertEpsToPdf converts an EPS file to PDF using Ghostscript,
// preserving the original BoundingBox to avoid full-page whitespace.
// If the PDF already exists, it checks whether it has the correct dimensions
// and only reconverts if necessary.
func convertEpsToPdf(epsPath, pdfPath string) {
	// Read BoundingBox from EPS file
	epsFile, err := os.Open(epsPath)
	if err != nil {
		return // EPS file not found; skip silently
	}
	defer epsFile.Close()

	var bbWidth, bbHeight int
	scanner := bufio.NewScanner(epsFile)
	bbRe := regexp.MustCompile(`^%%BoundingBox:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		if m := bbRe.FindStringSubmatch(line); len(m) == 5 {
			x1, _ := strconv.Atoi(m[1])
			y1, _ := strconv.Atoi(m[2])
			x2, _ := strconv.Atoi(m[3])
			y2, _ := strconv.Atoi(m[4])
			bbWidth = x2 - x1
			bbHeight = y2 - y1
			break
		}
		// Stop searching after PostScript comments end
		if len(line) > 0 && line[0] != '%' {
			break
		}
	}

	if bbWidth <= 0 || bbHeight <= 0 {
		return // No valid BoundingBox found
	}

	// Check if PDF already has correct dimensions
	if existingData, err := os.ReadFile(pdfPath); err == nil {
		mbRe := regexp.MustCompile(`/MediaBox\s*\[\s*0\s+0\s+(\d+)\s+(\d+)\s*\]`)
		if m := mbRe.FindStringSubmatch(string(existingData)); len(m) == 3 {
			existW, _ := strconv.Atoi(m[1])
			existH, _ := strconv.Atoi(m[2])
			if existW == bbWidth && existH == bbHeight {
				return // Already correctly cropped
			}
		}
	}

	// Find ghostscript
	gsPath, err := exec.LookPath("gs")
	if err != nil {
		return // gs not available
	}

	// Convert EPS to PDF with correct dimensions
	cmd := exec.Command(gsPath,
		"-q", "-dNOPAUSE", "-dBATCH",
		"-sDEVICE=pdfwrite",
		fmt.Sprintf("-dDEVICEWIDTHPOINTS=%d", bbWidth),
		fmt.Sprintf("-dDEVICEHEIGHTPOINTS=%d", bbHeight),
		"-dFIXEDMEDIA",
		fmt.Sprintf("-sOutputFile=%s", pdfPath),
		epsPath,
	)
	cmd.Run()
}
