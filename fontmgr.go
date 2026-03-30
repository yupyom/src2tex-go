package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// FontDef defines a monospace font with Japanese support.
type FontDef struct {
	Name            string  // CLI option name (hackgen, udev, etc.)
	DisplayName     string  // Human-readable name
	License         string  // License type
	Unified         bool    // true = Latin+CJK unified font (no scale needed)
	RegularFile     string  // Regular weight filename (mono)
	BoldFile        string  // Bold weight filename (mono)
	MainRegularFile string  // Main font regular (for non-unified, e.g. mincho)
	MainBoldFile    string  // Main font bold
	Scale           float64 // CJK scale factor (only for non-unified, e.g. ipaex)
	GithubRepo      string  // GitHub "owner/repo" for API
	ZipAsset        string  // Substring to match the ZIP asset name
	Description     string  // Short description
}

// CJKFontName is the selected font name (default: "ipaex")
var CJKFontName = "ipaex"

// FontExplicitlySet is true when -font was explicitly specified on the command line.
// When false, -unicode mode uses the legacy fallback preamble instead of
// generating custom font settings.
var FontExplicitlySet = false

// CommentFontName is the CJK font used for comment text (\rm mode).
// When empty, the same font as CJKFontName is used for \setCJKmainfont.
// When set, this font is used for \setCJKmainfont so that Japanese text
// in comments appears in a different typeface (e.g., mincho/serif) from
// the code font (gothic/sans-serif).
var CommentFontName = ""

// UseUnifiedFont is set to true when a unified Latin+CJK equal-width font
// (e.g., HackGen, UDEV Gothic, Firple) is selected. When true, the output
// uses literal spaces ("\ ") instead of \kern commands, and omits kerning
// adjustments around underscores and quotes, so that the PDF text can be
// cleanly copy-pasted.
var UseUnifiedFont = false

// FontDir is the directory where downloaded fonts are stored.
// Default: ~/.src2tex/fonts/, overridable with -fontdir.
var FontDir string

// ProxyURL is the HTTP proxy URL for font downloads.
// When empty, Go's default behavior is used (respects HTTP_PROXY/HTTPS_PROXY env vars).
var ProxyURL string

// SetupProxy configures the HTTP client to use the specified proxy.
func SetupProxy() {
	if ProxyURL == "" {
		return
	}
	proxyParsed, err := url.Parse(ProxyURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: invalid proxy URL %q: %v\n", ProxyURL, err)
		return
	}
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyParsed),
	}
	fmt.Fprintf(os.Stderr, "Using proxy: %s\n", ProxyURL)
}

// AutoDetectCommentFont sets CommentFontName to a mincho font if one is
// available. Called automatically when CommentFontName is not explicitly set.
// Priority: haranoaji > ipaexm > downloaded comment fonts > (none = same as code font).
func AutoDetectCommentFont() {
	if CommentFontName != "" {
		return // User explicitly specified a font
	}

	// Try TeX Live bundled mincho fonts in priority order
	for _, cfd := range BuiltinCommentFonts {
		if !cfd.TexLive {
			continue
		}
		for _, dir := range texLiveFontDirs() {
			path := filepath.Join(dir, cfd.RegularFile)
			if _, err := os.Stat(path); err == nil {
				CommentFontName = cfd.Name
				return
			}
		}
	}

	// Try downloaded comment fonts
	for _, cfd := range BuiltinCommentFonts {
		if cfd.TexLive {
			continue
		}
		if IsCommentFontInstalled(cfd.Name) {
			CommentFontName = cfd.Name
			return
		}
	}

	// No mincho font found — fall back to code font (CommentFontName stays "")
}
// BuiltinFonts defines the supported font options.
var BuiltinFonts = []FontDef{
	{
		Name:            "ipaex",
		DisplayName:     "IPAex Gothic",
		License:         "IPA",
		Unified:         false,
		RegularFile:     "ipaexg.ttf",
		BoldFile:        "ipaexg.ttf",
		MainRegularFile: "ipaexm.ttf",
		MainBoldFile:    "ipaexg.ttf",
		Scale:           1.05,
		GithubRepo:      "",
		ZipAsset:        "",
		Description:     "TeX Live 同梱の等幅日本語フォント。ダウンロード不要。",
	},
	{
		Name:        "hackgen",
		DisplayName: "HackGen",
		License:     "SIL OFL",
		Unified:     true,
		RegularFile: "HackGen-Regular.ttf",
		BoldFile:    "HackGen-Bold.ttf",
		Scale:       0,
		GithubRepo:  "yuru7/HackGen",
		ZipAsset:    "HackGen_v",
		Description: "Hack + 源柔ゴシック。プログラミング向けの人気フォント。",
	},
	{
		Name:        "udev",
		DisplayName: "UDEV Gothic",
		License:     "SIL OFL",
		Unified:     true,
		RegularFile: "UDEVGothic-Regular.ttf",
		BoldFile:    "UDEVGothic-Bold.ttf",
		Scale:       0,
		GithubRepo:  "yuru7/udev-gothic",
		ZipAsset:    "UDEVGothic_v",
		Description: "JetBrains Mono + BIZ UDゴシック。UD（ユニバーサルデザイン）対応。",
	},
	{
		Name:        "firple",
		DisplayName: "Firple",
		License:     "SIL OFL",
		Unified:     true,
		RegularFile: "Firple-Regular.ttf",
		BoldFile:    "Firple-Bold.ttf",
		Scale:       0,
		GithubRepo:  "negset/Firple",
		ZipAsset:    "Firple.zip",
		Description: "Fira Code + IBM Plex Sans JP。リガチャ対応。半角:全角 = 1:2。",
	},
}

// CommentFontDef defines a serif/mincho font for comments.
type CommentFontDef struct {
	Name        string // CLI option name
	DisplayName string // Human-readable name
	License     string // License type
	TexLive     bool   // true = bundled with TeX Live (no download needed)
	RegularFile string // Regular weight filename
	BoldFile    string // Bold weight filename (may be empty)
	Extension   string // File extension (.otf or .ttf)
	FontSpec    string // fontspec name (without extension, for setCJKmainfont)
	GithubRepo  string // GitHub "owner/repo" for download
	ZipAsset    string // Substring to match the ZIP asset name
	Description string // Short description
}

// BuiltinCommentFonts defines the supported comment font options (mincho/serif).
var BuiltinCommentFonts = []CommentFontDef{
	{
		Name:        "haranoaji",
		DisplayName: "原ノ味明朝 (Harano Aji)",
		License:     "SIL OFL",
		TexLive:     true,
		RegularFile: "HaranoAjiMincho-Regular.otf",
		BoldFile:    "HaranoAjiMincho-Bold.otf",
		Extension:   ".otf",
		FontSpec:    "HaranoAjiMincho-Regular",
		Description: "TeX Live 同梱の高品質明朝体。ダウンロード不要。",
	},
	{
		Name:        "ipaexm",
		DisplayName: "IPAex 明朝",
		License:     "IPA",
		TexLive:     true,
		RegularFile: "ipaexm.ttf",
		BoldFile:    "",
		Extension:   ".ttf",
		FontSpec:    "ipaexm",
		Description: "TeX Live 同梱の明朝体。ダウンロード不要。",
	},
	{
		Name:        "noto-serif",
		DisplayName: "Noto Serif JP",
		License:     "SIL OFL",
		TexLive:     false,
		RegularFile: "NotoSerifJP-Regular.otf",
		BoldFile:    "NotoSerifJP-Bold.otf",
		Extension:   ".otf",
		FontSpec:    "NotoSerifJP-Regular",
		GithubRepo:  "notofonts/noto-cjk",
		ZipAsset:    "NotoSerifJP",
		Description: "Google Noto明朝体。高品質で幅広い字形をカバー。",
	},
}

// GetDefaultFontDir returns the default font directory (~/.src2tex/fonts/).
func GetDefaultFontDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".src2tex", "fonts")
}

// GetFontDir returns the active font directory.
func GetFontDir() string {
	if FontDir != "" {
		return FontDir
	}
	return GetDefaultFontDir()
}

// GetFontDef returns the FontDef for the given name, or nil if not found.
func GetFontDef(name string) *FontDef {
	for i := range BuiltinFonts {
		if BuiltinFonts[i].Name == name {
			return &BuiltinFonts[i]
		}
	}
	return nil
}

// ListFonts prints the available fonts to stderr.
func ListFonts() {
	fmt.Fprintf(os.Stderr, "\nAvailable fonts for -font option:\n\n")
	fmt.Fprintf(os.Stderr, "  %-12s %-20s %-10s %s\n", "NAME", "DISPLAY NAME", "LICENSE", "DESCRIPTION")
	fmt.Fprintf(os.Stderr, "  %-12s %-20s %-10s %s\n",
		"----", "------------", "-------", "-----------")
	for _, f := range BuiltinFonts {
		installed := ""
		if f.Name == "ipaex" {
			installed = " [TeX Live]"
		} else if IsFontInstalled(f.Name) {
			installed = " [installed]"
		}
		fmt.Fprintf(os.Stderr, "  %-12s %-20s %-10s %s%s\n",
			f.Name, f.DisplayName, f.License, f.Description, installed)
	}

	// Check for custom fonts in config
	customFonts := LoadCustomFonts()
	if len(customFonts) > 0 {
		fmt.Fprintf(os.Stderr, "\nCustom fonts (from %s):\n", getConfigPath())
		for _, cf := range customFonts {
			installed := ""
			if IsFontInstalled(cf.Name) {
				installed = " [installed]"
			}
			fmt.Fprintf(os.Stderr, "  %-12s %-20s %-10s %s%s\n",
				cf.Name, cf.DisplayName, cf.License, cf.Description, installed)
		}
	}

	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  src2latex -font hackgen sample.rb     # Use HackGen font\n")
	fmt.Fprintf(os.Stderr, "  src2latex -font install hackgen        # Download and install HackGen\n")
	fmt.Fprintf(os.Stderr, "  src2latex -font install all            # Download all fonts\n")
	fmt.Fprintf(os.Stderr, "\n")
}

// GetCommentFontDef returns the CommentFontDef for the given name, or nil.
func GetCommentFontDef(name string) *CommentFontDef {
	for i := range BuiltinCommentFonts {
		if BuiltinCommentFonts[i].Name == name {
			return &BuiltinCommentFonts[i]
		}
	}
	return nil
}

// IsCommentFontInstalled checks if a comment font is available.
func IsCommentFontInstalled(name string) bool {
	cfd := GetCommentFontDef(name)
	if cfd == nil {
		return false
	}
	if cfd.TexLive {
		return true
	}
	fontPath := filepath.Join(GetFontDir(), "comment-"+name, cfd.RegularFile)
	_, err := os.Stat(fontPath)
	return err == nil
}

// GetCommentFontPath returns the directory path for a downloaded comment font.
func GetCommentFontPath(name string) string {
	return filepath.Join(GetFontDir(), "comment-"+name) + "/"
}

// ListCommentFonts prints the available comment fonts to stderr.
func ListCommentFonts() {
	fmt.Fprintf(os.Stderr, "\nAvailable fonts for -commentfont option (明朝体・セリフ体):\n\n")
	fmt.Fprintf(os.Stderr, "  %-12s %-22s %-10s %s\n", "NAME", "DISPLAY NAME", "LICENSE", "DESCRIPTION")
	fmt.Fprintf(os.Stderr, "  %-12s %-22s %-10s %s\n",
		"----", "------------", "-------", "-----------")
	for _, f := range BuiltinCommentFonts {
		installed := ""
		if f.TexLive {
			installed = " [TeX Live]"
		} else if IsCommentFontInstalled(f.Name) {
			installed = " [installed]"
		}
		fmt.Fprintf(os.Stderr, "  %-12s %-22s %-10s %s%s\n",
			f.Name, f.DisplayName, f.License, f.Description, installed)
	}
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  src2latex -commentfont haranoaji sample.c           # TeX Live 原ノ味明朝\n")
	fmt.Fprintf(os.Stderr, "  src2latex -commentfont install noto-serif           # Noto Serif JPをインストール\n")
	fmt.Fprintf(os.Stderr, "  src2latex -commentfont install all                  # 全明朝体をインストール\n")
	fmt.Fprintf(os.Stderr, "  src2latex -font firple -commentfont haranoaji sample.c\n")
	fmt.Fprintf(os.Stderr, "\n")
}

// InstallCommentFont downloads and installs a comment font.
func InstallCommentFont(name string) error {
	if name == "all" {
		for _, f := range BuiltinCommentFonts {
			if f.TexLive {
				continue
			}
			fmt.Fprintf(os.Stderr, "Installing %s...\n", f.DisplayName)
			if err := InstallCommentFont(f.Name); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: failed to install %s: %v\n", f.Name, err)
			}
		}
		return nil
	}

	cfd := GetCommentFontDef(name)
	if cfd == nil {
		return fmt.Errorf("unknown comment font: %s (use -commentfont list to see available fonts)", name)
	}
	if cfd.TexLive {
		fmt.Fprintf(os.Stderr, "%s is bundled with TeX Live. No download needed.\n", cfd.DisplayName)
		return nil
	}
	if cfd.GithubRepo == "" {
		return fmt.Errorf("font %s has no download source configured", name)
	}

	if IsCommentFontInstalled(name) {
		fmt.Fprintf(os.Stderr, "%s is already installed.\n", cfd.DisplayName)
		return nil
	}

	// For noto-cjk, the release tag pattern is different (Serif2.003 instead of latest)
	// Use the releases/latest API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", cfd.GithubRepo)
	fmt.Fprintf(os.Stderr, "Fetching release info from %s...\n", cfd.GithubRepo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %v", err)
	}

	// Find matching asset
	var downloadURL, assetName string
	for _, a := range release.Assets {
		if strings.Contains(a.Name, cfd.ZipAsset) {
			downloadURL = a.BrowserDownloadURL
			assetName = a.Name
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no matching asset found for %s in release %s", cfd.ZipAsset, release.TagName)
	}

	fmt.Fprintf(os.Stderr, "Downloading %s (%s)...\n", assetName, release.TagName)

	// Download ZIP
	zipResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer zipResp.Body.Close()

	// Save to temp file
	tmpFile, err := os.CreateTemp("", "src2tex-commentfont-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	written, err := io.Copy(tmpFile, zipResp.Body)
	if err != nil {
		return fmt.Errorf("download interrupted: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Downloaded %.1f MB\n", float64(written)/1024/1024)

	// Extract font files
	destDir := filepath.Join(GetFontDir(), "comment-"+name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create font directory: %v", err)
	}

	// Use a simple FontDef adapter for extractFontFiles
	adapter := &FontDef{
		RegularFile: cfd.RegularFile,
		BoldFile:    cfd.BoldFile,
	}
	return extractFontFiles(tmpFile.Name(), destDir, adapter)
}

// IsFontInstalled checks if the font files exist in the font directory.
func IsFontInstalled(name string) bool {
	fd := GetFontDef(name)
	if fd == nil {
		// Check custom fonts
		for _, cf := range LoadCustomFonts() {
			if cf.Name == name {
				fd = &cf
				break
			}
		}
	}
	if fd == nil {
		return false
	}
	if fd.Name == "ipaex" {
		return true // Always available via TeX Live
	}
	fontPath := filepath.Join(GetFontDir(), fd.Name, fd.RegularFile)
	_, err := os.Stat(fontPath)
	return err == nil
}

// IsIPAexAvailable checks if IPAex is available via TeX Live.
func IsIPAexAvailable() bool {
	// Check common TeX Live paths
	paths := []string{
		"/usr/local/texlive/2026/texmf-dist/fonts/truetype/public/ipaex/ipaexg.ttf",
		"/usr/local/texlive/2025/texmf-dist/fonts/truetype/public/ipaex/ipaexg.ttf",
		"/usr/local/texlive/2024/texmf-dist/fonts/truetype/public/ipaex/ipaexg.ttf",
	}

	// Also try platform-specific paths
	if runtime.GOOS == "linux" {
		paths = append(paths,
			"/usr/share/texlive/texmf-dist/fonts/truetype/public/ipaex/ipaexg.ttf",
			"/usr/share/fonts/truetype/ipaex/ipaexg.ttf",
		)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// InstallFont downloads and installs a font.
func InstallFont(name string) error {
	if name == "all" {
		for _, f := range BuiltinFonts {
			if f.Name == "ipaex" {
				continue
			}
			fmt.Fprintf(os.Stderr, "Installing %s...\n", f.DisplayName)
			if err := InstallFont(f.Name); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: failed to install %s: %v\n", f.Name, err)
			}
		}
		return nil
	}

	fd := GetFontDef(name)
	if fd == nil {
		return fmt.Errorf("unknown font: %s (use -font list to see available fonts)", name)
	}
	if fd.Name == "ipaex" {
		fmt.Fprintf(os.Stderr, "IPAex is bundled with TeX Live. No download needed.\n")
		return nil
	}
	if fd.GithubRepo == "" {
		return fmt.Errorf("font %s has no download source configured", name)
	}

	if IsFontInstalled(name) {
		fmt.Fprintf(os.Stderr, "%s is already installed.\n", fd.DisplayName)
		return nil
	}

	// Get latest release info from GitHub
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", fd.GithubRepo)
	fmt.Fprintf(os.Stderr, "Fetching release info from %s...\n", fd.GithubRepo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %v", err)
	}

	// Find matching asset
	var downloadURL, assetName string
	for _, a := range release.Assets {
		if strings.Contains(a.Name, fd.ZipAsset) {
			downloadURL = a.BrowserDownloadURL
			assetName = a.Name
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no matching asset found for %s in release %s", fd.ZipAsset, release.TagName)
	}

	fmt.Fprintf(os.Stderr, "Downloading %s (%s)...\n", assetName, release.TagName)

	// Download ZIP
	zipResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer zipResp.Body.Close()

	// Save to temp file
	tmpFile, err := os.CreateTemp("", "src2tex-font-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	written, err := io.Copy(tmpFile, zipResp.Body)
	if err != nil {
		return fmt.Errorf("download interrupted: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Downloaded %.1f MB\n", float64(written)/1024/1024)

	// Extract font files
	destDir := filepath.Join(GetFontDir(), name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create font directory: %v", err)
	}

	return extractFontFiles(tmpFile.Name(), destDir, fd)
}

// extractFontFiles extracts .ttf/.otf files from a ZIP archive.
func extractFontFiles(zipPath, destDir string, fd *FontDef) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %v", err)
	}
	defer r.Close()

	extracted := 0
	for _, f := range r.File {
		// Only extract .ttf and .otf files
		lower := strings.ToLower(f.Name)
		if !strings.HasSuffix(lower, ".ttf") && !strings.HasSuffix(lower, ".otf") {
			continue
		}

		// Get just the filename (ignore directory structure in ZIP)
		baseName := filepath.Base(f.Name)

		// For the font to work, we need at least the Regular and Bold files
		// Extract all font files for flexibility
		destPath := filepath.Join(destDir, baseName)

		rc, err := f.Open()
		if err != nil {
			continue
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			continue
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			continue
		}
		extracted++
	}

	if extracted == 0 {
		return fmt.Errorf("no font files found in archive")
	}

	// Verify required files exist
	regPath := filepath.Join(destDir, fd.RegularFile)
	if _, err := os.Stat(regPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: expected file %s not found.\n", fd.RegularFile)
		fmt.Fprintf(os.Stderr, "Extracted %d font files to %s\n", extracted, destDir)
		fmt.Fprintf(os.Stderr, "Available files:\n")
		entries, _ := os.ReadDir(destDir)
		for _, e := range entries {
			fmt.Fprintf(os.Stderr, "  %s\n", e.Name())
		}
		return fmt.Errorf("required font file not found: %s", fd.RegularFile)
	}

	fmt.Fprintf(os.Stderr, "Installed %d font files to %s\n", extracted, destDir)
	return nil
}

// --- Custom Font Config ---

// CustomFontConfig represents the JSON config file structure.
type CustomFontConfig struct {
	Fonts []FontDef `json:"fonts"`
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".src2tex", "fonts.json")
}

// LoadCustomFonts loads custom font definitions from ~/.src2tex/fonts.json
func LoadCustomFonts() []FontDef {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}
	var config CustomFontConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil
	}
	return config.Fonts
}

// LookupFont finds a font definition by name (builtin or custom).
func LookupFont(name string) *FontDef {
	// Check builtin first
	if fd := GetFontDef(name); fd != nil {
		return fd
	}
	// Check custom fonts
	for _, cf := range LoadCustomFonts() {
		if cf.Name == name {
			return &cf
		}
	}
	return nil
}

// GetFontPath returns the absolute path to the font directory for the given font.
func GetFontPath(name string) string {
	if name == "ipaex" {
		return "" // IPAex uses TeX Live's built-in path (file name lookup)
	}
	return filepath.Join(GetFontDir(), name) + "/"
}

// resolveCommentFontXeLaTeX resolves a comment font name to a \setCJKmainfont line.
// It checks: 1) known comment font definitions, 2) TeX Live paths, 3) downloaded fonts, 4) system fonts.
func resolveCommentFontXeLaTeX(fontName string) string {
	// Check if it's a known comment font definition
	cfd := GetCommentFontDef(fontName)
	if cfd != nil {
		if cfd.TexLive {
			// Find in TeX Live paths
			for _, dir := range texLiveFontDirs() {
				path := filepath.Join(dir, cfd.RegularFile)
				if _, err := os.Stat(path); err == nil {
					if cfd.BoldFile != "" {
						return fmt.Sprintf("  \\setCJKmainfont[Path=%s/, Extension=%s, BoldFont=%s]{%s}\n",
							dir, cfd.Extension, cfd.BoldFile, strings.TrimSuffix(cfd.RegularFile, cfd.Extension))
					}
					return fmt.Sprintf("  \\setCJKmainfont[Path=%s/, Extension=%s]{%s}\n",
						dir, cfd.Extension, strings.TrimSuffix(cfd.RegularFile, cfd.Extension))
				}
			}
		} else if IsCommentFontInstalled(fontName) {
			// Downloaded comment font
			fontPath := GetCommentFontPath(fontName)
			if cfd.BoldFile != "" {
				return fmt.Sprintf("  \\setCJKmainfont[Path=%s, Extension=%s, BoldFont=%s]{%s}\n",
					fontPath, cfd.Extension, cfd.BoldFile, strings.TrimSuffix(cfd.RegularFile, cfd.Extension))
			}
			return fmt.Sprintf("  \\setCJKmainfont[Path=%s, Extension=%s]{%s}\n",
				fontPath, cfd.Extension, strings.TrimSuffix(cfd.RegularFile, cfd.Extension))
		}
	}

	// Try raw font name in TeX Live paths
	for _, dir := range texLiveFontDirs() {
		for _, ext := range []string{".otf", ".ttf"} {
			path := filepath.Join(dir, fontName+ext)
			if _, err := os.Stat(path); err == nil {
				return fmt.Sprintf("  \\setCJKmainfont[Path=%s/, Extension=%s]{%s}\n",
					dir, ext, fontName)
			}
		}
	}

	// Check as absolute/relative path with extension
	if _, err := os.Stat(fontName); err == nil {
		dir := filepath.Dir(fontName)
		base := filepath.Base(fontName)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		return fmt.Sprintf("  \\setCJKmainfont[Path=%s/, Extension=%s]{%s}\n",
			dir, ext, name)
	}

	// Fall back to bare font name (system font)
	return fmt.Sprintf("  \\setCJKmainfont{%s}\n", fontName)
}

// resolveCommentFontLuaLaTeX resolves a comment font name to a \setmainjfont line.
func resolveCommentFontLuaLaTeX(fontName string) string {
	// Check if it's a known comment font definition
	cfd := GetCommentFontDef(fontName)
	if cfd != nil {
		if cfd.TexLive {
			// LuaLaTeX can find TeX Live fonts by filename
			return fmt.Sprintf("  \\setmainjfont{%s}\n", cfd.RegularFile)
		} else if IsCommentFontInstalled(fontName) {
			fontPath := GetCommentFontPath(fontName)
			return fmt.Sprintf("  \\setmainjfont[Path=%s]{%s}\n", fontPath, cfd.RegularFile)
		}
	}

	// Try raw font name in TeX Live
	for _, dir := range texLiveFontDirs() {
		for _, ext := range []string{".otf", ".ttf"} {
			path := filepath.Join(dir, fontName+ext)
			if _, err := os.Stat(path); err == nil {
				return fmt.Sprintf("  \\setmainjfont{%s%s}\n", fontName, ext)
			}
		}
	}

	// Fall back to bare font name
	return fmt.Sprintf("  \\setmainjfont{%s}\n", fontName)
}

// texLiveFontDirs returns common TeX Live font directories to search.
func texLiveFontDirs() []string {
	return []string{
		"/usr/local/texlive/2026/texmf-dist/fonts/opentype/public/haranoaji",
		"/usr/local/texlive/2025/texmf-dist/fonts/opentype/public/haranoaji",
		"/usr/local/texlive/2024/texmf-dist/fonts/opentype/public/haranoaji",
		"/usr/local/texlive/2026/texmf-dist/fonts/truetype/public/ipaex",
		"/usr/local/texlive/2025/texmf-dist/fonts/truetype/public/ipaex",
		"/usr/local/texlive/2024/texmf-dist/fonts/truetype/public/ipaex",
	}
}

// GenerateFontPreambleXeLaTeX generates the XeLaTeX preamble for the given font.
func GenerateFontPreambleXeLaTeX(fd *FontDef) string {
	var sb strings.Builder

	sb.WriteString("  \\usepackage{fontspec}\n")
	sb.WriteString("  \\usepackage{xeCJK}\n")

	if fd.Unified {
		// Unified font: use the same font for both Latin and CJK
		fontPath := GetFontPath(fd.Name)
		if fontPath != "" {
			sb.WriteString(fmt.Sprintf("  \\setmonofont[Path=%s, BoldFont=%s]{%s}\n",
				fontPath, fd.BoldFile, fd.RegularFile))
			sb.WriteString(fmt.Sprintf("  \\setCJKmonofont[Path=%s, BoldFont=%s]{%s}\n",
				fontPath, fd.BoldFile, fd.RegularFile))
		} else {
			sb.WriteString(fmt.Sprintf("  \\setmonofont[BoldFont=%s]{%s}\n",
				fd.BoldFile, fd.RegularFile))
			sb.WriteString(fmt.Sprintf("  \\setCJKmonofont[BoldFont=%s]{%s}\n",
				fd.BoldFile, fd.RegularFile))
		}
		// CJKmainfont: use comment font if specified, otherwise same as code font
		if CommentFontName != "" {
			sb.WriteString(resolveCommentFontXeLaTeX(CommentFontName))
		} else if fontPath != "" {
			sb.WriteString(fmt.Sprintf("  \\setCJKmainfont[Path=%s, BoldFont=%s]{%s}\n",
				fontPath, fd.BoldFile, fd.RegularFile))
		} else {
			sb.WriteString(fmt.Sprintf("  \\setCJKmainfont[BoldFont=%s]{%s}\n",
				fd.BoldFile, fd.RegularFile))
		}
	} else {
		// CJK-only font: keep Latin Modern TT for Latin, use CJK font with scale
		mainReg := fd.MainRegularFile
		if mainReg == "" {
			mainReg = fd.RegularFile
		}
		mainBold := fd.MainBoldFile
		if mainBold == "" {
			mainBold = fd.BoldFile
		}
		if CommentFontName != "" {
			sb.WriteString(resolveCommentFontXeLaTeX(CommentFontName))
		} else {
			sb.WriteString(fmt.Sprintf("  \\setCJKmainfont[BoldFont=%s, Scale=%.2f]{%s}\n",
				mainBold, fd.Scale, mainReg))
		}
		sb.WriteString(fmt.Sprintf("  \\setCJKmonofont[Scale=%.2f]{%s}\n",
			fd.Scale, fd.RegularFile))
	}

	// Suppress automatic spacing between CJK and Latin characters
	sb.WriteString("  \\XeTeXlinebreaklocale \"\"\n")
	sb.WriteString("  \\xeCJKsetup{CJKglue={},CJKecglue={}}\n")

	return sb.String()
}

// GenerateFontPreambleLuaLaTeX generates the LuaLaTeX preamble for the given font.
func GenerateFontPreambleLuaLaTeX(fd *FontDef) string {
	var sb strings.Builder

	if fd.Unified {
		fontPath := GetFontPath(fd.Name)
		sb.WriteString("  \\usepackage{luatexja-fontspec}\n")
		if fontPath != "" {
			sb.WriteString(fmt.Sprintf("  \\setmonofont[Path=%s, BoldFont=%s]{%s}\n",
				fontPath, fd.BoldFile, fd.RegularFile))
		} else {
			sb.WriteString(fmt.Sprintf("  \\setmonofont[BoldFont=%s]{%s}\n",
				fd.BoldFile, fd.RegularFile))
		}
		// mainjfont: use comment font if specified, otherwise same as code font
		if CommentFontName != "" {
			sb.WriteString(resolveCommentFontLuaLaTeX(CommentFontName))
		} else if fontPath != "" {
			sb.WriteString(fmt.Sprintf("  \\setmainjfont[Path=%s, BoldFont=%s]{%s}\n",
				fontPath, fd.BoldFile, fd.RegularFile))
		} else {
			sb.WriteString(fmt.Sprintf("  \\setmainjfont[BoldFont=%s]{%s}\n",
				fd.BoldFile, fd.RegularFile))
		}
	} else {
		if CommentFontName != "" {
			sb.WriteString("  \\usepackage{luatexja-fontspec}\n")
			sb.WriteString(resolveCommentFontLuaLaTeX(CommentFontName))
		} else {
			sb.WriteString("  \\usepackage[ipaex]{luatexja-preset}\n")
		}
	}

	// Suppress automatic spacing
	sb.WriteString("  \\ltjsetparameter{kanjiskip=0pt, xkanjiskip=0pt}\n")

	return sb.String()
}
