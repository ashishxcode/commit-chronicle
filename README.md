# 📊 Commit Chronicle

**Generate beautiful development reports from your git commits - organized by branch!**

Turn your git history into professional reports that show exactly what you worked on, with support for both monthly and daily reports, organized by feature branches, bug fixes, and more.

## ✨ Features

✅ **Daily or Monthly Reports** - Choose between single-day or full month analysis
✅ **Branch-Based Organization** - See commits grouped by feature branches
✅ **Professional Markdown Reports** - Clean formatting with tables and statistics
✅ **Multi-Repository Support** - Analyze multiple projects at once
✅ **GitHub Integration** - Includes PR reviews and authored PRs
✅ **Cross-Platform** - Works on macOS, Linux, and Windows (Git Bash/WSL)
✅ **Auto-Save** - Reports saved to your Downloads folder

## 🚀 Quick Start for Teams

### Prerequisites

- **git** (required) - Usually already installed
- **gh** (optional) - GitHub CLI for PR data
- **zsh** or **bash** shell

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/commit-chronicle.git
cd commit-chronicle

# 2. Make the script executable
chmod +x commit-chronicle

# 3. Optional: Add to PATH for global access
echo 'export PATH="$PATH:'"$(pwd)"'"' >> ~/.zshrc
source ~/.zshrc
```

### Configuration

Before running, configure your repository paths by editing the `commit-chronicle` script:

```bash
# Open the script
nano commit-chronicle  # or use your preferred editor

# Find this section (around line 264):
REPO_PATHS=(
    "/Users/ashish/work/forked/cx-saas-dashboard"
    "/Users/ashish/work/forked/cx-saas-server"
    "/Users/ashish/work/cx-partners"
    "/Users/ashish/work/saas-super-admin"
)

# Replace with YOUR repository paths:
REPO_PATHS=(
    "/Users/yourname/projects/my-project"
    "/Users/yourname/work/team-repo"
    # Add as many repos as you need
)
```

**💡 Tip:** Use `pwd` inside your project folders to get the full path.

## 📖 Usage

### Interactive Mode (Recommended)

Simply run the script and follow the prompts:

```bash
./commit-chronicle
```

You'll be asked to:
1. **Choose report type**: Monthly (1) or Daily (2)
2. **Enter date**:
   - Month format: `2025-10` (for monthly)
   - Date format: `2025-10-23` (for daily)
3. **Enter your name**: e.g., "John Doe"
4. **Enter GitHub username**: e.g., "johndoe"

### Example Session

```
📋 Development Report Generator
========================================
🖥️  Platform: macos (Darwin)

📊 Report Type:
   1) Month-based (entire month)
   2) Date-based (specific day)
Choose option [1]: 2

📅 Date-based report selected

📅 Enter date (YYYY-MM-DD) [2025-10-23]: 2025-10-18
👤 Enter your full name [Ashish Patel]: John Doe
🐙 Enter GitHub username [ashishxcode]: johndoe

✅ Found repository: /Users/john/projects/my-app
----------------------------------------
🔄 Extracting commits from 1 repositories...
✅ Report generated successfully!
```

## 📋 What You'll Get

### Monthly Report Example

```markdown
# 📊 Monthly Development Report - October 2025

> **Developer:** John Doe
> **GitHub Username:** `johndoe`
> **Reporting Period:** `2025-10-01 to 2025-10-31`

---

## 🔄 Pull Request Reviews

### PRs Reviewed by Me
**Total PRs Reviewed:** 5

### PRs Authored by Me
**Total PRs Authored:** 3

---

## 📝 Commits Summary

### 🌿 Commits Organized by Branch

##### 🌱 `my-app/feat/user-authentication` (15 commits)
```
┌─ 2025-10-15 14:23:45
└─ feat: add JWT token validation

┌─ 2025-10-15 16:45:12
└─ feat: implement refresh token mechanism
```
```

### Daily Report Example

```markdown
# 📊 Daily Development Report - October 18, 2025

> **Developer:** John Doe
> **GitHub Username:** `johndoe`
> **Reporting Period:** `2025-10-18`

---

## 📝 Commits Summary

##### 🌱 `my-app/fix/login-bug` (3 commits)
```
┌─ 2025-10-18 09:15:30
└─ fix: resolve login timeout issue

┌─ 2025-10-18 11:42:18
└─ test: add unit tests for login flow
```
```

## 🔧 GitHub CLI Setup (Optional but Recommended)

To include PR review data, install and authenticate with GitHub CLI:

### Installation

**macOS:**
```bash
brew install gh
```

**Linux:**
```bash
sudo apt install gh
# or
sudo snap install gh
```

**Windows:**
```bash
winget install GitHub.cli
# or
choco install gh
```

### Authentication

```bash
gh auth login
```

Follow the prompts to authenticate with your GitHub account.

## 🛠️ Advanced Configuration

### Multiple Authors

If you use different git author names across repos, the script automatically detects and uses the most common one.

### Custom Output Location

Reports are saved to:
- **macOS/Linux:** `~/Downloads/`
- **Windows:** `%USERPROFILE%/Downloads/`

Format: `October_2025_username_report.md` or `October_18_2025_username_report.md`

### Cross-Platform Compatibility

The script automatically detects your OS and adjusts:
- Date calculations
- Path handling
- Temp file creation
- Downloads directory location

## ❓ Troubleshooting

### "No commits found"

**Possible causes:**
1. Git author name mismatch - Check with: `git config user.name`
2. Wrong date range
3. No commits in configured repositories for that period

**Fix:** Make sure the name you enter matches your git configuration.

### "Not a git repository" errors

**Fix:** Update the `REPO_PATHS` array in the script with valid repository paths.

### Permissions error

**Fix:** Make the script executable:
```bash
chmod +x commit-chronicle
```

### GitHub CLI not working

**Fix:**
1. Install: `brew install gh` (macOS) or `sudo apt install gh` (Linux)
2. Authenticate: `gh auth login`
3. Verify: `gh auth status`

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Report bugs via GitHub Issues
- Submit feature requests
- Create pull requests with improvements

## 📝 License

MIT License - Feel free to use this in your team!

## 🎯 Pro Tips

💡 **Quarterly Reports:** Use month mode for each month, then combine reports
💡 **Daily Standup:** Use date mode to quickly see yesterday's work
💡 **Performance Reviews:** Generate reports for review periods
💡 **Team Sync:** Share reports in Slack/Teams to show progress
💡 **Automation:** Create a cron job for weekly auto-generated reports

---

**Made with ❤️ for developers who want to track their work efficiently**

Need help? Open an issue on GitHub!
