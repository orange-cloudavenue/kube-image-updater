package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
)

func generateDocImageStatus() {
	tmplFuncs := template.FuncMap{
		"imageStatusLastSync": func() string {
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, "api/v1alpha1/image_models.go", nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			imgStatusSlice := [][]string{}

			for _, decl := range astFile.Decls {
				if _, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range decl.(*ast.GenDecl).Specs {
						if _, ok := spec.(*ast.ValueSpec); ok {
							for _, ident := range spec.(*ast.ValueSpec).Names {
								imgStatusSlice = append(imgStatusSlice, []string{ident.Name, ident.Obj.Decl.(*ast.ValueSpec).Doc.Text()})
							}
						}
					}
				}
			}

			// pretty print table
			prettyPrintedTable, err := markdown.
				NewTableFormatterBuilder().
				WithAlphabeticalSortIn(markdown.ASCENDING_ORDER).
				Build("Last-Sync", "Description").
				Format(imgStatusSlice)
			if err != nil {
				panic(err)
			}

			return prettyPrintedTable
		},
	}

	// os read file
	file, err := os.ReadFile("docs/advanced/image-status.md.tmpl")
	if err != nil {
		log.Default().Printf("Failed to open file: %v", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("image-status").Funcs(tmplFuncs).Parse(string(file)))

	// write template to file
	f, err := os.Create("docs/advanced/image-status.md")
	defer f.Close()
	if err != nil {
		log.Default().Printf("Failed to create file: %v", err)
		f.Close()
		os.Exit(1) //nolint: gocritic
	}
	if err := tmpl.Execute(f, nil); err != nil {
		log.Default().Printf("Failed to execute template: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
