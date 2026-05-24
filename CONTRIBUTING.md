# Contributing to Searchbase

Contributions are welcome.

Before opening a pull request, please read:

- [`AGENTS.md`](./AGENTS.md) for architecture, data contracts, and documentation requirements.
- [`CLA.md`](./CLA.md) for the Contributor License Agreement.

## Contributor License Agreement

All contributors must accept the Searchbase Contributor License Agreement before a pull request can be merged.

The CLA Assistant workflow checks whether every pull request author has signed. If a signature is missing, the bot comments with signing instructions. Sign by commenting exactly:

```text
I have read the CLA Document and I hereby sign the CLA
```

Accepted signatures are stored in the repository on the `cla-signatures` branch under `.github/cla/signatures/v1/cla.json`.

## Documentation Requirement

If your change affects setup, configuration, provider behavior, API behavior, MCP usage, deployment, observability, or user-facing behavior, update the Hugo documentation site under `docs/`.

Run this before submitting documentation changes:

```bash
cd docs
hugo --destination /tmp/searchbase-docs-build
```
