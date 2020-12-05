package generate

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func NewGeneratedFile(filePath string) *GeneratedFile {
	res := &GeneratedFile{
		filePath: filePath,
	}
	res.WriteDoNotEdit()
	return res
}

// A GeneratedFile is a generated file.
type GeneratedFile struct {
	buf      bytes.Buffer
	filePath string
}

// P prints a line to the generated output. It converts each parameter to a
// string following the same rules as fmt.Print. It never inserts spaces
// between parameters.
func (g *GeneratedFile) P(v ...interface{}) {
	var err error
	for _, x := range v {
		_, err = fmt.Fprint(&g.buf, x)
		if err != nil {
			logrus.Fatalf("GeneratedFile P %v have an err: %v", x, err)
		}
	}
	_, err = fmt.Fprintln(&g.buf)
	if err != nil {
		logrus.Fatalf("GeneratedFile P new line have an err: %v", err)
	}
}

func (g *GeneratedFile) WriteFile() (int, error) {
	dir := filepath.Dir(g.filePath)

	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return -1, err
		}
	}

	f, err := os.Create(g.filePath)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	n3, err := f.Write(g.buf.Bytes())
	if err != nil {
		return -1, err
	}

	if err := f.Sync(); err != nil {
		return -1, err
	}

	return n3, nil
}

func (g *GeneratedFile) WriteDoNotEdit() {
	g.P("// Code generated by go-gen. DO NOT EDIT.")
}
