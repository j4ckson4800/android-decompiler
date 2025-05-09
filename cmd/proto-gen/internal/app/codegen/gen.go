package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/j4ckson4800/android-decompiler/cmd/proto-gen/internal/app/defs"
)

type codeGen struct {
	OutDir     string
	templateFs *template.Template
}

func NewCodegen(outDir string) (*codeGen, error) {
	tpl := template.New("")
	tpl.Funcs(
		template.FuncMap{
			"indent":  indent,
			"include": bindInclude(tpl),
		},
	)

	tpl, err := tpl.ParseFS(fs, "tpl/**.proto.tpl")
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	return &codeGen{
		OutDir:     outDir,
		templateFs: tpl,
	}, nil
}

func (c *codeGen) WritePackage(pkg *defs.ProtoPackage) error {
	packageDir := filepath.Join(c.OutDir, strings.TrimSuffix(pkg.GoPackageName, pkg.FileName))
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(packageDir + "/" + pkg.FileName + ".proto")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := c.templateFs.ExecuteTemplate(file, "package.proto.tpl", pkg); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}
