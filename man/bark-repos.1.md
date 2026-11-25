% "BARK-REPOS" "1" "" "bark" "User Commands"

# NAME
bark-repos \- list repos in the worktree root

# SYNOPSIS
**bark repos** [**--json**] [**--worktree-root-dir** *dir*]

# DESCRIPTION
Print the name of each repo in the worktree root, one per line, sorted by
name. A repo is a directory directly under the root that contains a *.git*
entry, as created by **bark-clone**(1). Other directories and files are
skipped.

If the worktree root doesn't exist, nothing is printed (**[]** with
**--json**) and bark exits 0.

The output is meant for scripts and pickers such as **fzf**(1). Names can be
passed straight to **bark-add**(1) and **bark-status**(1).

# OPTIONS
**--json**
: Print a JSON array instead. Each repo is an object with **name**,
  **path** and **worktrees**, a list of objects with **branch** and **path**.
  **branch** is empty for a worktree with a detached HEAD.

**--worktree-root-dir** *dir*
: List repos under *dir* instead of the configured worktree root.

# EXAMPLES
```
repo=$(bark repos | fzf) && bark add "$repo" my-feature

# Every worktree path, e.g. to feed a session picker
bark repos --json | jq -r '.[].worktrees[].path'
```

# SEE ALSO
**bark**(1), **bark-add**(1), **bark-status**(1)
