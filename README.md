## 🐦 WingIt-MCP

> An MCP server that helps birders *wing it* — generating personalized target checklists for new lifers.

[![Release](https://img.shields.io/github/v/release/kpb/wingit-mcp?include_prereleases&display_name=tag&sort=semver)](https://github.com/kpb/wingit-mcp/releases)

WingIt-MCP is a Go-based Model Context Protocol server that combines your personal eBird checklists with 
recent community sightings to recommend nearby species you’ve never seen — your next lifers — in a ready-to-use field checklist.

WingIt-MCP is written in idiomatic Go and implements the [Model Context Protocol (MCP)](https://modelcontextprotocol.io)
to make its tools and resources discoverable to AI hosts like Claude Desktop.

## Roadmap

### v0.2.0 – Live eBird data + MCP polish
Planned improvements for the next release:

- Live eBird API support for recent nearby sightings (env-gated)
- Simple caching to reduce API calls and improve demo repeatability
- Clear runtime mode selection:
  - offline fixtures via `WINGIT_RECENT_JSON`
  - live API via `WINGIT_EBIRD_TOKEN`
- MCP integration smoke test (tools / prompts / resources)
- Friendlier output formatting (frequency bands: high / medium / low)

Track progress in the GitHub milestone: **v0.2.0 – Live eBird data + MCP polish**.
