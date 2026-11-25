% "BARK-CONFIG" "1" "" "bark" "User Commands"

# NAME
bark-config \- print a config value

# SYNOPSIS
**bark config** *key*

# DESCRIPTION
Print the value bark is using for *key*: the value from the config file, or
the default if the file doesn't set it.

# KEYS
**root**
: The worktree root directory (**root_dir** in the config file).

**tmux.startup_command**
: The command sent to new tmux sessions (**startup_command** in the
  **[tmux]** table).

# EXAMPLES
```
bark config root
cd "$(bark config root)"
```

# SEE ALSO
**bark**(1), **bark**(5)
