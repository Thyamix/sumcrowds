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

### Never Push Directly to Main/Master

Never push directly to `main` or `master` branch unless the user specifically tells you to. Always create a feature branch and push to that instead.

### Creating Pull Requests

Use the Gitea CLI (`tea`) to create pull requests automatically:

```bash
tea pr create --title "type(scope): description" --description "Detailed description of changes"
```

This project uses Gitea at `git.thyamix.com`, not GitHub.

## Linear Workflow

When working on Linear issues:

1. **Assign yourself** to the issue when starting work
2. **Mark as "In Progress"** when you begin working on the issue
3. **Mark as "Awaiting Review"** once the PR is created
4. **Link the PR** to the issue using the `links` field

Example workflow:
```
1. Update issue: assignee=me, state="In Progress"
2. Create branch and implement fix
3. Create PR
4. Update issue: state="Awaiting Review", add PR link
```

## Mobile App Builds

Always build the **release/production** version of the mobile app, never the debug version:

```bash
cd mobile/android && ./gradlew assembleRelease
```

## Documentation

### Keep READMEs Up to Date

**Always update README files after making changes.** When adding new features, modifying configuration, or making significant changes, update the relevant README files before committing:

- `README.md` - Main project README with feature list and configuration
- `mobile/README.md` - Mobile app README with features and project structure

Ensure READMEs accurately reflect the current state of the codebase.

## Configuration System

The project uses a centralized configuration system with:
- **Config files** (TOML) - Non-secret settings like endpoints, CORS origins, ports
- **Env files** - Secrets like database passwords

### File Structure

```
config.dev.toml      # Development config
config.staging.toml  # Staging config
config.prod.toml     # Production config
.env.dev             # Dev secrets (gitignored)
.env.staging         # Staging secrets (gitignored)
.env.prod            # Production secrets (gitignored)
.env.example         # Template for required env vars
```

### Adding New Config Values

1. Add the value to all `config.{env}.toml` files
2. Update `backend/sharedlib/config/config.go` struct if needed
3. If it's a secret, add to `.env.{env}` files and `.env.example`

### Mobile Config Generation

Mobile config is generated from root config files. Run before building:

```bash
cd mobile
npm run generate-config:prod    # For production build
npm run generate-config:dev     # For development
```

This generates `mobile/src/config.ts` from the corresponding config file.

## Database / SQLC

### Regenerate SQLC After Query Changes

If any database schema or query changes occur (files in `backend/sharedlib/database/` or `backend/sharedlib/database/queries/`), always regenerate the SQLC code before committing:

```bash
cd backend && sqlc generate
```

This ensures the generated Go code in `backend/sharedlib/database/sqlcdb/` stays in sync with the SQL queries.

## Communication

### "Remember" Instructions

When the user says "remember something", ask for confirmation before adding it to this file. Only add to claude.md if the user confirms.
