# 📊 Commit Chronicle

A zsh script that generates professional monthly development reports from your git commits and GitHub activity.

## 🚀 Quick Setup

### 1. Install Prerequisites

**macOS:**
```bash
brew install gh  # GitHub CLI (optional for PR data)
```

**Linux:**
```bash
sudo apt install gh
```

**Windows:**
```bash
winget install GitHub.cli
```

### 2. Setup Script

```bash
# Make executable
chmod +x commit-chronicle

# Add to PATH (permanent)
echo 'export PATH="$PATH:$(pwd)"' >> ~/.zshrc
source ~/.zshrc
```

### 3. Configure Repositories

Edit the script (around line 111):
```bash
REPO_PATHS=(
    "/path/to/your/repo1"
    "/path/to/your/repo2"
)
```

### 4. GitHub CLI (Optional)

```bash
gh auth login
```

## 📊 Usage

```bash
commit-chronicle
```

Follow the prompts:
- Month (YYYY-MM)
- Your full name
- GitHub username

**Output:** `MonthName_YYYY_username_report.md`

## 🔧 Features

- Professional markdown reports with tables and statistics
- Multi-repository analysis
- Case-insensitive author matching
- GitHub PR integration
- Cross-platform compatibility

## ⚠️ Troubleshooting

| Issue | Solution |
|-------|----------|
| Command not found | Add script directory to PATH |
| No commits found | Check git author name matches input |
| Not a git repo | Run from git repository or configure REPO_PATHS |

---

*Generate professional monthly development reports in seconds.*