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
**bark clone** and **bark add** register the new directories with
**zoxide**(1), if installed. Unless **--stay** is passed, they also create a
tmux session named *repo*/*branch*, switch the current tmux client to it, and
send it the **startup_command** from the config file. Pass **--stay** when
scripting or running non-interactively.

# ENVIRONMENT
**LOG_LEVEL**
: Log level for bark's JSON logs on stdout: **debug**, **info**, **error**,
  or the default, **warn**.

**XDG_CONFIG_HOME**
: Where to look for the config file; see **bark**(5).

# SEE ALSO
**bark**(5), **bark-add**(1), **bark-clone**(1), **bark-config**(1),
**git-worktree**(1)
