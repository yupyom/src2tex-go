# src2tex-go

`src2tex-go` は、各種プログラミング言語のソースコードを読みやすくフォーマットされた TeX / LaTeX 形式のファイルに変換するユーティリティです。長年利用されてきた C言語版ユーティリティ `src2tex` (version 2.12) を、できる限りの互換性を維持したまま Go言語で移植・拡張したものです。

## 特徴

- **多言語対応**: C, C++, Java, Pascal, Lisp, Scheme, BASIC, Fortran, Perl, Tcl/Tk, MATLAB など、多種多様な言語のソースコードやコメントを自動判別し、適切に TeX / LaTeX 形式へ整形します。
- **モダン言語の追加サポート**: Go版では Python, Ruby, Rust, Go, JavaScript, TypeScript, Kotlin, Swift にも新たに対応しています。
- **TeX/LaTeX両対応**: TeX（通常）モードとLaTeXモードの両方をサポートしています。
- **Unicodeモード（デフォルト）**: デフォルトで XeLaTeX / LuaLaTeX / Tectonic に対応した Unicode 出力を生成します。日本語 CJK フォントも自動的に設定されます。生成された `.tex` ファイルは `iftex` パッケージによるエンジン自動検出を行うため、同一ファイルがどのエンジンでもそのままコンパイルできます。
- **CJKフォント管理**: `-font` オプションにより、IPAex Gothic（デフォルト）のほか、HackGen, UDEV Gothic, Firple などの等幅日本語フォントを選択可能です。GitHub からのフォントダウンロード・インストール機能も内蔵しています。
- **EPS図版の自動変換**: ソース内の `\special{epsfile=...}` を `\includegraphics` に自動変換し、Ghostscript による EPS→PDF 変換も行います。
- **互換性の維持**: C言語版 (`src2tex-212`) と同一の出力（日本語の `\mc` 等のpTeXマクロの制御やインデント計算を含む）を安定して生成するように最適化されています。

## コマンド名について

本プロジェクトは以下の名前体系を使用します：

| 名前 | 用途 |
|---|---|
| `src2tex-go` | プロジェクト名（Go版 src2tex） |
| `src2latex` | 実行バイナリ名（LaTeXモードで自動動作） |
| `src2tex` | バイナリをこの名前にリネーム/リンクすると plain TeX モードで動作 |

> **注意**: 旧名称 `src2latexg` も後方互換のため引き続き認識されます。

## 前提環境

| ツール | 用途 | 必須/任意 |
|---|---|---|
| [Go](https://go.dev/dl/) 1.21+ | ビルド | 必須 |
| [Task (go-task)](https://taskfile.dev/installation/) | タスクランナー | 任意（手動ビルドも可） |
| [XeLaTeX](https://tug.org/xetex/)、[LuaLaTeX](https://www.luatex.org/)、または [Tectonic](https://tectonic-typesetting.github.io/) | PDF 生成 | PDF を生成する場合にいずれか1つが必要 |
| [Ghostscript](https://www.ghostscript.com/) (`gs`) | EPS→PDF 変換 | 図版入りサンプルで必要 |

依存ライブラリはGoの標準ライブラリのみです（外部依存としてエンコーディング変換用の `golang.org/x/text` を利用）。

## コンパイル方法

### go-task を使う方法（推奨）

[go-task](https://taskfile.dev/) をインストール済みであれば、`Taskfile.yml` に定義されたタスクでビルドからPDF生成まで一貫して実行できます。go-taskはクロスプラットフォーム対応のタスクランナーで、Windows / macOS / Linux で同じコマンドが使えます。

```bash
# 現在のプラットフォーム向けにビルド
task build

# 全プラットフォーム向けにクロスコンパイル
task build:all

# 個別のプラットフォーム向けビルド
task build:darwin-arm64    # macOS (Apple Silicon)
task build:darwin-amd64    # macOS (Intel)
task build:linux-amd64     # Linux (amd64)
task build:windows-amd64   # Windows (amd64)

# サンプルの変換（TeX生成）
task samples

# サンプルの PDF 生成（変換 + PDF コンパイル）
task pdf

# Tectonic を使って PDF 生成
task pdf:compile TEX_ENGINE=tectonic

# クリーンアップ
task clean          # バイナリのみ削除
task clean:samples  # 生成された TeX/PDF を削除
task clean:all      # すべて削除
```

### 手動でビルドする方法

Go 言語がインストールされていれば、以下の標準コマンドで簡単にビルドできます。

```bash
go build -o src2latex
```

### クロスコンパイル（手動）

各種プラットフォーム向けのクロスコンパイルにも対応しています。
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

## 使用方法

### 基本的な使い方

デフォルトのバイナリ名は `src2latex` です。そのまま実行すると **LaTeX モード** かつ **Unicode モード** で動作します。

```bash
./src2latex <入力ファイル>
```

変換されたファイルは、元のファイル名に `.tex` が追加された名前で保存されます。
- 例: `hanoi.c` → `hanoi.c.tex`

### モードの切り替え

出力モードは実行ファイル名によって自動的に決まります。

| 実行ファイル名 | デフォルトモード |
|---|---|
| `src2latex` | **LaTeX モード**（デフォルトのビルド名） |
| `src2tex` | **plain TeX モード**（シンボリックリンク等で使用） |

```bash
# plain TeX モードで使いたい場合
ln -s src2latex src2tex
./src2tex <入力ファイル>
```

`-latex` / `-tex` オプションを使えば、実行ファイル名に関係なくモードを明示的に切り替えることもできます。

### オプション一覧

| オプション | 説明 |
|---|---|
| `-latex` | LaTeX モードを明示的に指定 |
| `-tex` | plain TeX モードを明示的に指定 |
| `-unicode` | XeLaTeX / LuaLaTeX / Tectonic 向けの Unicode 出力モード（**デフォルト**） |
| `-legacy` | 従来の pTeX / pLaTeX 向け出力モード（Unicode モードを無効化） |
| `-euc` | 入力ファイルのエンコーディングを EUC-JP として指定します |
| `-sjis` | 入力ファイルのエンコーディングを Shift_JIS として指定します |
| `-utf8` | 入力ファイルのエンコーディングを UTF-8 として指定します |
| `-font <name>` | CJK フォントを選択します（デフォルト: `ipaex`） |
| `-commentfont <name>` | コメント部分の CJK フォントを選択します（明朝体等） |
| `-fontdir <path>` | フォントのインストール先ディレクトリを指定します（デフォルト: `~/.src2tex/fonts/`） |
| `-proxy <url>` | フォントダウンロード時の HTTP プロキシを指定します |
| `-<n>` | 1ページあたりの行数を `n` に制限します（行番号付き出力）|
| `-0` | 行番号付き出力（行数制限なし） |
| `-header <file>` | カスタムヘッダーのプリアンブルファイルを指定（fancyhdr ヘッダー部を置換）|
| `-footer <file>` | カスタムフッターのプリアンブルファイルを指定（fancyhdr フッター部を置換）|
| `-v` | バージョン情報を表示します |

### 使用例

```bash
# 基本的な変換（Unicode・LaTeXモード、デフォルト）
./src2latex samples/hanoi.c

# EUC-JPのソースを変換
./src2latex -euc samples/newton.c

# HackGen フォントで変換
./src2latex -font hackgen samples/hanoi.rb

# コメントの日本語を明朝体で変換
./src2latex -font firple -commentfont haranoaji samples/hanoi.c

# 行番号付き出力
./src2latex -0 samples/hanoi.c

# カスタムヘッダー/フッターを使用
./src2latex -header my_header.tex -footer my_footer.tex samples/hanoi.c

# 従来の pLaTeX 向け出力
./src2latex -legacy samples/hanoi.c
```

### PDF の生成

変換されたTeXファイルは、お手持ちのTeX環境でコンパイルしてPDFを生成できます。デフォルトの Unicode モードで生成された `.tex` ファイルは、以下の3つのエンジンでの動作を確認しています。

```bash
# XeLaTeX を使う場合（TeX Live / MacTeX に標準搭載）
xelatex hanoi.c.tex

# LuaLaTeX を使う場合（TeX Live / MacTeX に標準搭載）
lualatex hanoi.c.tex

# Tectonic を使う場合（軽量でインストールが簡単）
tectonic hanoi.c.tex
```

オリジナル版の`src2tex`は、1990年代の古いバージョンの $\TeX$ や $\LaTeX$ を前提としていましたが、このバージョンでは新しい環境に即したプリアンブルを生成するように調整しています。`-font` や `-commentfont` を指定することで、コードやコメント部に好きな日本語書体を割り当てることができます。

ただし、コード部に適用される`-font`については、全角が半角の2倍幅になっている設計の書体を利用することが前提となっています。このモードを使うと、空白部分をboxで生成せず、リテラルのスペース（またはタブ）で埋める形を取ります。

> **注意**: 生成される `.tex` ファイルは `iftex` パッケージを使ってエンジンを自動検出するため、同一の `.tex` ファイルを上記いずれのエンジンでもそのままコンパイルできます。
> - **XeLaTeX** の場合: `fontspec` + `xeCJK` + `zxjafont` が使用されます
> - **LuaLaTeX** の場合: `luatexja-preset` が使用されます
> - **Tectonic** の場合: 内部で XeTeX エンジンが使用されるため、XeLaTeX と同じパッケージが適用されます

```bash
# 従来の pLaTeX を使う場合（-legacy で変換したファイル向け）
platex hanoi.c.tex
dvipdfmx hanoi.c.dvi
```

## CJK フォント管理

`-font` オプションで、コメント部分のフォント（CJKフォント）を切り替えることができます。

### 内蔵フォント一覧

| 名前 | フォント名 | タイプ | ライセンス | 説明 |
|---|---|---|---|---|
| `ipaex` | IPAex Gothic | 非統合 | IPA | TeX Live 同梱。ダウンロード不要（**デフォルト**） |
| `hackgen` | HackGen | 統合 | SIL OFL | Hack + 源柔ゴシック。プログラミング向けの人気フォント |
| `udev` | UDEV Gothic | 統合 | SIL OFL | JetBrains Mono + BIZ UDゴシック。UD対応 |
| `firple` | Firple | 統合 | SIL OFL | Fira Code + IBM Plex Sans JP。リガチャ対応。半角:全角 = 1:2 |

**統合フォント**（HackGen, UDEV Gothic, Firple）は、Latin文字と日本語文字の幅が統一されているため、PDF上でのテキスト選択・コピーが正確に動作します。

### フォントのインストールと使用

```bash
# 利用可能なフォントの一覧表示
./src2latex -font list

# フォントのダウンロード・インストール
./src2latex -font install hackgen       # HackGen をインストール
./src2latex -font install all           # 全フォントをインストール

# インストールしたフォントで変換
./src2latex -font hackgen samples/hanoi.rb

# フォントディレクトリのカスタマイズ
./src2latex -fontdir /path/to/fonts -font hackgen samples/hanoi.rb
```

フォントは `~/.src2tex/fonts/` にインストールされます。`-fontdir` オプションで変更可能です。

### カスタムフォントの追加

`~/.src2tex/fonts.json` にフォント定義を追加することで、独自のフォントも使用できます：

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
      "Description": "カスタムフォントの説明"
    }
  ]
}
```

## コメント用 CJK フォント管理

`-commentfont` オプションで、コメント部分の日本語フォントをコード部分と別に指定できます。統合フォント（HackGen, Firple 等）ではコード部分がゴシック体ですが、コメント内の日本語を明朝体にすることで、Computer Modern と調和した上品な組版が得られます。

### コメント用フォント一覧

| 名前 | フォント名 | ライセンス | 説明 |
|---|---|---|---|
| `haranoaji` | 原ノ味明朝 (Harano Aji Mincho) | SIL OFL | TeX Live 同梱。ダウンロード不要 |
| `ipaexm` | IPAex 明朝 | IPA | TeX Live 同梱。ダウンロード不要 |
| `noto-serif` | Noto Serif JP | SIL OFL | Google Noto 明朝体。要ダウンロード |

### コメントフォントのインストールと使用

```bash
# 利用可能なコメントフォントの一覧表示
./src2latex -commentfont list

# コメントフォントのダウンロード・インストール
./src2latex -commentfont install noto-serif    # Noto Serif JP をインストール
./src2latex -commentfont install all           # 全コメントフォントをインストール

# TeX Live 同梱のフォントを使用（ダウンロード不要）
./src2latex -font firple -commentfont haranoaji samples/hanoi.c
./src2latex -font hackgen -commentfont ipaexm samples/hanoi.rb

# ダウンロードしたフォントを使用
./src2latex -font firple -commentfont noto-serif samples/hanoi.c
```

> **仕組み**: `-commentfont` で指定されたフォントは LaTeX の `\setCJKmainfont`（XeLaTeX）/ `\setmainjfont`（LuaLaTeX）に設定されます。コメント内で `\rm`（ローマン体）モードに切り替わると、日本語テキストにこの明朝体が適用されます。コード部分は従来通り `\setCJKmonofont` のゴシック体が使われます。

### プロキシ経由でのフォントダウンロード

プロキシサーバ経由でインターネットに接続している環境では、`-proxy` オプションまたは環境変数を使用できます：

```bash
# -proxy オプションで指定
./src2latex -proxy http://proxy.example.com:8080 -font install hackgen
./src2latex -proxy http://proxy.example.com:8080 -commentfont install noto-serif

# 認証付きプロキシ
./src2latex -proxy http://user:pass@proxy.example.com:8080 -font install hackgen

# 環境変数で指定（Go の http パッケージが自動認識）
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
./src2latex -font install all
./src2latex -commentfont install all
```

## 対応言語一覧

ファイル拡張子に基づいて言語を自動判別します。判別できない場合は、ソース内のキーワードから推定を試みます。

### 元の src2tex-212 と同じ対応言語

| 拡張子 | 言語 |
|---|---|
| `.tex`, `.txt`, `.doc` | TEXT（テキスト） |
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

### Go版で追加された言語

| 拡張子 | 言語 |
|---|---|
| `.py` | Python |
| `.rb` | Ruby |
| `.rs` | Rust |
| `.go` | Go |
| `.js` | JavaScript |
| `.ts` | TypeScript |
| `.kt` | Kotlin |
| `.swift` | Swift |

## FAQ（よくある質問）

### Q1: src2tex形式でコメントを書くにはどうすればよいですか？

ソースコードのコメント領域内に、直接TeXの数式やコマンドを書くことができます。例えば、C言語の場合：

```c
r = sqrt(x*x + y*y);   /* 半径 $r=\sqrt{x^2+y^2}$ */
```

REDUCEの場合：
```reduce
int(x/sqrt(1-x^2), x);  % 積分 $\int{x\over\sqrt{1-x^2}}\,dx$
```

### Q2: EPS/PS画像をソースに埋め込むにはどうすればよいですか？

コメント領域（TeXモード部分）内で `\special{epsfile=...}` コマンドを使用します：

```c
/* See the following numerical simulation.
                    {\special{epsfile=simulation.eps}} */
```

Unicode モードでは、`\special{epsfile=...}` は自動的に `\includegraphics` に変換され、Ghostscriptを使ったEPS→PDF変換も自動的に行われます。

### Q3: コメント領域の文字デザインを変更できますか？

`\src2tex{...}` エスケープシーケンスを使用して、コメントエリアのフォントを変更できます。

```c
/* {\src2tex{texfont=tt}} */     /* 以降のコメントを typewriter font に */
/* {\src2tex{texfont=rm}} */     /* 以降のコメントを roman font に戻す */
/* {\src2tex{texfont=bf}} */     /* 以降のコメントを bold font に */
/* {\src2tex{texfont=it}} */     /* 以降のコメントを italic に */
/* {\src2tex{texfont=sl}} */     /* 以降のコメントを slant font に */
/* {\src2tex{texfont=sc}} */     /* 以降のコメントを small caps に */
```

プログラムエリアのフォントも変更可能です：
```c
/* {\src2tex{textfont=bf}} */    /* プログラムエリアを bold font に */
```

### Q4: タブ幅やインデント幅を変更できますか？

```c
/* {\src2tex{htab=4}} */         /* 水平タブサイズを4に設定 */
/* {\src2tex{vtab=2}} */         /* 垂直タブサイズを2に設定 */
```

### Q5: LaTeXのスタイルファイルを変更できますか？

ソースファイルの先頭行のコメントに、使用したいdocumentstyleを記述できます。
例えば、Pascalの場合：
```pascal
(* {\documentstyle[twocolumn,12pt]{article}} *)
```

ただし Unicode モード（デフォルト）では、XeLaTeX互換のプリアンブルが自動的に使用されるため、ソース内のdocumentstyle指定は無視されます。`-legacy` モードで変換した場合に有効です。

### Q6: コメントの位置揃えがうまくいかない場合は？

コメント内で揃えたい場合は、`\hfill` や `\hrulefill` などのTeXコマンドを使うと良いでしょう：

```c
/* {\hrulefill} */
/* {\hfill SUBROUTINE1 \hfill} */
/* {\hrulefill} */
```

あるいは、統合フォント（HackGen, UDEV Gothic, Firple）を使用すると、全角・半角文字の幅が統一されるため、位置揃えの問題が軽減されます：

```bash
./src2latex -font hackgen samples/hanoi.rb
```

### Q7: 行番号を出力できますか？

はい、`-0` オプションまたは `-<n>` オプションを使用します：

```bash
./src2latex -0 samples/hanoi.c       # 行番号付き（行数制限なし）
./src2latex -35 samples/hanoi.c      # 1ページ35行で行番号付き
```

### Q8: ヘッダーやフッターをカスタマイズできますか？

はい、`-header` / `-footer` オプションを使用して、プリアンブルの fancyhdr 設定部分をカスタムファイルの内容で置き換えることができます。

まず、カスタムヘッダーファイル（例: `my_header.tex`）を作成します：

```tex
\renewcommand{\headrulewidth}{0.4pt}
\fancyhf{}
\fancyhead[L]{\rm My Project Name}
\fancyhead[R]{\rm \today}
```

カスタムフッターファイル（例: `my_footer.tex`）も同様に作成します：

```tex
\fancyfoot[C]{\thepage}
\fancyfoot[L]{\rm Confidential}
\fancyfoot[R]{\rm Draft}
```

変換時にこれらを指定します：

```bash
# ヘッダーのみカスタマイズ
./src2latex -header my_header.tex samples/hanoi.c

# フッターのみカスタマイズ
./src2latex -footer my_footer.tex samples/hanoi.c

# 両方カスタマイズ
./src2latex -header my_header.tex -footer my_footer.tex samples/hanoi.c
```

> **注意**: `\usepackage{fancyhdr}` と `\pagestyle{fancy}` は自動的に出力されるため、カスタムファイルにはそれ以降の設定（`\fancyhead`, `\fancyfoot`, `\renewcommand{\headrulewidth}` 等）のみを記述してください。

> **ヒント**: ヘッダーやフッターを完全に非表示にしたい場合は、空の内容（`\fancyhf{}` のみ）のファイルを指定します。

### Q9: 特定のキーワードを太字にするには？

src2tex自体にはその機能はありませんが、変換後のTeXファイルに対して `sed` や他のテキスト処理ツールで後処理できます：

```bash
./src2latex sample.c
sed -e 's/}main(){/}{\\bf main()}{/g' sample.c.tex > sample_modified.tex
```

src2texはプログラムエリアのキーワードを `}keyword{` という形式で出力するため、この性質を利用してスタイルの変更が可能です。

## 注意事項

### エンコーディングについて

- 元の src2tex-212 のサンプル（`newton.c`, `simpson.c`, `farmer+hen.scm`）は EUC-JP エンコーディングです。これらを変換する場合は `-euc` を指定してください。
- Go版で新規作成されたサンプル（`hanoi.c`, `hanoi.go`, `hanoi.py` など）は UTF-8 です。`-utf8` を指定するか、デフォルトのまま使用してください。
- `-legacy` モードで変換した場合、出力は pTeX / pLaTeX 用の形式になります。

### Unicode モードについて

- デフォルトで Unicode モードが有効になっています。`-legacy` オプションで無効化できます。
- Unicode モードでは XeLaTeX / LuaLaTeX / Tectonic に対応したプリアンブルが自動生成されます。
- `iftex` パッケージによるエンジン自動検出を行い、XeLaTeX では `xeCJK` + `zxjafont`（またはカスタムフォント設定）、LuaLaTeX では `luatexja-preset` を自動的に使い分けます。これにより、同一の `.tex` ファイルがどのエンジンでもそのままコンパイルできます。
- 原始的な TeX コマンド（`\eqalign`, `\pmatrix` 等）の互換性マクロも自動的に含まれます。
- ソース内の `\special{epsfile=...}` は `\includegraphics` に自動変換されます。
- EPS ファイルは Ghostscript を使って自動的に PDF に変換されます（`gs` コマンドが必要）。

#### 対応 TeX エンジン

| エンジン | 説明 | インストール方法 |
|---|---|---|
| [XeLaTeX](https://tug.org/xetex/) | TeX Live / MacTeX に標準搭載。Unicode + OpenType フォントに対応 | `brew install --cask mactex` (macOS) |
| [LuaLaTeX](https://www.luatex.org/) | TeX Live / MacTeX に標準搭載。Lua 拡張による高度なフォント制御 | TeX Live に同梱 |
| [Tectonic](https://tectonic-typesetting.github.io/) | 軽量な XeTeX ベースのエンジン。パッケージの自動ダウンロード機能付き | `brew install tectonic` (macOS) |

### サンプルファイルについて

`samples/` フォルダには以下のサンプルが含まれています：

| ファイル | 言語 | 内容 |
|---|---|---|
| `newton.c` | C | Newton-Raphson法（数式・図入り）|
| `simpson.c` | C | Simpson積分公式（数式・図入り）|
| `hanoi.c` | C | ハノイの塔（英語/日本語コメント付き）|
| `hanoi.go` | Go | ハノイの塔 Go版 |
| `hanoi.pas` | Palcal | ハノイの塔 Pascal版 |
| `hanoi.pl` | Perl | ハノイの塔 Perl版 |
| `hanoi.py` | Python | ハノイの塔 Python版 |
| `hanoi.rb` | Ruby | ハノイの塔 Ruby版 |
| `hanoi.rs` | Rust | ハノイの塔 Rust版 |
| `hanoi.js` | JavaScript | ハノイの塔 JavaScript版 |
| `hanoi.ts` | TypeScript | ハノイの塔 TypeScript版 |
| `hanoi.kt` | Kotlin | ハノイの塔 Kotlin版 |
| `hanoi.sh` | Shell | ハノイの塔 ShellScript版 |
| `hanoi.swift` | Swift | ハノイの塔 Swift版 |
| `farmer+hen.scm` | Scheme | 「おじいさんと鶏」パズル |
| `popgen.red` | REDUCE | 集団遺伝学の偏微分方程式 |
| `sqrt_mat.red` | REDUCE | 正方行列の平方根 |

### サンプルの変換例

```bash
# UTF-8のサンプル（Go版で作成）— Unicode モード（デフォルト）
./src2latex samples/hanoi.c
./src2latex samples/hanoi.go
./src2latex samples/hanoi.py
./src2latex -font hackgen samples/hanoi.rb    # HackGen フォントで

# EUC-JPのサンプル（元のsrc2tex-212から）
./src2latex -euc samples/newton.c
./src2latex -euc samples/simpson.c
./src2latex -euc samples/farmer+hen.scm

# ASCII のサンプル
./src2latex samples/popgen.red
./src2latex samples/sqrt_mat.red

# PDFの生成（いずれかのエンジンを使用）
tectonic samples/hanoi.c.tex     # Tectonic
xelatex samples/hanoi.c.tex      # XeLaTeX
lualatex samples/hanoi.c.tex     # LuaLaTeX
```
## バージョンについて

go言語と数字の「5」の日本読みにあやかり、$ \sqrt 5 $ の小数点表現にしております。現在のバージョンは、version 2.23 です。

## オリジナルについて

オリジナルの `src2tex` version 2.12 は、天野一男（Kazuo AMANO）氏と野本慎一（Shinichi NOMOTO）氏によって開発されました。

オリジナルに付属する下記のサンプルファイル（および関連リソース）は、原作者に著作権があります。

- samples/farmer+hen.scm
- samples/hanoi.c
- samples/newton.c
- samples/popgen.red
- samples/simpson.c
- samples/sqrt_mat.red

> ** NOTE **: hanoi.c は、もともと hanoi89.c のような形で付属していましたが、C11でコンパイルできるように調整し、これをもとに他の言語版のhanoiを作成しています。


## ライセンス

オリジナルの src2tex-212 のライセンス条件に準じます。
