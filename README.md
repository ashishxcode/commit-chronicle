# 📓 commit-chronicle

Find your commits **and** pull requests across all your repos for a date range,
pick what matters in an interactive terminal UI, and export a clean **worklog**
(Markdown or JSON).

A single self-contained binary — no `node`, `python`, `fzf`, or `gum` to
install. The only requirement is **git**; **gh** (optional) unlocks PR & review
discovery.

```
range  →  PICK (fuzzy filter + multi-select + live preview)  →  EDIT  →  EXPORT
```

For a window you choose, it gathers **everything you did**:

- commits from git history (matched by author, across all branches)
- commits on pull requests you authored
- pull requests you authored
- pull requests you reviewed

…deduped into one **tagged** picker (`commit` / `PR` / `review`).

---

## Install

### Option 1 — `go install` (needs Go 1.24+)

```bash
go install github.com/ashishxcode/commit-chronicle/cmd/commit-chronicle@latest
```

This drops the `commit-chronicle` binary in `$(go env GOPATH)/bin` (usually
`~/go/bin`). Make sure that's on your `PATH`:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc   # or ~/.bashrc
exec $SHELL
```

### Option 2 — build from source

```bash
git clone https://github.com/ashishxcode/commit-chronicle
cd commit-chronicle
make install          # builds + installs to ~/go/bin
# or: make build      # just produces ./bin/commit-chronicle
```

### Option 3 — download a release binary

Grab the binary for your OS/arch from the
[Releases](https://github.com/ashishxcode/commit-chronicle/releases) page, then:

```bash
chmod +x commit-chronicle-*        # macOS/Linux
mv commit-chronicle-* /usr/local/bin/commit-chronicle
```

Maintainers can produce all platform binaries with `make release` (output in
`dist/`).

---

## Quick start

```bash
commit-chronicle
```

On the **first run** with no repos configured, it walks you through a one-time
setup: it scans the usual places (`~/projects`, `~/work`, `~/code`, the current
folder, …), shows how many git repos each holds, lets you pick one, and offers
to remember it. It also checks whether `gh` is authenticated so PRs and reviews
can be included.

After that, just run `commit-chronicle` from anywhere:

```bash
commit-chronicle                       # pick range → pick items → edit → export
commit-chronicle --since "7 days ago"
commit-chronicle --date today --copy   # also: yesterday, "3 days ago", etc.
```

> Re-run setup any time with `commit-chronicle --setup`.

---

## Usage

Run it inside a git repo, or configure repos/roots (below) to scan many at once:

```bash
commit-chronicle                       # interactive: pick range → pick items → edit → export
commit-chronicle --since "7 days ago"
commit-chronicle --month 2026-05 --copy
commit-chronicle --date 2026-05-25 --all --format json --out ./today.json
```

### Picker keys

| Key            | Action               |
| -------------- | -------------------- |
| `↑` / `↓`      | move                 |
| `space` / `tab`| toggle selection     |
| `a`            | select / clear all   |
| `/`            | fuzzy filter         |
| `enter`        | confirm selection    |
| `q` / `esc`    | cancel               |

In the editor: `ctrl+s` save · `esc` cancel.

### Options

```
--since <when>      relative range, e.g. "7 days ago"
--from YYYY-MM-DD   start date (inclusive)
--to   YYYY-MM-DD   end date (inclusive; default today)
--month YYYY-MM     whole calendar month
--date  YYYY-MM-DD  single day
--author "Name"     author to match (default: git config user.name)
--user <login>      GitHub login for PR discovery (default: gh user)
--repos a,b,c       comma-separated repo paths
--root  ~/work      comma-separated dirs to auto-discover git repos under
--out <path>        output path (default: ~/Downloads, timestamped)
--format md|json    output format (default: md)
--all               select everything (skip the picker)
--no-edit           skip the editor step
--no-pr             skip GitHub PR + review discovery (git commits only)
--copy              copy the whole worklog to the clipboard (skips the picker)
-h, --help          show help
```

---

## Configuring which repos to scan

Repo sources are **unioned** and checked in this order:

1. `--repos a,b,c` — explicit repo paths
2. `--root ~/work` — directories to auto-discover git repos under
3. `./.commit-chronicle` — explicit repo paths (one per line)
4. `~/.config/commit-chronicle/repos` — same, global
5. `~/.config/commit-chronicle/roots` — directories to auto-discover under
6. fallback: the current directory, if it's a git repo

The most convenient setup — point it at the folder that holds your projects,
once:

```bash
mkdir -p ~/.config/commit-chronicle
echo '~/work' > ~/.config/commit-chronicle/roots
```

Now `commit-chronicle` scans every git repo under `~/work` from anywhere, with
no flags. See [`.commit-chronicle.example`](.commit-chronicle.example) for the
file format.

---

## Pull requests & reviews

PRs and reviews are **included by default** — there's no flag to turn them on.
All you need is the GitHub CLI, authenticated once:

```bash
gh auth login        # one-time
gh auth status       # verify
```

With that in place, every run gathers, alongside your commits:

- pull requests **you authored** (tag `PR`)
- pull requests **you reviewed** (tag `review`, dated by your review)
- commits on your PRs that the plain author match might miss

**Fork workflows just work.** If you push to your own `origin` fork but open
PRs and submit reviews against an `upstream` parent, discovery queries *every*
remote — so your reviews on the upstream repo are found automatically.

Pass `--no-pr` if you ever want commits only. No `gh` installed (or not
authenticated) also falls back to git-only, with a one-line note telling you how
to enable PRs/reviews.

---

## How it works

- **Commits** come from `git log --all --author=<you>` across every ref.
- **PRs / reviews** come from `gh` (the GitHub CLI). It lists your PRs in the
  window, then fetches commit/review details per-PR — GitHub searches are
  date-bounded so it only inspects PRs that could fall in range.
- **Fork-aware:** discovery follows *every* remote of a repo, not just
  `origin`. In a fork workflow you push to your `origin` fork but open PRs and
  submit reviews against the `upstream` parent, so both are queried.
- Everything is keyed by hash (commits) or repo+number (PRs) and de-duplicated,
  so a commit that shows up both in history and on a PR appears once.
- Output is grouped by date; commits, PRs and reviews each render as distinct,
  link-bearing lines.

No `gh`, or pass `--no-pr`, and it runs git-only.

---

## Requirements

- **git** (required)
- **gh**, authenticated (`gh auth login`) — optional, for PR & review discovery
- a clipboard tool for `--copy`: `pbcopy` (macOS), `wl-copy` or `xclip` (Linux)

## Security

All `git`/`gh` calls use an explicit argument vector (no shell), `gh` handles
GitHub auth (no tokens touched here), and commit/PR text is stripped of terminal
escape sequences before display. CI runs `govulncheck` on every push. See
[SECURITY.md](SECURITY.md) for details and how to report issues.

## License

See [LICENSE](LICENSE).
