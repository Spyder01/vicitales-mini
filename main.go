package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // file exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // definitely does not exist
	}

	return false
}

var md = goldmark.New()
var tpl = template.Must(template.ParseFiles("templates/story.html"))

func renderMarkdown(path string) (string, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := md.Convert(input, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Redirect(w, r, "/stories", http.StatusFound)
		return
	}

	// Map URL to markdown file
	mdPath := filepath.Join("content", path) + ".md"
	htmlContent, err := renderMarkdown(mdPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	baseName := filepath.Base(path)
	pageVar := map[string]any{
		"Content": template.HTML(htmlContent),
		"Title":   baseName,
		"Year":    time.Now().Year(),
	}

	dirPath := filepath.Dir(path)

	num, err := strconv.Atoi(baseName)
	if err != nil {
		tpl.Execute(w, pageVar)
		return
	}

	prevLink := filepath.Join(dirPath, strconv.Itoa(num-1))
	if fileExists(filepath.Join("content", prevLink+".md")) {
		pageVar["PrevLink"] = prevLink
	}

	nextLink := filepath.Join(dirPath, strconv.Itoa(num+1))
	if fileExists(filepath.Join("content", nextLink+".md")) {
		pageVar["NextLink"] = nextLink
	}

	tpl.Execute(w, pageVar)
}

func main() {
	http.HandleFunc("/", storyHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
