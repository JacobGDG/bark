% "BARK-CLONE" "1" "" "bark" "User Commands"

# NAME
bark-clone \- clone a repo into the worktree root

# SYNOPSIS
**bark clone** *url* [**--stay**] [**--worktree-root-dir** *dir*]

# DESCRIPTION
Clone *url* as a bare repo, then check out a worktree for the remote's default
branch:

```
<worktree-root>/<repo>/
  .git/      a bare clone
  <branch>/  worktree for the default branch, named after it (e.g. main/)
```

*repo* is the last path segment of *url*, without a trailing *.git*. bark
checks the remote exists first, and refuses to clone if
*worktree-root*/*repo* already exists. The default branch is set to track
its upstream on origin.

After cloning, the *repo* and *branch* directories are added to **zoxide**(1)
if it is installed. Unless **--stay** is passed, a tmux session named
*repo*/*branch* is created, the current tmux client is switched to it, and the
**startup_command** from **bark**(5) is sent to it.

# OPTIONS
**--stay**
: Don't create or switch to a tmux session. Use this when scripting or
  running non-interactively.

**--worktree-root-dir** *dir*
: Clone under *dir* instead of the configured worktree root.

# EXAMPLES
```
bark clone git@github.com:JacobGDG/bark.git --stay
cd "$(bark config root)/bark/main"
```

# SEE ALSO
**bark**(1), **bark-add**(1), **bark**(5)
