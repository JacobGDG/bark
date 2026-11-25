% "BARK-STATUS" "1" "" "bark" "User Commands"

# NAME
bark-status \- show where a branch exists and its state

# SYNOPSIS
**bark status** *repo*|**.** *branch* [**--json**] [**--fetch**]
[**--worktree-root-dir** *dir*]

# DESCRIPTION
Report whether *branch* exists in *repo*, locally, on origin or as a worktree,
and if so what state it is in. *repo* is resolved as in **bark-add**(1): a
directory name under the worktree root, or **.** for the current repo.

Remote information comes from the last fetch. Pass **--fetch** to fetch
origin first. bark status never changes branches or worktrees.

A branch that doesn't exist is not an error: bark prints its status with
**exists** set to no and exits 0.

# FIELDS
Plain output prints one *field*: *value* line per field, with **-** for fields
that don't apply. **--json** prints one object with the same fields, using
**null** for fields that don't apply.

**repo**, **branch**
: The repo and branch asked about.

**exists**
: Whether the branch exists locally or on origin.

**local**, **remote**
: Whether it exists as a local branch, and as *origin*/*branch*.

**worktree**
: The path of the branch's worktree, if it has one.

**changed_files**
: The number of uncommitted changes in the worktree, including untracked
  files. Only set if there is a worktree.

**ahead**, **behind**
: Commits on the local branch that aren't on origin's copy, and the other way
  round. Only set if the branch exists both locally and on origin.

**last_commit**
: The committer date of the branch's latest commit, in RFC 3339 format. Uses
  the local branch if there is one, otherwise origin's copy.

**merged**
: Whether every commit on the branch is in origin's default branch (or the
  local default branch, if origin doesn't have it). A new branch with no
  commits of its own counts as merged. Branches merged by squash or rebase
  show as not merged, because their commits were rewritten.

# OPTIONS
**--json**
: Print the status as a JSON object.

**--fetch**
: Run **git fetch --prune origin** first. If the fetch fails, bark logs a
  warning and reports what it already has.

**--worktree-root-dir** *dir*
: Look for *repo* under *dir* instead of the configured worktree root.

# EXAMPLES
```
bark status bark my-feature

# Ask before reusing a branch, from a script
if bark status bark "$branch" --fetch --json | jq -e .exists >/dev/null; then
  ...
fi
```

# SEE ALSO
**bark**(1), **bark-add**(1), **bark-repos**(1)
