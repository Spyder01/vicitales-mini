# 📖 Markdown Story Blog

A lightweight **Go-powered story blog** that serves **numeric-chapter Markdown files** dynamically, with **Next/Previous navigation**, **light/dark mode**, and a clean, book-like reading experience.

Ideal for authors who write in Markdown and want to publish stories online with minimal setup.

---

## ✨ Features

* Chapters must be **strictly numeric** (`1.md`, `2.md`, …).
* Serve Markdown files dynamically via Go.
* Automatic **Next / Previous buttons** for chapter navigation.
* **Light/Dark mode** support using `prefers-color-scheme`.
* Star-divider styling for headers.
* Fully responsive and mobile-friendly.

---

## 🗂️ Folder Structure

```
markdown-story-blog/
├── content/
│   └── the-child/
│       ├── 1.md
│       ├── 2.md
│       └── 3.md
├── templates/
│   └── story.html        # HTML template
├── static/
│   └── style.css         # CSS for light/dark mode and navigation
├── main.go               # Go server
├── go.mod
└── Dockerfile            # Optional for deployment
```

---

## ⚡ Getting Started

### Requirements

* Go 1.23+

---

### Run Locally

```bash
# Clone the repo
git clone https://github.com/spyder01/vicitales-mini.git
cd <repo-name>

# Install dependencies
go mod tidy

# Run server
go run main.go
```

Open [http://localhost:8080](http://localhost:8080) to view your stories.

---

## 📝 Adding Stories

* **Filenames must be numeric**: `1.md`, `2.md`, etc.
* Each chapter corresponds to `/stories/<folder>/<number>` URL.
* Example:

  ```
  content/the-child/1.md  → /stories/the-child/1
  content/the-child/2.md  → /stories/the-child/2
  ```
* **Navigation**: Next/Previous links are generated automatically based on the numeric order.

---

## 🎨 Styling

* Light/Dark mode automatically detected.
* Responsive design for mobile and desktop.
* Star-divider separates sections elegantly.

---

## ❤️ Contributing

1. Fork the repository
2. Add new numeric-chapter Markdown files to `content/`
3. Submit a pull request

---

## 📜 License

MIT License © 2025 Suhan Bangera

---
