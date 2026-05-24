---
title: "Searchbase"
type: "docs"
---
# Searchbase

**SEARCHBASE** is a self-hosted, highly efficient, privacy-focused Search Engine designed explicitly for **AI Agents and LLMs**. 

It natively supports the **Model Context Protocol (MCP)**, making it a plug-and-play web search tool for LLM UIs like Claude Desktop, Cursor, and OpenWebUI.

## 🎯 The Mission

In an era where Big AI companies profit from your data and gatekeep intelligence behind expensive token economies, privacy and autonomy have become luxury commodities. While small, local LLMs have the potential to be incredibly powerful, they are currently tethered to proprietary search APIs that demand high subscriptions and offer zero transparency regarding what happens to your queries or agent interactions.

**The status quo is built on dependency:**
* **Expensive Gatekeepers:** Services like Tavily, Bing, and Google Search API force developers into tiered subscription models, creating a "pay-to-play" barrier for innovation.
* **Privacy Erosion:** Using proprietary middleware means your agents' research paths—and the data they uncover—are being fed into corporate training loops.
* **Inefficient Proxies:** Standard wrappers (like SearXNG) struggle with high-frequency, high-concurrency demands—frequently hitting rate limits and Captchas—while delivering unoptimized, token-heavy content that bloats LLM context windows.

**Searchbase is built to break this cycle.** 

Searchbase mission is to provide a sovereign, privacy-first search layer that enables any LLM—regardless of size—to access the live web with maximum efficiency and zero dependency on corporate gatekeepers. Searchbase turns the web into a clean, token-optimized stream of knowledge, ready for your context window.

This repository is my independent open-source project. It will remain separate from any managed cloud service I may build around it.

**The Searchbase Way:**
* **Economic Sovereignty:** Bypass paid API keys entirely using lightweight, intelligent scraping.
* **Token Optimization:** searchbase desn't just fetch web pages; it distills them. By stripping boilerplate and ads, it returns pristine Markdown designed specifically to minimize token usage and maximize LLM reasoning.
* **Native Integration:** With out-of-the-box **MCP** support, your AI tools gain "eyes" on the web with zero glue-code.
* **Architectural Transparency:** A developer-first REST API that provides a robust foundation for building autonomous agents without breaking the bank.

## ⚖️ License & Business Model

Searchbase is dual-licensed under the **GNU AGPLv3** and a **Commercial License**.

### My Promise: The Core Stays Open
I want to be completely transparent about my business intentions from day one. Let's face it: open-source is cool, but open-source doesn't pay the bills. To make this project sustainable, I am building and offering a managed, hosted Cloud version of Searchbase (bootstrapped, operated, and fully controlled by me).

However, **the core project will always remain free and open-source.** Regardless of the commercial cloud offering, the companies I form, or who handles the business licenses in the future, this repository is my strict commitment to the community. You will always be able to self-host the core of Searchbase for free.

### Why AGPLv3? (The "Enterprise Repellent")
I specifically chose the AGPLv3 license instead of a permissive license like MIT or Apache 2.0 to protect this open-source and independence vision. 

My goal is to keep Searchbase **100% free for homelabbers, individuals, and the open-source community**. What I *do not* want is large corporations or massive SaaS platforms taking my hard work, packaging it into a paid product, and offering it as a service without contributing back or purchasing a commercial license.

The AGPLv3 ensures that if a company modifies or uses Searchbase as part of a network service, they must open-source their entire application. 

### Dual License Options
1. **Open Source (AGPLv3):** Free for personal use, homelabs, and open-source projects. If you use it to provide a service over a network, you must comply with the AGPLv3 terms.
2. **Commercial License:** If you wish to use Searchbase in a commercial, proprietary SaaS or enterprise environment without the obligations of the AGPLv3, please contact me to purchase a commercial license.
