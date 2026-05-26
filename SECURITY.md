# Security

## Reporting a vulnerability

Please report security issues privately via
[GitHub Security Advisories](https://github.com/ashishxcode/commit-chronicle/security/advisories/new)
rather than opening a public issue. We aim to acknowledge reports within a few
days.

## Design notes

commit-chronicle is a local, single-user CLI. It reads your own repositories and
talks to GitHub through the `gh` CLI you have already authenticated. With that in
mind:

- **No shell interpolation.** All `git`/`gh` invocations use `exec.Command` with
  an explicit argument vector — never a shell — so repo paths, author names, PR
  numbers, and date strings cannot be used for command injection.
- **No credential handling.** GitHub authentication is delegated entirely to
  `gh`; this tool never reads, stores, or transmits tokens.
- **Untrusted text is sanitized.** Commit messages and PR titles come from other
  people. Before they are shown in the picker, preview, or a worklog printed to
  the terminal, `model.CleanText` strips ANSI/terminal escape sequences and
  control characters so a crafted message cannot hijack your terminal.
- **Dependency scanning.** CI runs
  [`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on every
  push to catch known vulnerabilities in the module and its dependencies.
