package ctl

import (
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

func PrintAsJSON(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	os.Stdout.Write(b)
	return nil
}

func PrintAsTable(handler func(table.Writer)) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	handler(t)

	t.Render()
}
