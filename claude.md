# Sumcrowds Project Guidelines

## Wiki Documentation

The `wiki/` folder contains project documentation that must stay synchronized with the codebase:

- `wiki/error-codes.md` - API error codes reference

**Important:** The wiki must always be up to date:
- When making code changes that affect documented features, update the corresponding wiki pages
- When updating wiki documentation, implement the described changes in code
- This is a bidirectional sync - code and wiki should always match

## Git Workflow

### Before Starting Any Work

Always ensure all branches are up to date before doing any work:

```bash
git fetch origin
git checkout master
git pull origin master
```

### Check for Remote Changes

Before starting work and frequently during development, check for remote changes:

```bash
git fetch origin
git status
```

### Merge Main Before Pushing

Always merge the main branch into your working branch before pushing changes:

```bash
git fetch origin
git merge origin/master
# Resolve any conflicts
git push origin <branch>
```

### Clean Up Merged Branches

After a branch has been merged, delete the local branch:

```bash
git branch -d <branch-name>
```

To delete all local branches that have been merged into master:

```bash
git branch --merged master | grep -v "master" | xargs -r git branch -d
```

### Push Reliability

Pushing commits can be unreliable. Always attempt up to 3 times before reporting a failure:

```bash
git push origin <branch> || git push origin <branch> || git push origin <branch>
```

### Commit and Push When Done

Always commit and push changes when finished with a task, unless explicitly told otherwise. Push to the current branch or create a new branch as appropriate.

## Communication

### "Remember" Instructions

When the user says "remember something", ask for confirmation before adding it to this file. Only add to claude.md if the user confirms.
