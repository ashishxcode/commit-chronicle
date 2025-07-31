#!/bin/bash

# Monthly Development Report Generator
# Uses only git and gh commands to generate markdown report

set -e

echo "📋 Monthly Development Report Generator"
echo "========================================"

# Interactive input collection
read -p "📅 Enter month (YYYY-MM) [$(date +"%Y-%m")]: " MONTH_INPUT
MONTH=${MONTH_INPUT:-$(date +"%Y-%m")}

read -p "👤 Enter your full name [Ashish Patel]: " AUTHOR_INPUT
AUTHOR=${AUTHOR_INPUT:-"Ashish Patel"}

read -p "🐙 Enter GitHub username [ashishxcode]: " GITHUB_USERNAME_INPUT
GITHUB_USERNAME=${GITHUB_USERNAME_INPUT:-"ashishxcode"}

# Validate month format
if [[ ! $MONTH =~ ^[0-9]{4}-[0-9]{2}$ ]]; then
    echo "❌ Invalid month format. Use YYYY-MM (e.g., 2025-07)"
    exit 1
fi

# Parse month for date range
YEAR=$(echo $MONTH | cut -d'-' -f1)
MONTH_NUMBER=$(echo $MONTH | cut -d'-' -f2)
START_DATE="${YEAR}-${MONTH_NUMBER}-01"
END_DATE="${YEAR}-${MONTH_NUMBER}-31"

# Get month name
case $MONTH_NUMBER in
    01) MONTH_NAME="January" ;;
    02) MONTH_NAME="February" ;;
    03) MONTH_NAME="March" ;;
    04) MONTH_NAME="April" ;;
    05) MONTH_NAME="May" ;;
    06) MONTH_NAME="June" ;;  
    07) MONTH_NAME="July" ;;
    08) MONTH_NAME="August" ;;
    09) MONTH_NAME="September" ;;
    10) MONTH_NAME="October" ;;
    11) MONTH_NAME="November" ;;
    12) MONTH_NAME="December" ;;
esac

echo ""
echo "📊 Generating report for:"
echo "   Month: $MONTH_NAME $YEAR ($START_DATE to $END_DATE)"
echo "   Author: $AUTHOR"
echo "   GitHub: $GITHUB_USERNAME"
echo "----------------------------------------"

# Output markdown file
REPORT_FILE="${MONTH_NAME}_${YEAR}_${GITHUB_USERNAME}_report.md"

# Start building the markdown report
cat > "$REPORT_FILE" << EOF
# Monthly Development Report - $MONTH_NAME $YEAR

**Developer:** $AUTHOR  
**GitHub Username:** $GITHUB_USERNAME  
**Reporting Period:** $START_DATE to $END_DATE  
**Generated:** $(date +"%B %d, %Y at %H:%M")

---

## 📝 Commits Summary

EOF

# Step 1: Extract commits
echo "📝 Extracting commits..."
git fetch --all >/dev/null 2>&1

# Get commit count
COMMIT_COUNT=$(git log --all --author="$AUTHOR" --since="$START_DATE" --until="$END_DATE" --oneline --no-merges | wc -l | tr -d ' ')

echo "**Total Commits:** $COMMIT_COUNT" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

if [ "$COMMIT_COUNT" -gt 0 ]; then
    echo "### Detailed Commits" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    
    # Get detailed commit information
    git log --all --author="$AUTHOR" --since="$START_DATE" --until="$END_DATE" \
        --pretty=format:"Hash: %H%nDate: %ai%nMessage: %s%nBody: %b%n---" --no-merges >> "$REPORT_FILE"
    
    echo '```' >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Add commit statistics
    echo "### Commit Statistics" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    echo "Commits by day:" >> "$REPORT_FILE"
    git log --all --author="$AUTHOR" --since="$START_DATE" --until="$END_DATE" \
        --pretty=format:"%ad" --date=short --no-merges | sort | uniq -c | sort -rn >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "Most active days:" >> "$REPORT_FILE"
    git log --all --author="$AUTHOR" --since="$START_DATE" --until="$END_DATE" \
        --pretty=format:"%ad" --date=format:"%A" --no-merges | sort | uniq -c | sort -rn >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
else
    echo "No commits found for this period." >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
fi

# Step 2: Extract PR reviews
echo "🔄 Extracting PR reviews..."

cat >> "$REPORT_FILE" << EOF

---

## 🔄 Pull Request Reviews

EOF

if command -v gh &> /dev/null; then
    echo "### PRs Reviewed by Me" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Get PRs reviewed by user
    TEMP_REVIEWED=$(mktemp)
    gh pr list --search "reviewed-by:$GITHUB_USERNAME created:$START_DATE..$END_DATE" \
        --limit 100 --state all 2>/dev/null > "$TEMP_REVIEWED" || echo "" > "$TEMP_REVIEWED"
    
    if [ -s "$TEMP_REVIEWED" ]; then
        echo '```' >> "$REPORT_FILE"
        cat "$TEMP_REVIEWED" >> "$REPORT_FILE"
        echo '```' >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        
        # Count reviewed PRs
        REVIEWED_COUNT=$(wc -l < "$TEMP_REVIEWED" | tr -d ' ')
        echo "**Total PRs Reviewed:** $REVIEWED_COUNT" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    else
        echo "No PRs reviewed during this period." >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    fi
    
    echo "### PRs Authored by Me" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Get PRs authored by user
    TEMP_AUTHORED=$(mktemp)
    gh pr list --search "author:$GITHUB_USERNAME created:$START_DATE..$END_DATE" \
        --limit 100 --state all 2>/dev/null > "$TEMP_AUTHORED" || echo "" > "$TEMP_AUTHORED"
    
    if [ -s "$TEMP_AUTHORED" ]; then
        echo '```' >> "$REPORT_FILE"
        cat "$TEMP_AUTHORED" >> "$REPORT_FILE"
        echo '```' >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        
        # Count authored PRs
        AUTHORED_COUNT=$(wc -l < "$TEMP_AUTHORED" | tr -d ' ')
        echo "**Total PRs Authored:** $AUTHORED_COUNT" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    else
        echo "No PRs authored during this period." >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    fi
    
    # Clean up temp files
    rm -f "$TEMP_REVIEWED" "$TEMP_AUTHORED"
    
else
    echo "⚠️ GitHub CLI not found. Skipping PR data." >> "$REPORT_FILE"
    echo "Install with: \`brew install gh\` or visit https://cli.github.com/" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
fi

# Add footer
cat >> "$REPORT_FILE" << EOF

---

## 📊 Summary

- **Commits:** $COMMIT_COUNT
- **Period:** $START_DATE to $END_DATE
- **Generated:** $(date +"%Y-%m-%d %H:%M:%S")

*This report was generated using git and GitHub CLI commands.*
EOF

echo "✅ Report generated: $REPORT_FILE"
echo ""
echo "📄 View your report:"
echo "   cat $REPORT_FILE"