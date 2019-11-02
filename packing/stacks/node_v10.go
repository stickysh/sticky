package stacks

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stickysh/sticky/pkg/fs"
	"log"
	"text/template"
)

type nodeStack struct {
	name               string
	runnerPath         string
	tmplPath           string
	filenames          []string
	entryPointFile     string
	entryPointFunction string
}



func NewNodev10() *nodeStack {
	name := "node10"
	return &nodeStack{
		name:               name,
		runnerPath:         fmt.Sprintf("/%s/%s", name, "runner.js"),
		tmplPath:           fmt.Sprintf("/%s/%s", name, "index.js.tmpl"),
		filenames:          []string{"__func__.js", "config.json", "params.json"},
		entryPointFile:     "__func__.js",
		entryPointFunction: "action",
	}
}

func (ns *nodeStack) Build(wd string) (map[string][]byte, error) {
	files := make(map[string][]byte,0)
	// Read action
	for i := 0; i < len(ns.filenames); i++ {
		fileName := ns.filenames[i]
		content := fs.ReadFile(fmt.Sprintf("%s/%s", wd, fileName))
		if content == nil {
			continue
		}

		files[fileName] = content
	}

	return files, nil
}

func (ns *nodeStack) Deploy(path string) error {

	// Write Runner
	data, err := fs.ReadStaticFile(ns.runnerPath)
	if err != nil {
		log.Fatal(err)
	}

	err = fs.WriteFile(fmt.Sprintf("%s/runner.js", path), data)
	if err != nil {
		return err
	}

	// Write Template
	buf := bytes.Buffer{}

	w := bufio.NewWriter(&buf)

	tmpl, _ := ns.getTemplate()

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: ns.entryPointFile,
		HandleName:  ns.entryPointFunction,
	})

	w.Flush()

	err = fs.WriteFile(fmt.Sprintf("%s/index.js", path), buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}

func (ns *nodeStack) getTemplate() (*template.Template, error) {

	buf, err := fs.ReadStaticFile(ns.tmplPath)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("index").Parse(string(buf))
	if err != nil {
		return nil, fmt.Errorf("template for stack %s was not found, %s", ns.tmplPath, err)
	}

	return tmpl, nil
}
