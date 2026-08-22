# Project agents

Subagent definitions for this repo live here as `<name>.md` files. Claude Code
picks them up automatically in both local and web sessions, so anything checked
in here is available to every contributor.

Each agent is a markdown file with YAML frontmatter:

```markdown
---
name: my-agent
description: When this agent should be used.
tools: Read, Grep, Glob, Bash   # optional; omit to inherit all tools
model: sonnet                   # optional
---

System prompt for the agent goes here.
```

TODO: add agent definitions (e.g. a devops/release agent).
