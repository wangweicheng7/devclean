# devclean

🧹 **devclean** is a cross-platform developer environment cleanup tool.

It helps you **safely scan and clean development junk files** such as build caches,
dependency directories, and temporary artifacts — **without touching source code or git history**.

---

## ✨ Features

- 🔍 Scan development junk and report disk usage
- 🗑️ Clean build caches and dependency directories safely
- 🛡️ Never deletes `.git` or source files
- 🚀 Cross-platform (macOS / Linux / Windows)
- 📦 Distributed via Homebrew

---

## 📦 Installation

### Homebrew (recommended)

```bash
brew tap wangweicheng7/tap
brew install devclean
```

### Verify:

```bash
devclean --help
```

## 🚀 Quick Start
### Scan junk files (no deletion)
```bash
devclean scan
```

### Example output:

```bash
Xcode DerivedData        12435.22 MB
Node.js node_modules     8123.10 MB
Flutter build cache      2311.44 MB
--------------------------------
Total                    22869.76 MB
```

### Clean safely (interactive)
```
devclean clean
```

### Dry-run (recommended)
```
devclean clean --dry-run
```

### Non-interactive (scripts / cron)
```
devclean clean --yes
```

## 📖 Commands
`devclean scan`

### Scan development junk files and directories.
- Read-only
- No deletion
- Safe to run anytime

---

`devclean clean`

### Clean scanned junk files.

### Flags:

- --dry-run Show what would be deleted
- --yes Skip confirmation
- --rule Only clean specific rules (comma-separated)

### Example:

```
devclean clean --rule xcode,node
```

---

`devclean doctor` (coming soon)

Diagnose common development environment issues.

---

`devclean stats` (coming soon)

Show historical cleanup statistics.

---

## 🛡️ Safety Design

`devclean` is designed with multiple safety layers:

- Never deletes .git
- Never follows symlinks
- Never deletes $HOME or /
- Rule-based cleanup only
- Interactive confirmation by default

---

## ⚙️ Configuration (planned)

Future versions will support:

```
# ~/.devclean.yaml
ignore:
  - ~/Projects/important-project
rules:
  - node
  - flutter
```

## 🧑‍💻 Development

### Requirements:

- Go 1.22+

### Build:

- go build ./...

## 📄 License

MIT License
