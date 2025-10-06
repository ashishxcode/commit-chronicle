# 📊 Commit Chronicle

**Generate beautiful monthly development reports from your git commits - organized by branch!**

Turn your git history into professional reports that show exactly what you worked on, organized by feature branches, bug fixes, and more.

## 🚀 Super Simple Setup

### Step 1: Copy this to your terminal

```bash
# Add this line to your ~/.zshrc file
echo 'source /Users/ashish/projects/commit-chronicle/commit-chronicle' >> ~/.zshrc

# Reload your terminal
source ~/.zshrc
```

### Step 2: Configure your repositories

Open the script file and find this section (around line 162):

```bash
REPO_PATHS=(
    "/Users/ashish/work/forked/cx-saas-dashboard"
    "/Users/ashish/work/forked/cx-saas-server"
    # Add your repository paths here
)
```

**Replace with your actual repository paths!** For example:
```bash
REPO_PATHS=(
    "/Users/yourname/projects/my-awesome-app"
    "/Users/yourname/work/company-project"
)
```

### Step 3: Done! 🎉

Now you can run `commit_chronicle` from anywhere in your terminal:

```bash
commit_chronicle
```

## 📋 What You'll Get

✅ **Branch-Based Organization** - See commits grouped by feature branches  
✅ **Professional Reports** - Clean markdown with tables and statistics  
✅ **Multi-Repository Support** - Analyze multiple projects at once  
✅ **GitHub Integration** - Includes PR data when available  
✅ **Saved to Downloads** - Reports automatically saved to your Downloads folder


### Sample Output:
```
📊 Monthly Development Report - September 2025

🌿 Commits Organized by Branch:
├── cx-saas-dashboard/feat/user-authentication (15 commits)
├── cx-saas-dashboard/fix/login-bug (3 commits)
├── my-app/feature/dark-mode (8 commits)
└── my-app/main (2 commits)
```

## 🔧 Optional: GitHub CLI Setup

For pull request data, install GitHub CLI:

**macOS:** `brew install gh`  
**Linux:** `sudo apt install gh`  
**Windows:** `winget install GitHub.cli`

Then: `gh auth login`

## ❓ Common Questions

**Q: Where do I find my repository paths?**  
A: Use `pwd` command when you're inside your project folder

**Q: Can I change the output location?**  
A: Reports are automatically saved to your Downloads folder

**Q: What if I don't see my commits?**  
A: Make sure your git name matches what you enter when running the script

**Q: Does this work on Windows?**  
A: Yes! Works on Windows Git Bash, macOS, and Linux

---

**🎯 Generate your development report in under 30 seconds!**