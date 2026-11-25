% "BARK-ADD" "1" "" "bark" "User Commands"

# NAME
bark-add \- add a worktree for a branch of a cloned repo

# SYNOPSIS
**bark add** *repo*|**.** *branch* [**--on-exists** *policy*] [**--from** *ref*]
[**--print-path**] [**--stay**] [**--worktree-root-dir** *dir*]

# DESCRIPTION
Create a worktree for *branch* at *worktree-root*/*repo*/*branch*, inside a
repo cloned with **bark-clone**(1).

*repo* must exactly match a directory under the worktree root. Pass **.** to
use the repo containing the current directory, which must be somewhere inside
*worktree-root*/*repo*/.

bark first runs **git fetch --prune origin**, so branches pushed since the
clone are found and new branches start from an up-to-date base. If the fetch
fails, for example when offline, bark logs a warning and uses the remote
branches it already has.

The branch is chosen like this:

- If *branch* exists locally, the worktree checks it out.
- Otherwise, if *branch* exists on origin, a local branch is created tracking it.
- Otherwise, a new branch is created from **--from** *ref*, or by default from
  origin's copy of the default branch (e.g. *origin/main*). The new branch has
  no upstream until you push it.

If *branch* already has a worktree, **--on-exists** decides what happens. If
the directory *worktree-root*/*repo*/*branch* exists but isn't a worktree for
*branch*, bark refuses to touch it.

After adding, the *repo* and *branch* directories are added to **zoxide**(1)
if it is installed. Unless **--stay** is passed, a tmux session named
*repo*/*branch* is created, the current tmux client is switched to it, and the
**startup_command** from **bark**(5) is sent to it. If that session already
exists, bark switches to it instead. Any "." or ":" in the session name is
replaced with "_", as tmux requires.

bark add never prompts, so it can be called from scripts. To choose a repo or
react to an existing branch interactively, build that in your script with
**bark-repos**(1) and **bark-status**(1).

# OPTIONS
**--on-exists** *policy*
: What to do when *branch* already has a worktree. **error**, the default,
  exits with status 2 and changes nothing. **open** uses the existing
  worktree: its directory is added to zoxide again and its tmux session is
  opened (unless **--stay**). **increment** applies whenever *branch* is in
  use, meaning a local or remote branch, a worktree or a directory has that
  name; bark then uses the first free *branch*-*N*, counting from 1. If
  *branch* is free it is used as given. With **--print-path** you can see
  which name was chosen.

**--from** *ref*
: Start a new branch from *ref*, such as *origin/release* or a commit,
  instead of origin's default branch. Ignored, with a warning, if *branch*
  already exists.

**--print-path**
: Print the worktree's path to stdout when done. Nothing else is printed to
  stdout.

**--stay**
: Don't create or switch to a tmux session. Use this when scripting or
  running non-interactively.

**--worktree-root-dir** *dir*
: Look for *repo* under *dir* instead of the configured worktree root.

# EXIT STATUS
**0**
: The worktree was created, or opened with **--on-exists=open**.

**1**
: Something failed.

**2**
: *branch* already has a worktree and **--on-exists** is **error**.

# EXAMPLES
```
bark add bark my-feature --stay
bark add . my-feature --stay   # from anywhere inside <worktree-root>/bark/

# Start or resume work on a ticket
bark add bark JIR-123-fix-bug --on-exists=open

# Always get a fresh branch: JIR-123-fix-bug, then -1, -2, ...
path=$(bark add bark JIR-123-fix-bug --on-exists=increment --print-path --stay)
```

# SEE ALSO
**bark**(1), **bark-clone**(1), **bark-repos**(1), **bark-status**(1),
**bark**(5)
