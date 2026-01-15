package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

var md = goldmark.New()
var storyTpl = template.Must(template.ParseFiles("templates/story.html"))
var indexTpl = template.Must(template.ParseFiles("templates/index.html"))

// ChapterEntry represents a single chapter
type ChapterEntry struct {
	Number string
	URL    string
}

// StoryEntry represents a story with all its chapters (and optional cover)
type StoryEntry struct {
	Genre    string
	Story    string
	Chapters []ChapterEntry
	CoverURL string
}

// renderMarkdown converts Markdown file to HTML
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

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// copyStatic recursively copies static files
func copyStatic(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, src)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

// prerender generates HTML for all chapters and index
func prerender() error {
	contentDir := "content"
	outputDir := "public"

	_ = os.MkdirAll(outputDir, 0755)

	var allStories []StoryEntry

	// Walk genres
	genres, err := os.ReadDir(contentDir)
	if err != nil {
		return err
	}

	for _, g := range genres {
		if !g.IsDir() {
			continue
		}
		genreName := g.Name()
		stories, err := os.ReadDir(filepath.Join(contentDir, genreName))
		if err != nil {
			return err
		}

		for _, s := range stories {
			if !s.IsDir() {
				continue
			}
			storyName := s.Name()
			chapterDir := filepath.Join(contentDir, genreName, storyName)
			chapters, err := os.ReadDir(chapterDir)
			if err != nil {
				return err
			}

			// Detect cover file (cover.png, cover.jpg, cover.jpeg)
			coverPath := ""
			possibleCovers := []string{"cover.png", "cover.jpg", "cover.jpeg"}

			for _, name := range possibleCovers {
				p := filepath.Join(chapterDir, name)
				if fileExists(p) {
					coverPath = fmt.Sprintf("%s/%s/%s", genreName, storyName, name)
					break
				}
			}

			var chapterList []ChapterEntry

			for i, c := range chapters {
				if !strings.HasSuffix(c.Name(), ".md") {
					continue
				}
				num := strconv.Itoa(i + 1)

				htmlPath := fmt.Sprintf("%s/%s/%s.html", genreName, storyName, num)
				chapterList = append(chapterList, ChapterEntry{
					Number: num,
					URL:    htmlPath,
				})

				// Render chapter HTML
				mdPath := filepath.Join(chapterDir, c.Name())
				htmlContent, err := renderMarkdown(mdPath)
				if err != nil {
					return err
				}

				pageVars := map[string]any{
					"Content": template.HTML(htmlContent),
					"Title":   num,
					"Year":    time.Now().Year(),
				}

				// Next/Prev links
				numInt, _ := strconv.Atoi(num)
				prev := filepath.Join(contentDir, genreName, storyName, fmt.Sprintf("%d.md", numInt-1))
				next := filepath.Join(contentDir, genreName, storyName, fmt.Sprintf("%d.md", numInt+1))
				if fileExists(prev) {
					pageVars["PrevLink"] = fmt.Sprintf("./%d.html", numInt-1)
				}
				if fileExists(next) {
					pageVars["NextLink"] = fmt.Sprintf("./%d.html", numInt+1)
				}

				// Ensure output folder exists
				outDir := filepath.Join(outputDir, genreName, storyName)
				_ = os.MkdirAll(outDir, 0755)
				f, err := os.Create(filepath.Join(outDir, num+".html"))
				if err != nil {
					return err
				}
				if err := storyTpl.Execute(f, pageVars); err != nil {
					f.Close()
					return err
				}
				f.Close()
				fmt.Println("✅ Generated:", htmlPath)
			}

			allStories = append(allStories, StoryEntry{
				Genre:    genreName,
				Story:    storyName,
				Chapters: chapterList,
				CoverURL: coverPath,
			})
		}
	}

	// Copy static files
	if err := copyStatic("static", filepath.Join(outputDir, "static")); err != nil {
		return err
	}

	// Generate index.html
	indexFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer indexFile.Close()

	if err := indexTpl.Execute(indexFile, map[string]any{
		"Stories": allStories,
		"Year":    time.Now().Year(),
	}); err != nil {
		return err
	}

	fmt.Println("✨ All stories rendered in ./public/")
	return nil
}

func main() {
	if err := prerender(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}
