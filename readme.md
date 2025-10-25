# 📖 Markdown Story Blog

A lightweight **Go-powered story blog** that **prerenders numeric-chapter Markdown files** into static HTML, with **Next/Previous navigation**, **breadcrumbs**, **light/dark mode**, and a clean, book-like reading experience.

Ideal for authors who write in Markdown and want a static, deployable online story site with minimal setup.

---

## ✨ Features

* Chapters must be **strictly numeric** (`1.md`, `2.md`, …).
* **Pre-render Markdown** to static HTML (`public/` folder).
* Automatic **Next / Previous buttons** for chapter navigation.
* **Breadcrumbs** showing `Genre / Story / Chapter`.
* **Light/Dark mode** support using `prefers-color-scheme`.
* Fully responsive and mobile-friendly.
* Copies `static/` folder automatically into `public/` for CSS, JS, and images.
* Ready for **GitHub Pages** or any static hosting.

---

## 🗂️ Folder Structure

```
markdown-story-blog/
├── content/
│   ├── fantasy/
│   │   └── red-lily/
│   │       ├── 1.md
│   │       ├── 2.md
│   │       └── 3.md
├── templates/
│   └── story.html        # HTML template for each chapter
├── static/
│   └── style.css         # CSS for light/dark mode and navigation
├── main.go               # Prerender Go program
├── public/               # Generated HTML & static assets
├── go.mod
└── .github/workflows/deploy.yml  # GitHub Actions workflow
```

---

## ⚡ Getting Started

### Requirements

* Go 1.23+

---

### Generate Static Site Locally

```bash
# Clone the repo
git clone https://github.com/spyder01/vicitales-mini.git
cd <repo-name>

# Install dependencies
go mod tidy

# Prerender all stories to public/
go run main.go
```

Your HTML and static assets will be in `public/`.

---

### Serve Locally for Testing

```bash
# Option 1: Using Go
go run serve.go  # or create a simple file server serving public/

# Option 2: Using Python
cd public
python3 -m http.server 8080
```

Open [http://localhost:8080](http://localhost:8080) to view the site.

---

## 📝 Adding Stories

* **Filenames must be numeric**: `1.md`, `2.md`, etc.
* Directory structure: `content/<genre>/<story>/<chapter>.md`
* Example:

```
content/fantasy/red-lily/1.md  → public/fantasy/red-lily/1.html
content/fantasy/red-lily/2.md  → public/fantasy/red-lily/2.html
```

* **Navigation**: Next/Previous links are generated automatically based on numeric chapters.
* **Breadcrumbs**: Automatically show Genre → Story → Chapter.

---

## 🌐 Deployment

* Deploy the contents of `public/` to **GitHub Pages**, **Netlify**, **Vercel**, or any static hosting.
* Optional: Use **GitHub Actions** to automatically prerender and deploy on push. Example workflow:

```
.github/workflows/deploy.yml
```

* Push Markdown changes → GitHub Actions will generate HTML → deploy to `gh-pages` branch.

---

## 🎨 Styling

* Light/Dark mode automatically detected by system preference.
* Fully responsive design for desktop and mobile.
* Star-divider separates sections elegantly.

---

## ❤️ Contributing

1. Fork the repository
2. Add new numeric-chapter Markdown files to `content/`
3. Submit a pull request

---

## 📜 License

MIT License © 2025 Suhan Bangera

