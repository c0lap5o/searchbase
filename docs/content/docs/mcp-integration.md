---
title: "MCP Integration"
weight: 20
---

If your LLM client supports MCP (like Claude Desktop, Cursor, LM Studio, OpenWebUI, or Opencode), you can plug Searchbase directly into it using one of the supported transports.

- **MCP SSE Endpoint:** `http://localhost:8080/mcp/sse`
- **MCP Streamable HTTP Endpoint:** `http://localhost:8080/mcp/http`

## Available Tools

- `web_search`: Search the live web and return compact search results.
- `fetch_url`: Fetch one exact URL and extract optimized Markdown.

## OpenWebUI

OpenWebUI works best with the Streamable HTTP transport.

1. Open OpenWebUI.
2. Click your profile.
3. Go to **Admin Panel** > **Settings** > **Integrations**.
4. Click **+ Add connection**.
5. Set the connection type to **MCP Streamable HTTP**.
6. Set the ID to `searchbase`.
7. Set the target URL to `http://localhost:8080/mcp/http`.
8. Save the connection.

After connection, OpenWebUI should discover the `web_search` and `fetch_url` tools automatically.

{{< hint warning >}}
If OpenWebUI runs inside Docker, `localhost` points to the OpenWebUI container, not your host. Use the Searchbase service name if both services share a Docker network, or use your host/server address instead.
{{< /hint >}}

## LM Studio

LM Studio can connect to Searchbase through the MCP SSE endpoint.

Use this configuration in LM Studio's MCP server settings:

```json
{
  "mcpServers": {
    "searchbase": {
      "transport": "sse",
      "url": "http://localhost:8080/mcp/sse"
    }
  }
}
```

If LM Studio asks for fields instead of raw JSON, use:

| Field | Value |
| :--- | :--- |
| Name | `searchbase` |
| Transport | `sse` |
| URL | `http://localhost:8080/mcp/sse` |

Once enabled, LM Studio can call `web_search` for discovery and `fetch_url` when it needs full page content.

## Opencode

Opencode can connect to Searchbase through the remote MCP SSE endpoint.

Add the Searchbase MCP server to your `opencode.json` or `opencode.jsonc` config:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "searchbase": {
      "type": "remote",
      "url": "http://localhost:8080/mcp/sse",
      "enabled": true
    }
  }
}
```

If you restrict tools in Opencode, enable the Searchbase MCP server tools:

```json
{
  "tools": {
    "searchbase": true
  }
}
```

Replace `localhost` with your Searchbase server address when Opencode runs on a different machine.

## Claude Desktop / Cursor

Clients that support SSE MCP servers can use the same JSON shape:

```json
{
  "mcpServers": {
    "searchbase": {
      "transport": "sse",
      "url": "http://localhost:8080/mcp/sse"
    }
  }
}
```

Replace `localhost` with the actual server address if Searchbase is not running on the same machine as your MCP client.
