package pkg

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func TabWriterInit(header string, tabWriter *tabwriter.Writer) error {

	tabWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, err := fmt.Fprintln(tabWriter, header)
	if err != nil {
		return err
	}
	return nil
}

func TabWriterWrite(toPrint []string, tabWriter *tabwriter.Writer) error {
	_, err := fmt.Fprintln(tabWriter, strings.Join(toPrint, "\t"))
	if err != nil {
		return err
	}
	return nil

}
