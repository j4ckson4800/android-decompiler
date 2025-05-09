package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func bindInclude(tFs *template.Template) any {
	return func(name string, data any) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tFs.ExecuteTemplate(buf, name, data); err != nil {
			return buf.String(), fmt.Errorf("execute template: %w", err)
		}
		return buf.String(), nil
	}
}

// pasted from https://github.com/Masterminds/sprig/blob/master/strings.go#L109-L112
func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}
