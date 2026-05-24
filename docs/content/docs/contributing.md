---
title: "Contributing"
weight: 60
---
# Contributing

Contributions are welcome. Useful areas include new search providers, crawler improvements, API ergonomics, deployment examples, observability, and documentation.

Please see `AGENTS.md` for internal system architecture and design decisions before making changes to endpoints, data contracts, deployment shape, or service boundaries.

## Development

Searchbase uses [Task](https://taskfile.dev/) to automate development and deployment tasks.

## Common Tasks

You can run these tasks from the root directory.

| Command | Description |
| :--- | :--- |
| `task dev:up` | Start the local development environment (Docker Compose). |
| `task dev:down` | Stop the local development environment. |
| `task docs:hugo` | Run the documentation site locally. |

## Service-Specific Tasks

Each service has its own `Taskfile.yml` for more granular control. For example, in the `search-gateway` directory:

- `task test:gateway`: Run Go tests.
- `task test:fmt`: Check code formatting.
- `task fmt`: Format Go code.
- `task docs:swagger`: Regenerate Swagger/OpenAPI documentation.
- `task build`: Build the `search-gateway` binary locally.

## AI Disclosure & Contribution Policy

I believe in transparency regarding how this project is built. While AI tools (like Copilot, Claude, or local LLMs) are extensively used in the development of Searchbase, **the core architecture, design decisions, and critical business logic are human-led and human-written** with AI assistance. AI acts as an assistant, not a replacement.

However, I do heavily delegate boilerplate tasks, CI/CD pipeline generation, testing setups, and documentation writing to AI agents to maximize efficiency.

When contributing to this project, using AI to assist you is absolutely fine and encouraged! However, please ensure you deeply understand the code you are submitting. **Low-effort, blindly copy-pasted AI-generated Pull Requests will not be accepted.** Quality, intent, and architectural alignment matter far more than speed.
