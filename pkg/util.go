package util

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type Writer interface {
	Init(header interface{}, w *io.Writer)
	Write(w *io.Writer, toPrint ...interface{})
}

type TabWriter struct{}

func (tabWriter TabWriter) Init(header string, w *tabwriter.Writer) {

	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, header)

}

func (tabWriter TabWriter) Write(w *tabwriter.Writer, toPrint ...string) {
	fmt.Fprintln(w, strings.Join(toPrint, "\t"))
}
