package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

var md = goldmark.New()

var (
	storyTpl = template.Must(template.ParseFiles("templates/story.html"))
	indexTpl = template.Must(template.ParseFiles("templates/index.html"))
)

/* ------------------------------
   DATA STRUCTURES
------------------------------ */

type ChapterEntry struct {
	Number string
	URL    string
}

type StoryEntry struct {
	Genre    string
	Story    string
	Chapters []ChapterEntry
	CoverURL string
}

type GenreGroup struct {
	Name    string
	Stories []StoryEntry
}

/* ------------------------------
   UTILITIES
------------------------------ */

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

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

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
		return copyFile(path, target)
	})
}

/* ------------------------------
   MAIN RENDERER
------------------------------ */

func prerender() error {
	contentDir := "content"
	outputDir := "public"

	_ = os.MkdirAll(outputDir, 0755)

	var allStories []StoryEntry

	genres, err := os.ReadDir(contentDir)
	if err != nil {
		return err
	}

	for _, g := range genres {
		if !g.IsDir() {
			continue
		}

		genreName := g.Name()
		stories, _ := os.ReadDir(filepath.Join(contentDir, genreName))

		for _, s := range stories {
			if !s.IsDir() {
				continue
			}

			storyName := s.Name()
			storyDir := filepath.Join(contentDir, genreName, storyName)

			files, _ := os.ReadDir(storyDir)

			/* ---- DETECT COVER ---- */
			coverPath := ""
			for _, ext := range []string{"png", "jpg", "jpeg"} {
				p := filepath.Join(storyDir, "cover."+ext)
				if fileExists(p) {
					coverPath = fmt.Sprintf("%s/%s/cover.%s", genreName, storyName, ext)
					// Copy to public
					dst := filepath.Join(outputDir, coverPath)
					os.MkdirAll(filepath.Dir(dst), 0755)
					_ = copyFile(p, dst)
					break
				}
			}

			/* ---- COLLECT CHAPTERS ---- */
			var chapterList []ChapterEntry
			var mdFiles []string

			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".md") {
					mdFiles = append(mdFiles, f.Name())
				}
			}

			sort.Strings(mdFiles)

			for _, fname := range mdFiles {
				chNum := strings.TrimSuffix(fname, ".md")
				numInt, _ := strconv.Atoi(chNum)

				htmlPath := fmt.Sprintf("%s/%s/%d.html", genreName, storyName, numInt)
				chapterList = append(chapterList, ChapterEntry{
					Number: chNum,
					URL:    htmlPath,
				})

				// Render chapter markdown
				mdPath := filepath.Join(storyDir, fname)
				htmlContent, _ := renderMarkdown(mdPath)

				pageVars := map[string]any{
					"Content": template.HTML(htmlContent),
					"Title":   storyName + " — Chapter " + chNum,
					"Year":    time.Now().Year(),
				}

				// Prev/Next navigation
				prevFile := filepath.Join(storyDir, fmt.Sprintf("%d.md", numInt-1))
				nextFile := filepath.Join(storyDir, fmt.Sprintf("%d.md", numInt+1))

				if fileExists(prevFile) {
					pageVars["PrevLink"] = fmt.Sprintf("./%d.html", numInt-1)
				}
				if fileExists(nextFile) {
					pageVars["NextLink"] = fmt.Sprintf("./%d.html", numInt+1)
				}

				// Write output HTML
				outDir := filepath.Join(outputDir, genreName, storyName)
				os.MkdirAll(outDir, 0755)

				f, _ := os.Create(filepath.Join(outDir, fmt.Sprintf("%d.html", numInt)))
				_ = storyTpl.Execute(f, pageVars)
				f.Close()

				fmt.Println("✔ Built:", htmlPath)
			}

			allStories = append(allStories, StoryEntry{
				Genre:    genreName,
				Story:    storyName,
				Chapters: chapterList,
				CoverURL: coverPath,
			})
		}
	}

	/* ------------------------------
	   GROUP BY GENRE FOR TEMPLATE
	------------------------------ */
	genreMap := make(map[string][]StoryEntry)
	for _, s := range allStories {
		genreMap[s.Genre] = append(genreMap[s.Genre], s)
	}

	var grouped []GenreGroup
	for g, list := range genreMap {
		grouped = append(grouped, GenreGroup{g, list})
	}

	/* ------------------------------
	   COPY STATIC FILES
	------------------------------ */
	_ = copyStatic("static", filepath.Join(outputDir, "static"))

	/* ------------------------------
	   BUILD INDEX.HTML
	------------------------------ */
	indexFile, _ := os.Create(filepath.Join(outputDir, "index.html"))
	defer indexFile.Close()

	indexTpl.Execute(indexFile, map[string]any{
		"Genres": grouped,
		"Year":   time.Now().Year(),
	})

	return nil
}

func main() {
	if err := prerender(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
	fmt.Println("✨ Site built in ./public/")
}
