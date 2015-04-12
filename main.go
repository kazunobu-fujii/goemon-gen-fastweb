package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	rootPath     = "."
	assetsPath   = rootPath + "/assets"
	appServer    = "localhost:8080"
	goemonServer = "localhost:35730"
)

type Template struct {
	Path string
	Name string
	Body string
}

type Scaffold struct {
	Assets     string
	Server     string
	Livereload string
}

var templateFiles = []Template{
	Template{
		Path: rootPath,
		Name: "main.go",
		Body: `package main

import (
	"net/http"
)

func cacheHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
	})
}

func main() {
	h := http.FileServer(http.Dir("{{.Assets}}"))
	h = cacheHandler(h)
	http.Handle("/", h)
	http.ListenAndServe("{{.Server}}", nil)
}
`,
	},
	Template{
		Path: rootPath,
		Name: "goemon.yml",
		Body: `livereload: {{.Livereload}}
tasks:
- match: '{{.Assets}}/*.gcss'
  commands:
  - cat ${GOEMON_TARGET_FILE} | gcss > ${GOEMON_TARGET_DIR}/${GOEMON_TARGET_NAME}
  - :livereload /
- match: '{{.Assets}}/*.ace'
  commands:
  - cat ${GOEMON_TARGET_FILE} | ace > ${GOEMON_TARGET_DIR}/${GOEMON_TARGET_NAME}
  - echo "<script src=\"http://{{.Livereload}}/livereload.js\"></script>" >> ${GOEMON_TARGET_DIR}/${GOEMON_TARGET_NAME}
  - :livereload /
- match: '*.go'
  commands:
  - go build
  - :restart
  - :livereload /
`,
	},
	Template{
		Path: assetsPath,
		Name: "sample.html.ace",
		Body: `= doctype html
html lang=en
  head
    title Hello Ace
    link rel="stylesheet" type="text/css" href="sample.css"
  body
    h1 goemon + ace + gcss
    #container.wrapper
      p..
        Ace is an HTML template engine for Go.
        This engine simplifies HTML coding in Go web application development.
      p GCSS is a pure Go CSS preprocessor.
      p goemon is Go Extensible Monitoring
    = javascript
      console.log('Welcome to Ace');
`,
	},
	Template{
		Path: assetsPath,
		Name: "sample.css.gcss",
		Body: `$main-color: blue

h1
  color: $main-color
`,
	},
}

func makeFile(file string, body string, s *Scaffold) {
	tmpl, err := template.New(file).Parse(body)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, s)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(file, doc.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func generateFile() {
	scaffold := Scaffold{
		Assets:     assetsPath,
		Server:     appServer,
		Livereload: goemonServer,
	}

	os.Mkdir(assetsPath, 0755)
	for _, tpl := range templateFiles {
		file := tpl.Path + "/" + tpl.Name
		makeFile(file, tpl.Body, &scaffold)
		fmt.Fprintf(os.Stderr, "generated %s\n", file)
	}
}

func isEmptyDir(path string) bool {
	if fis, _ := ioutil.ReadDir(path); len(fis) > 0 {
		return false
	}
	return true
}

func checkEmpty() bool {
	pwd, _ := os.Getwd()
	return isEmptyDir(pwd)
}

func main() {
	if !checkEmpty() {
		fmt.Fprintln(os.Stderr, "file exists")
		os.Exit(1)
	}
	generateFile()

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "install goemon , Ace and GCSS")
	fmt.Fprintln(os.Stderr, "  go get github.com/mattn/goemon/cmd/goemon")
	fmt.Fprintln(os.Stderr, "  go get github.com/yosssi/ace/cmd/ace")
	fmt.Fprintln(os.Stderr, "  go get github.com/yosssi/gcss/cmd/gcss")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "run command")
	fmt.Fprintln(os.Stderr, "  goemon go run main.go")
	fmt.Fprintln(os.Stderr, "  touch assets/*")
	fmt.Fprintln(os.Stderr, "  open http://localhost:8080/sample.html")
	fmt.Fprintln(os.Stderr, "")
}
