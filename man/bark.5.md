% "BARK" "5" "" "bark" "File Formats"

# NAME
bark \- configuration file for bark

# SYNOPSIS
*$XDG_CONFIG_HOME/bark/bark.toml*

*~/.config/bark/bark.toml*

# DESCRIPTION
**bark**(1) reads its configuration from a TOML file at
*$XDG_CONFIG_HOME/bark/bark.toml*, or *~/.config/bark/bark.toml* when
**XDG_CONFIG_HOME** is unset.

The file is optional. Every key is optional and falls back to its default.
Unknown keys are an error, so typos are caught rather than ignored. An empty
file is allowed and logs a warning.

Print the value bark is using for a key with **bark config** *key*; see
**bark-config**(1).

# KEYS
**root_dir** (string)
: The worktree root: the directory repos are cloned into and worktrees are
  created under, as *root_dir*/*repo*/*branch*. Default: *~/worktrees*.
  Shown by **bark config root**. **bark clone** and **bark add** can override
  it with **--worktree-root-dir**.

### [tmux]
**startup_command** (string)
: A command sent to each new tmux session that **bark clone** and **bark add**
  create, e.g. to open an editor. Default: empty (send nothing).
  Shown by **bark config tmux.startup_command**.

# EXAMPLE
```
root_dir = "/home/me/src"

[tmux]
startup_command = "nvim ."
```

# FILES
*$XDG_CONFIG_HOME/bark/bark.toml*
: The configuration file.

# SEE ALSO
**bark**(1), **bark-config**(1)
