# 📊 Commit Chronicle - Monthly Development Report Generator

A powerful, cross-platform zsh script that generates **professional markdown reports** from your git commit history and GitHub activity. Perfect for developers, team leads, and managers who need comprehensive monthly development summaries.

## ✨ Features

- 📈 **Professional Markdown Reports** - Beautifully formatted with tables, charts, and structured data
- 🔍 **Smart Author Detection** - Case-insensitive matching for names and usernames
- 📁 **Multi-Repository Support** - Analyze commits across multiple repositories
- 🎯 **Executive Summaries** - Key metrics and performance indicators
- 📊 **Detailed Statistics** - Commit patterns, active days, and activity breakdown
- 🔄 **GitHub Integration** - Pull request reviews and authored PRs
- 🌐 **Cross-Platform** - Works on Windows, macOS, and Linux
- ⚙️ **Zero Configuration** - Uses current directory by default
- 🎨 **Rich Formatting** - Tables, blockquotes, emojis, and proper markdown structure

## 🚀 Quick Start

1. **Make executable and add to PATH:**
   ```bash
   chmod +x commit-chronicle
   export PATH="$PATH:$(pwd)"  # Add to your ~/.zshrc for permanent access
   ```

2. **Run from any directory:**
   ```bash
   commit-chronicle
   ```

3. **Follow the prompts:**
   - Month (YYYY-MM) 
   - Your full name (for git commits)
   - GitHub username (for PR data)

4. **Get your professional report:**
   - 📄 `MonthName_YYYY_username_report.md`
   - Ready for management, HR, or portfolio use

## 📋 Prerequisites

### Required Tools
- **Zsh** - Modern shell (default on macOS 10.15+, available on all platforms)
- **Git** - For commit extraction (usually pre-installed)
- **GitHub CLI (gh)** - For PR data (optional but recommended)

### Installation

#### macOS/Linux
```bash
# Install GitHub CLI
brew install gh
```

#### Windows

**Option 1: Git Bash (Recommended)**
- Install [Git for Windows](https://git-scm.windows.com/) (includes Git Bash)
- Install [GitHub CLI](https://cli.github.com/)

**Option 2: WSL**
```bash
sudo apt update
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh
```

**Option 3: PowerShell with Package Managers**
```powershell
# Using Chocolatey
choco install git gh

# Using winget
winget install Git.Git GitHub.cli
```

## 🔧 Setup & Configuration

### 1. Repository Configuration

By default, the script analyzes the **current directory**. To analyze multiple repositories, edit the `REPO_PATHS` array in the script:

```bash
# Edit the script around line 100
REPO_PATHS=(
    "/Users/username/projects/repo1"
    "/Users/username/projects/repo2" 
    "../other-project"
    "/path/to/work-repos/api-service"
)
```

### 2. GitHub CLI Authentication (Optional but Recommended)
```bash
gh auth login
```
*Required for PR review data. Without this, commit data will still be extracted.*

## 📊 Usage Examples

### Basic Usage
```bash
# Run from any git repository
cd /path/to/your/project
commit-chronicle

# Or run from anywhere if repositories are configured in script
commit-chronicle
```

### Sample Interactive Session
```
📋 Monthly Development Report Generator
========================================
🖥️  Platform: macos (Darwin)

📅 Enter month (YYYY-MM) [2025-09]: 2025-08
👤 Enter your full name [John Doe]: John Doe  
🐙 Enter GitHub username [johndoe]: johndoe

📂 Using current directory: /Users/john/projects/my-app

📊 Generating report for:
   Month: August 2025 (2025-08-01 to 2025-08-31)
   Author: John Doe
   GitHub: johndoe
   Repositories: 1 repo(s)
     - my-app (/Users/john/projects/my-app)
----------------------------------------
📝 Extracting commits from 1 repositories...
  🔍 Processing: my-app
🔄 Extracting PR reviews...
✅ Report generated: August_2025_johndoe_report.md
```

## 📄 Report Structure

The generated markdown report includes:

### 📊 **Executive Summary**
- Total commits, repositories analyzed, active days
- Key performance indicators in table format

### 📝 **Detailed Commits by Repository**  
- Repository-specific commit breakdowns
- Commit details in organized code blocks
- Author information and commit messages

### 📊 **Commit Statistics**
- Daily commit activity table
- Most active days of the week
- Visual patterns and trends

### 🔄 **Pull Request Reviews**
- PRs reviewed by you (with GitHub CLI)
- PRs authored by you  
- Professional table format with status

### 📈 **Final Summary & Metrics**
- Comprehensive KPI breakdown
- Repository summary table
- Technical notes and methodology

## 🎨 Sample Report Output

```markdown
# 📊 Monthly Development Report - August 2025

> **Developer:** John Doe  
> **GitHub Username:** `johndoe`  
> **Reporting Period:** `2025-08-01` to `2025-08-31`  
> **Generated:** September 02, 2025 at 10:30

---

## 📝 Commits Summary

### 🎯 Executive Summary

| Metric | Value |
|--------|-------|
| **Total Commits** | 23 |
| **Repositories Analyzed** | 3 |
| **Period** | 2025-08-01 to 2025-08-31 |
| **Active Days** | 12 |

### Detailed Commits by Repository

## 📁 Repository: my-web-app

**Path:** `/Users/john/projects/my-web-app`
**Commits Found:** 15
**Primary Author Pattern:** John Doe

```
Commit: a1b2c3d
Date: 2025-08-15 10:30:00 +0000
Message: feat: add user authentication system
Author: John Doe <john@company.com>

Implemented JWT-based authentication with refresh tokens...

-------------------------------------------
```

### 📊 Commit Statistics

#### Commits by Date

| Date | Commits |
|------|---------|
| 2025-08-15 | 5 |
| 2025-08-10 | 3 |
| 2025-08-05 | 2 |

#### Most Active Days of Week

| Day | Commits |
|-----|---------|
| Tuesday | 8 |
| Monday | 6 |
| Wednesday | 5 |

---

## 🔄 Pull Request Reviews

### PRs Reviewed by Me

**Total PRs Reviewed:** 4

| PR # | Title | Status | Repository |
|------|-------|--------|------------|
| #123 | Fix login bug | MERGED | web-app |
| #124 | Add dashboard | OPEN | admin-panel |

### PRs Authored by Me

**Total PRs Authored:** 2

| PR # | Title | Status | Repository |
|------|-------|--------|------------|
| #125 | Authentication system | MERGED | web-app |
| #126 | API optimization | OPEN | backend |
```

## 🤖 AI/LLM Integration

The generated **markdown reports** are perfect for AI analysis and enhancement:

### 📋 **Management Summary Generation**
```
Analyze this monthly development report and create an executive summary highlighting:
- Key accomplishments and impact
- Productivity metrics and trends  
- Technical contributions and growth areas

[Paste generated markdown report]
```

### 📊 **Performance Analysis**
```
Based on this commit and PR data, provide insights on:
- Development velocity and patterns
- Code review participation
- Areas for improvement and recognition

[Paste markdown report sections]
```

### 🎯 **Goal Setting Assistant** 
```
Using this monthly report, suggest SMART goals for next month focusing on:
- Commit consistency and quality
- Code review engagement
- Technical skill development

[Include relevant report sections]
```

## ⚠️ Troubleshooting

### Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| **"command not found: commit-chronicle"** | Add script directory to PATH: `export PATH="$PATH:/path/to/script"` |
| **"gh: command not found"** | Install GitHub CLI or continue without PR data |
| **"No commits found"** | Check git author name matches input, verify repository has commits in period |
| **"Not a git repository"** | Run from git repository or configure `REPO_PATHS` in script |
| **Permission denied (Windows)** | Use Git Bash, WSL, or PowerShell with proper permissions |

### 🔍 **Advanced Debugging**

**Check git configuration:**
```bash
git config user.name      # Verify this matches your input
git config user.email     # Check email configuration
```

**Test git log manually:**
```bash
git log --author="Your Name" --since="2025-08-01" --until="2025-08-31" --oneline
```

**Enable debug mode:**
```bash
# Add to script after 'set -e'
set -x  # Shows all executed commands
```

### 🔧 **Author Name Issues**

The script uses **case-insensitive matching** and checks both:
- Your full name (e.g., "John Doe") 
- Your GitHub username (e.g., "johndoe")

If commits still aren't found:
1. Check different name variations in your git history
2. Use exact name from `git log --pretty=format:"%an" | sort | uniq`
3. Verify commits exist in the specified date range

## 🔄 Automation & CI/CD Integration

### 📅 **Automated Monthly Reports**

**GitHub Actions Example:**
```yaml
name: Monthly Development Report
on:
  schedule:
    - cron: '0 9 1 * *'  # 9 AM on 1st of each month
  workflow_dispatch:     # Manual trigger
  
jobs:
  generate-report:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup GitHub CLI
        run: |
          gh auth login --with-token <<< "${{ secrets.GITHUB_TOKEN }}"
      - name: Generate Monthly Report
        run: |
          chmod +x commit-chronicle
          echo "$(date -d 'last month' +%Y-%m)\n${{ github.actor }}\n${{ github.actor }}" | ./commit-chronicle
      - name: Upload Report
        uses: actions/upload-artifact@v4
        with:
          name: monthly-report
          path: "*_report.md"
```

### 🔧 **Script Customization**

**Custom Date Ranges:**
```bash
# Edit the script to change date calculations:
START_DATE="${YEAR}-${MONTH_NUMBER}-01"
END_DATE="${YEAR}-${MONTH_NUMBER}-31"

# For quarterly reports:
START_DATE="${YEAR}-01-01"  # Jan 1
END_DATE="${YEAR}-03-31"    # Mar 31
```

**Multiple Team Members:**
```bash
# Create team report script
TEAM_MEMBERS=("alice" "bob" "charlie")
for member in "${TEAM_MEMBERS[@]}"; do
    echo "2025-08\n$member\n$member" | ./commit-chronicle
done
```

## 🎨 Advanced Customization

### 📊 **Custom Report Sections**
Add your own metrics by editing the script:

```bash
# Add custom metrics to executive summary
echo "| **Lines of Code** | $(git log --author=\"$AUTHOR\" --since=\"$START_DATE\" --until=\"$END_DATE\" --pretty=tformat: --numstat | awk '{add+=$1; del+=$2} END {print add}') |" >> "$REPORT_FILE"
```

### 🏷️ **Commit Filtering**
Filter commits by type or pattern:

```bash
# Only feature commits
git log --author="$AUTHOR" --grep="feat:" --since="$START_DATE" --until="$END_DATE"

# Exclude specific patterns
git log --author="$AUTHOR" --invert-grep --grep="Merge\|merge" --since="$START_DATE" --until="$END_DATE"
```

### 📁 **Repository Groups**
Organize repositories by teams or projects:

```bash
# In script configuration
FRONTEND_REPOS=("/path/to/web-app" "/path/to/mobile-app")
BACKEND_REPOS=("/path/to/api" "/path/to/services")
REPO_PATHS=("${FRONTEND_REPOS[@]}" "${BACKEND_REPOS[@]}")
```

## 🤝 Contributing

We welcome contributions! Here are some areas for improvement:

- 🪟 **PowerShell version** for native Windows support
- 🎨 **HTML report generation** with charts and graphs  
- 📊 **Additional metrics** (lines of code, file changes, etc.)
- ⚙️ **Configuration file support** (YAML/JSON config)
- 🌐 **GitLab/Bitbucket integration** beyond GitHub
- 📱 **Mobile-friendly report formats**

### 🔧 Development Setup
```bash
git clone https://github.com/your-username/commit-chronicle
cd commit-chronicle
chmod +x commit-chronicle
./commit-chronicle  # Test the script
```

## 📄 License

MIT License - feel free to modify and distribute for your development workflow.

## 🆘 Support

- 🐛 **Issues:** [GitHub Issues](https://github.com/your-username/commit-chronicle/issues)
- 💡 **Feature Requests:** Create an issue with the `enhancement` label
- 📖 **Documentation:** Check this README and inline script comments
- 💬 **Discussions:** [GitHub Discussions](https://github.com/your-username/commit-chronicle/discussions)

---

**⭐ Star this repo if it helps you create better development reports!**