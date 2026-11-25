---
name: sync-man-pages
description: Update bark's hand-written man pages in man/ to match the code. Use after adding, removing or changing a command, flag, config key or user-visible behaviour, or when asked to check or sync the man pages or docs.
---

# Sync man pages

bark's man pages are hand-written Markdown in `man/`, rendered by go-md2man at
Nix build time. Nothing generates them from the code, so they drift unless
updated alongside it.

| Code | Page |
|---|---|
| root command (`internal/cli/root.go`), env vars, logging | `man/bark.1.md` |
| `bark <cmd>` (`internal/cli/<cmd>.go`) | `man/bark-<cmd>.1.md` |
| nested `bark <a> <b>` | `man/bark-<a>-<b>.1.md` |
| config file keys (`internal/config/config.go`, `configFile`) | `man/bark.5.md` and `man/bark-config.1.md` |

## Steps

1. Find what changed: `git diff main...HEAD` plus `git diff HEAD`, focusing on
   `internal/cli/` and `internal/config/`.
2. For each change, read the code (not just the diff or the `--help` text) and
   update the matching page:
   - **New command:** copy an existing page, e.g. `man/bark-add.1.md`. Add it to
     COMMANDS and SEE ALSO in `man/bark.1.md` and to SEE ALSO in related pages.
     Make its `Long` text end with `Full manual: man bark-<cmd>`.
   - **Removed command:** delete its page and every `bark-<cmd>(1)` reference to it.
   - **Flag:** update SYNOPSIS and OPTIONS. Write flags as `**--name**`; the
     test looks for exactly that.
   - **Config key:** update KEYS and EXAMPLE in `bark.5.md`. If `bark config`
     can print it, update `bark-config.1.md` and `ValidArgs` in `config.go`.
   - **Behaviour:** update DESCRIPTION, and SIDE EFFECTS in `bark.1.md` if it
     touches git, zoxide or tmux.
3. Keep `--help` short. Detail belongs in the page; `Long` is a summary and the
   `Full manual:` pointer.
4. Check:
   ```sh
   go test ./internal/cli -run TestManPagesMatchCommands
   ```
   This checks that every command has a page, every flag is documented, every
   page has a command and every `--help` points at its page. It can't check
   prose or config keys, so re-read what you changed.
5. Preview any page you changed:
   ```sh
   nix shell nixpkgs#go-md2man -c sh -c 'go-md2man -in man/bark-add.1.md | man -l -'
   ```
6. `git add` any new page. Flakes only see tracked files, so an untracked page
   is silently left out of `nix build`.

## Page conventions

- The first line is the title block: `% "BARK-ADD" "1" "" "bark" "User Commands"`.
  Section 5 pages use `"File Formats"`. Leave the date empty.
- Sections in this order, as needed: NAME, SYNOPSIS, DESCRIPTION, COMMANDS,
  OPTIONS or KEYS, SIDE EFFECTS, ENVIRONMENT, EXAMPLES, SEE ALSO.
- NAME is `bark-add \- short description`, in lower case.
- Use `**bold**` for literal text you type and `*italic*` for placeholders.
  Don't write `<placeholder>` outside code blocks; it's dropped as HTML.
- Options and keys use definition lists: a `**--flag** *arg*` line, then a line
  starting with `: `.
