# go-ghostscript

[![GoDoc](https://godoc.org/github.com/mrsaints/go-ghostscript/ghostscript?status.svg)](https://godoc.org/github.com/mrsaints/go-ghostscript/ghostscript)

Simple, and _idiomatic_ Go bindings for [Ghostscript][ghostscript] [Interpreter C API][api].

_Idiomatic_ is italicised because no true Go code should include [cgo][]. Ironic, I know.

> Ghostscript is a suite of software based on an interpreter for Adobe Systems' PostScript and Portable Document Format (PDF) page description languages. Its main purposes are the rasterization or rendering of such page description language files, for the display or printing of document pages, and the conversion between PostScript and PDF files. - [Wikipedia][wiki]

I would not recommend using this on production. I only worked on it to experiment with [cgo][], and I do not plan on maintaining it very often. Contributions are nevertheless, welcomed.


## Dependencies

To build, and run the package, you must have `libgs-dev` installed.

On Debian systems, this can be achieved using
`apt-get install libgs-dev`.


## Usage

1. Download, and install `go-ghostscript/ghostscript`:

    ```shell
    go get github.com/MrSaints/go-ghostscript/ghostscript
    ```

2. Import the package into your code:

    ```go
    import "github.com/MrSaints/go-ghostscript/ghostscript"
    ```

View the [GoDoc][], [examples][] or [code][] for more information.


[ghostscript]: http://ghostscript.com/
[api]: http://www.ghostscript.com/doc/current/API.htm
[wiki]: https://en.wikipedia.org/wiki/Ghostscript
[cgo]: https://golang.org/cmd/cgo/
[GoDoc]: https://godoc.org/github.com/mrsaints/go-ghostscript/ghostscript
[examples]: examples/
[code]: ghostscript/