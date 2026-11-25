% "BARK-ADD" "1" "" "bark" "User Commands"

# NAME
bark-add \- add a worktree for a branch of a cloned repo

# SYNOPSIS
**bark add** *repo*|**.** *branch* [**--stay**] [**--worktree-root-dir** *dir*]

# DESCRIPTION
Create a worktree for *branch* at *worktree-root*/*repo*/*branch*, inside a
repo cloned with **bark-clone**(1).

*repo* must exactly match a directory under the worktree root. Pass **.** to
use the repo containing the current directory, which must be somewhere inside
*worktree-root*/*repo*/.

The branch is chosen like this:

- If *branch* exists locally, the worktree checks it out.
- Otherwise, if *branch* exists on origin, a local branch is created tracking it.
- Otherwise, a new branch is created from the repo's current HEAD.

bark refuses to add a worktree whose directory already exists.

After adding, the *repo* and *branch* directories are added to **zoxide**(1)
if it is installed. Unless **--stay** is passed, a tmux session named
*repo*/*branch* is created, the current tmux client is switched to it, and the
**startup_command** from **bark**(5) is sent to it.

# OPTIONS
**--stay**
: Don't create or switch to a tmux session. Use this when scripting or
  running non-interactively.

**--worktree-root-dir** *dir*
: Look for *repo* under *dir* instead of the configured worktree root.

# EXAMPLES
```
bark add bark my-feature --stay
bark add . my-feature --stay   # from anywhere inside <worktree-root>/bark/
```

# SEE ALSO
**bark**(1), **bark-clone**(1), **bark**(5)
