package pkg

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Writer interface {
	Init(header interface{})
	Write(toPrint ...interface{})
}

type TabWriter struct {
	Writer *tabwriter.Writer
}

func (w TabWriter) Init(header string) {

	w.Writer.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w.Writer, header)

}

func (w TabWriter) Write(toPrint ...string) {
	fmt.Fprintln(w.Writer, strings.Join(toPrint, "\t"))
}
