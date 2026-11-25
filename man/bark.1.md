% "BARK" "1" "" "bark" "User Commands"

# NAME
bark \- git worktree manager

# SYNOPSIS
**bark** *command* [*flags*]

**bark** [**-h**|**--help**] [**-v**|**--version**]

# DESCRIPTION
**bark** clones each repo as a bare git repo and checks out one git worktree
per branch, all under a single worktree root:

```
<worktree-root>/
  <repo>/
    .git/      bare clone (no files checked out here)
    main/      worktree for the default branch
    <branch>/  one worktree per branch added with "bark add"
```

Each *branch* directory is a normal git worktree: cd into it and run git as
usual. Branch names containing "/" become nested directories
(*feature/x* becomes *repo*/*feature/x*/).

The worktree root defaults to *~/worktrees*. Change it with **root_dir** in
the config file (see **bark**(5)), and print the current value with
**bark config root**.

# COMMANDS
**bark-clone**(1)
: Clone a repo into the worktree root.

**bark-add**(1)
: Add a worktree for a branch of a cloned repo.

**bark-repos**(1)
: List repos in the worktree root.

**bark-status**(1)
: Show where a branch exists and its state.

**bark-config**(1)
: Print a config value.

**bark completion** *shell*
: Print a shell completion script. See **bark completion --help**.

# OPTIONS
**-h**, **--help**
: Show help for bark or a command.

**-v**, **--version**
: Print the version.

# SIDE EFFECTS
**bark add** runs **git fetch --prune origin** in the repo first, and
**bark status --fetch** does the same.

**bark clone** and **bark add** register the new directories with
**zoxide**(1), if installed. Unless **--stay** is passed, they also create a
tmux session named *repo*/*branch*, switch the current tmux client to it, and
send it the **startup_command** from the config file, or just switch to that
session if it already exists. Pass **--stay** when scripting or running
non-interactively.

# SCRIPTING
**bark** is built to be called from scripts, such as a tmux popup that picks
a ticket and starts work on it. It never prompts: **bark-repos**(1) and
**bark-status**(1) give a script what it needs to ask its own questions, and
**bark-add**(1) has flags to act on the answers. Command output goes to
stdout and logs go to stderr.

```
repo=$(bark repos | fzf) || exit
bark add "$repo" "$branch" --on-exists=open
```

# ENVIRONMENT
**LOG_LEVEL**
: Log level for bark's JSON logs on stderr: **debug**, **info**, **error**,
  or the default, **warn**.

**XDG_CONFIG_HOME**
: Where to look for the config file; see **bark**(5).

# EXIT STATUS
**0** on success, **1** on failure. **bark-add**(1) exits **2** when the
branch already has a worktree and **--on-exists** is **error**.

# SEE ALSO
**bark**(5), **bark-add**(1), **bark-clone**(1), **bark-config**(1),
**bark-repos**(1), **bark-status**(1), **git-worktree**(1)
