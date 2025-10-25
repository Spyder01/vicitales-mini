package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

var md = goldmark.New()
var tpl = template.Must(template.ParseFiles("templates/story.html"))

type Breadcrumb struct {
	Name string
	Link string
}

// RenderMarkdown converts a Markdown file to HTML.
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Recursively copy a folder
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

// Pre-render all markdown files into public/ preserving folder structure
func prerender() error {
	contentDir := "content"
	outputDir := "public"

	// Copy static folder first
	if err := copyDir("static", filepath.Join(outputDir, "static")); err != nil {
		return err
	}

	return filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Relative path like "fantasy/red-lily/1.md"
		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}

		// Output path like "public/fantasy/red-lily/1.html"
		outPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".md")+".html")
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		html, err := renderMarkdown(path)
		if err != nil {
			return err
		}

		pageVars := map[string]any{
			"Content": template.HTML(html),
			"Title":   info.Name(),
			"Year":    time.Now().Year(),
		}

		// Generate breadcrumbs
		parts := strings.Split(relPath, string(os.PathSeparator))
		var breadcrumbs []Breadcrumb
		for i := 0; i < len(parts)-1; i++ {
			dirPath := filepath.Join(parts[:i+1]...)
			breadcrumbs = append(breadcrumbs, Breadcrumb{
				Name: parts[i],
				Link: filepath.ToSlash(filepath.Join(dirPath, "1.html")),
			})
		}
		pageVars["Breadcrumbs"] = breadcrumbs

		// Next / Prev links (numeric chapters)
		chapterNum, err := strconv.Atoi(strings.TrimSuffix(info.Name(), ".md"))
		if err == nil {
			storyDir := filepath.Dir(path)
			prev := filepath.Join(storyDir, fmt.Sprintf("%d.md", chapterNum-1))
			next := filepath.Join(storyDir, fmt.Sprintf("%d.md", chapterNum+1))

			if fileExists(prev) {
				relPrev, _ := filepath.Rel(contentDir, prev)
				pageVars["PrevLink"] = strings.TrimSuffix(filepath.ToSlash(relPrev), ".md") + ".html"
			}
			if fileExists(next) {
				relNext, _ := filepath.Rel(contentDir, next)
				pageVars["NextLink"] = strings.TrimSuffix(filepath.ToSlash(relNext), ".md") + ".html"
			}
		}

		// Generate HTML
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := tpl.Execute(f, pageVars); err != nil {
			return err
		}

		fmt.Println("✅ Generated:", outPath)
		return nil
	})
}

func main() {
	if err := prerender(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
	fmt.Println("✨ All stories rendered in ./public/ including static assets!")
}
