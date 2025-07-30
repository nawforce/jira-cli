# AI-USAGE.md

This file provides guidance for AI coding assistants (Claude Code, GitHub Copilot, etc.) when using jira-cli to interact with Jira instances.

## Basic Command Structure

```bash
# Core pattern: jira <resource> <action> [flags] [arguments]
jira issue list                    # List issues
jira issue view PROJ-123          # View specific issue
jira issue create                 # Create new issue (interactive)
jira sprint list                  # List sprints
jira epic list                    # List epics
```

## Key Principles for AI Usage

1. **Always use issue keys in full format** (e.g., `PROJ-123`, not just `123`)
2. **Check available flags with `--help`** before suggesting complex commands
3. **Use JQL for filtering** when listing issues: `jira issue list --jql "assignee = currentUser()"`
4. **Prefer interactive modes** for complex operations (create, edit)
5. **Use `--plain` flag** for machine-readable output when parsing results

## Common Workflows

### Issue Management

```bash
# List my current issues
jira issue list --jql "assignee = currentUser() AND status != Done"

# View issue details
jira issue view PROJ-123

# Create issue (interactive - best for AIs to suggest)
jira issue create

# Quick status change
jira issue move PROJ-123 --status "In Progress"

# Add comment
jira issue comment add PROJ-123 --message "Status update from AI assistant"

# Assign issue
jira issue assign PROJ-123 --assignee "user@example.com"

# Link issues
jira issue link PROJ-123 PROJ-456 --type "blocks"
```

### Sprint/Epic Management

```bash
# List active sprints
jira sprint list --state active

# List epic issues
jira epic list --table

# Add issue to epic
jira epic add EPIC-123 PROJ-456

# Add issue to sprint
jira sprint add --id 123 PROJ-456
```

### Searching and Filtering

```bash
# Use JQL for complex queries
jira issue list --jql "project = PROJ AND updated >= -7d"
jira issue list --jql "sprint in openSprints() AND assignee = currentUser()"
jira issue list --jql "status = 'To Do' AND priority = High"
jira issue list --jql "created >= -1w AND assignee = currentUser()"

# Common JQL patterns
jira issue list --jql "assignee = currentUser() AND status = 'In Progress'"
jira issue list --jql "project = PROJ AND fixVersion = 'v1.0'"
jira issue list --jql "reporter = currentUser() AND created >= -30d"
```

## AI-Specific Guidelines

### When Suggesting Commands:

1. **Ask for project context** first: "What's your project key?"
2. **Suggest `--help` flag** for unfamiliar commands
3. **Use `--plain` output** when you need to parse results
4. **Always validate issue keys** before operating on them

### Error Handling:

```bash
# Check if issue exists first
jira issue view PROJ-123 --plain > /dev/null 2>&1 && echo "exists" || echo "not found"

# Verify project exists
jira project list --plain | grep "^PROJ"

# Test connectivity
jira me > /dev/null 2>&1 && echo "connected" || echo "connection failed"
```

### Output Processing:

```bash
# Get machine-readable output for parsing
jira issue list --plain --jql "assignee = currentUser()" | head -10

# Format for specific use cases
jira issue list --columns key,summary,status --jql "project = PROJ"

# Export to different formats
jira issue list --jql "project = PROJ" --table
jira issue list --jql "project = PROJ" --plain
```

## Environment Setup Hints for AIs

```bash
# Check configuration
jira init  # Run if not configured

# Test connection
jira me    # Should show current user

# Check available projects
jira project list

# Verify server info
jira serverinfo

# Check board configuration
jira board list
```

## Cloudflare Access Support

This version of jira-cli includes enhanced Cloudflare Access support for organizations using Cloudflare Zero Trust.

### Token Configuration

```bash
# Manual token (highest priority)
export CF_ACCESS_TOKEN="your-cloudflare-access-token"

# Automatic token generation (requires cloudflared)
export CF_ACCESS_TOKEN="auto"

# No token (normal operation)
# CF_ACCESS_TOKEN unset or empty
```

### How It Works

When `CF_ACCESS_TOKEN="auto"`, jira-cli automatically:
1. Runs `cloudflared access login <jira-server-url>`
2. Opens a browser for Cloudflare authentication
3. Extracts the token from cloudflared output
4. Includes `cf-access-token` header in all requests

### Troubleshooting Cloudflare Issues

```bash
# Debug Cloudflare token issues
jira --debug me  # Shows token retrieval warnings

# Test cloudflared availability
cloudflared --version

# Manual token generation
cloudflared access login https://your-jira.company.com
```

### AI Guidance for Cloudflare Environments

When suggesting commands to users in Cloudflare-protected environments:

1. **Check for Cloudflare errors**: 403 Forbidden or CF-specific error messages
2. **Suggest automatic mode first**: `export CF_ACCESS_TOKEN="auto"`
3. **Verify cloudflared is installed**: `which cloudflared`
4. **Fallback to manual token**: If cloudflared is unavailable
5. **Debug with --debug flag**: To see token-related warnings

### Configuration Examples

```bash
# For most users (easiest setup)
export CF_ACCESS_TOKEN="auto"
jira issue list

# For automation/CI (manual token)
export CF_ACCESS_TOKEN="eyJhbGciOiJSUzI1NiI..."
jira issue list

# Temporary override
CF_ACCESS_TOKEN="auto" jira issue view PROJ-123
```

## Advanced Patterns

### Bulk Operations (Use with Caution)

```bash
# Bulk status update (always confirm with user first)
jira issue list --jql "status = 'To Do' AND assignee = currentUser()" --plain | \
  head -5 | xargs -I {} jira issue move {} --status "In Progress"

# Bulk assignment
jira issue list --jql "status = 'To Do' AND assignee is EMPTY" --plain | \
  head -3 | xargs -I {} jira issue assign {} --assignee "user@example.com"
```

### Custom Field Handling

```bash
# Create issue with custom fields
jira issue create --custom field1=value1,field2=value2

# View available custom fields
jira issue createmeta --project PROJ
```

### Template Usage

```bash
# Use templates for consistent issue creation
jira issue create --template story
jira issue create --template bug
jira issue create --template task
```

## Things AIs Should Avoid

- ❌ **Never guess issue keys** - always verify they exist
- ❌ **Don't assume project structure** - check available fields/workflows first
- ❌ **Avoid destructive operations** without explicit user confirmation
- ❌ **Don't hardcode server URLs** - use configured defaults
- ❌ **Don't bulk-modify without user approval** - especially status changes
- ❌ **Avoid assuming field names** - different Jira instances have different custom fields

## Helpful Debugging

```bash
# Debug mode for troubleshooting
jira --debug issue view PROJ-123

# Check current configuration
jira me
jira serverinfo

# Validate JQL syntax
jira issue list --jql "invalid syntax" --dry-run  # If available
```

## Integration Examples

### Git Integration

```bash
# Find Jira issues mentioned in recent commits
git log --oneline -10 | grep -E "[A-Z]+-[0-9]+" | while read line; do
  issue=$(echo $line | grep -oE "[A-Z]+-[0-9]+")
  echo "Issue: $issue"
  jira issue view $issue --plain
done

# Update issue when pushing commits
git log -1 --pretty=format:"%s" | grep -E "[A-Z]+-[0-9]+" | while read line; do
  issue=$(echo $line | grep -oE "[A-Z]+-[0-9]+")
  jira issue comment add $issue --message "Code updated: $(git log -1 --pretty=format:'%h %s')"
done
```

### Daily Standup Helper

```bash
# What am I working on?
echo "=== In Progress ==="
jira issue list --jql "assignee = currentUser() AND status = 'In Progress'" \
  --columns key,summary

echo "=== Recently Updated ==="
jira issue list --jql "assignee = currentUser() AND updated >= -1d" \
  --columns key,summary,status | head -5
```

### Sprint Planning

```bash
# Show current sprint workload
jira sprint list --current --table

# Estimate story points for current sprint
jira issue list --jql "sprint in openSprints() AND assignee = currentUser()" \
  --columns key,summary,storypoints
```

### Release Management

```bash
# Issues ready for release
jira issue list --jql "status = Done AND fixVersion = 'v1.0'" \
  --columns key,summary,assignee

# Blockers for release
jira issue list --jql "priority = Blocker AND status != Done" \
  --columns key,summary,assignee,status
```

## Common JQL Patterns for AIs

```bash
# My work
"assignee = currentUser()"

# Recent activity
"updated >= -7d"
"created >= -1w"

# Sprint-related
"sprint in openSprints()"
"sprint in closedSprints()"

# Status filtering
"status != Done"
"status in ('To Do', 'In Progress')"

# Priority filtering
"priority in (High, Highest)"
"priority = Blocker"

# Project and version
"project = PROJ AND fixVersion = '1.0'"
"project in (PROJ1, PROJ2)"

# Combining conditions
"assignee = currentUser() AND status = 'In Progress' AND updated >= -3d"
```

## Error Messages and Troubleshooting

### Common Issues AIs Should Help Debug:

1. **Authentication errors** → Suggest `jira init`
2. **Invalid issue keys** → Verify with `jira issue view KEY`
3. **JQL syntax errors** → Simplify query or check Jira documentation
4. **Permission errors** → Check user permissions in Jira web interface
5. **Network issues** → Test with `jira me` or `jira serverinfo`

### When to Suggest Alternative Approaches:

- If bulk operations fail → Suggest smaller batches
- If JQL is complex → Break into simpler queries
- If custom fields aren't working → Check available fields with createmeta
- If transitions fail → Check available transitions with `jira issue transitions KEY`

## Best Practices for AI Responses

1. **Always provide context** for suggested commands
2. **Explain what each command does** before suggesting it
3. **Ask for confirmation** before destructive operations
4. **Provide fallback options** if the primary approach fails
5. **Include relevant flags** like `--help`, `--plain`, `--table` as appropriate
6. **Suggest testing with a single item** before bulk operations