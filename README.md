# Monthly Development Report Generator

A simple bash script that generates comprehensive monthly development reports using only `git` and `gh` commands. Perfect for team members to create standardized monthly reports.

## 🚀 Quick Start

1. **Clone or download this repository**
2. **Run the script:**
   ```bash
   ./extract_monthly_data.sh
   ```
3. **Follow the prompts:**
   - Enter month (YYYY-MM)
   - Enter your full name  
   - Enter your GitHub username
4. **Get your report:**
   - Complete markdown report: `MonthName_YYYY_username_report.md`
   - Ready to share with managers and HR

## 📋 Prerequisites

### Required Tools
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

## 🔧 Setup

### 1. Make Script Executable (Unix/Linux/macOS)
```bash
chmod +x extract_monthly_data.sh
```

### 2. GitHub CLI Authentication
```bash
gh auth login
```
Follow the prompts to authenticate with your GitHub account.

## 📊 Usage

### Interactive Mode (Recommended)
```bash
./extract_monthly_data.sh
```

The script will prompt you for:
- **Month**: YYYY-MM format (e.g., 2025-07)
- **Full Name**: Your name as it appears in git commits
- **GitHub Username**: Your GitHub username for PR searches

### Example Session
```
🔍 Monthly Development Data Extractor
=====================================+

📅 Enter month (YYYY-MM) [2025-07]: 2025-07
👤 Enter your full name [Ashish Patel]: John Doe
🐙 Enter GitHub username [ashishxcode]: johndoe

📊 Extracting data for:
   Month: 2025-07 (2025-07-01 to 2025-07-31)
   Author: John Doe
   GitHub: johndoe
----------------------------------------

📝 Extracting commits...
✅ Commits saved to: monthly_data_2025_07/johndoe_commits_2025-07.json

🔄 Extracting PR reviews...
✅ PR reviews saved to: monthly_data_2025_07/johndoe_pr_reviews_2025-07.json

✅ Data extraction complete!
```

## 📁 Output Structure

```
monthly_data_YYYY_MM/
├── username_commits_YYYY-MM.json      # All commits data
└── username_pr_reviews_YYYY-MM.json   # PR reviews data
```

### Commits JSON Structure
```json
[
  {
    "hash": "abc123...",
    "date": "2025-07-15T10:30:00+00:00",
    "message": "feat: add user authentication",
    "body": "Detailed commit description..."
  }
]
```

### PR Reviews JSON Structure
```json
{
  "month": "2025-07",
  "user": "username",
  "prs_reviewed": [
    {
      "number": 123,
      "title": "Fix login bug",
      "author": {"login": "colleague"},
      "state": "merged",
      "mergedAt": "2025-07-15T14:20:00Z",
      "url": "https://github.com/repo/pull/123"
    }
  ],
  "prs_authored": [...],
  "extracted_at": "2025-07-31 10:30:00"
}
```

## 🤖 Using with LLMs

The generated JSON files are perfect for LLM processing. Example prompts:

### For Monthly Reports
```
Analyze this commits and PR data to generate a professional monthly development report:

[Paste commits JSON]
[Paste PR reviews JSON]

Include: impact summary, key accomplishments, technical contributions, and metrics.
```

### For Code Review Analysis
```
Based on this PR review data, analyze my code review patterns and provide insights:

[Paste PR reviews JSON]
```

## ⚠️ Troubleshooting

### Common Issues

**1. "jq: command not found"**
```bash
# Install jq first
brew install jq  # macOS
sudo apt install jq  # Linux
```

**2. "gh: command not found"**
- The script will still work but won't extract PR data
- Install GitHub CLI for complete functionality

**3. "Invalid month format"**
- Use YYYY-MM format (e.g., 2025-07, not 07-2025)

**4. No commits found**
- Check if your git name matches what you entered
- Verify the repository has commits in the specified month
- Check git config: `git config user.name`

**5. Permission denied (Windows)**
- Use Git Bash or WSL
- Or create a PowerShell version (let me know if needed)

### Debug Mode
Add debug output by editing the script:
```bash
# Add after 'set -e'
set -x  # This will show all commands being executed
```

## 🔄 Integration with CI/CD

You can automate this for regular reporting:

```yaml
# GitHub Actions example
name: Monthly Report Data
on:
  schedule:
    - cron: '0 0 1 * *'  # First day of each month
  
jobs:
  extract:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Extract monthly data
        run: |
          echo "$(date -d 'last month' +%Y-%m)\nDeveloper Name\ngithub-username" | ./extract_monthly_data.sh
```

## 📝 Customization

### Modify Date Range
Edit the script to change the date calculation:
```bash
# For custom date ranges, modify these lines:
START_DATE="${YEAR}-${MONTH_NUMBER}-01"
END_DATE="${YEAR}-${MONTH_NUMBER}-31"
```

### Add More Git Data
Extend the git log format:
```bash
# Add more fields to the JSON output
--pretty=format:'{"hash":"%H","date":"%ai","message":"%s","body":"%b","author":"%an","email":"%ae"}'
```

### Filter Commits
Add commit filtering:
```bash
# Example: Only feature commits
git log ... | grep -E "feat:|feature:"
```

## 🤝 Contributing

Feel free to improve this script:
- Add PowerShell version for Windows
- Add more git data extraction
- Improve error handling
- Add configuration file support

## 📄 License

This script is provided as-is for development productivity. Modify as needed for your workflow.

---

**Need help?** Check the troubleshooting section or modify the script for your specific needs.