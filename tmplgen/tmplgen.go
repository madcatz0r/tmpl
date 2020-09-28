package main

import (
	"flag"
	"github.com/madcatz0r/tmpl/snake_case"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"
)

var (
	modelsDir       string
	templatesDir    string
	fileTemplate, _ = template.New("fileTemplate").Parse(`package {{ .Namespace }}

const (
{{range $key, $value := .FieldMap}}{{ $key }} = "{{ $value }}"
{{end}})
`)
)

type tmplFile struct {
	Namespace string
	FieldMap  map[string]string
}

func main() {
	modelsPtr := flag.String("model", "", "-model=./app/models")
	outputDir := flag.String("out", "", "-out=./app/common/vars")
	flag.Parse()
	if len(os.Args) != 3 {
		flag.PrintDefaults()
	}
	modelsDir = *modelsPtr
	templatesDir = *outputDir
	err := generateModelPackageFields()
	if err != nil {
		panic(err)
	}
}

func generateModelPackageFields() error {
	fSet := token.NewFileSet()
	dir, err := parser.ParseDir(fSet, modelsDir, func(info os.FileInfo) bool { return true }, 0)
	if err != nil {
		return err
	}
	for _, f := range dir["models"].Files {
		for key, declName := range f.Scope.Objects {
			if declName.Kind == ast.Typ {
				// create dir key
				path := filepath.Join(templatesDir, key)
				_ = os.Mkdir(path, os.ModePerm)
				tmpl := tmplFile{Namespace: key, FieldMap: make(map[string]string)}
				test, ok := declName.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
				if ok {
					for _, rez := range test.Fields.List {
						name := rez.Names[0].Name
						tmpl.FieldMap[name] = snake_case.ToSnakeCase(name)
					}
					// create file key.go
					fileOut, err := os.Create(filepath.Join(path, key+".go"))
					if err != nil {
						return err
					}
					// write file using template
					err = fileTemplate.Execute(fileOut, tmpl)
					if err != nil {
						// file writer close
						_ = fileOut.Close()
						return err
					}
					// file writer close
					err = fileOut.Close()
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
