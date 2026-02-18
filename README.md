## 🐦 WingIt-MCP

> An MCP server that helps birders *wing it* — generating personalized target checklists for new lifers.

WingIt-MCP is a Go-based Model Context Protocol server that combines your personal eBird checklists with 
recent community sightings to recommend nearby species you’ve never seen — your next lifers — in a ready-to-use field checklist.

WingIt-MCP is written in idiomatic Go and implements the [Model Context Protocol (MCP)](https://modelcontextprotocol.io)
to make its tools and resources discoverable to AI hosts like Claude Desktop.

## Usage

WingIt-MCP needs your personal checklist and a source of recent observations.

### Required environment variables

- `WINGIT_PERSONAL_JSON`: path to your personal eBird checklist export (JSON).

### Recent observations

You can run in one of two modes:

- **Offline**: set `WINGIT_RECENT_JSON` to a JSON file of recent observations.
- **Live**: set `WINGIT_EBIRD_TOKEN` to your eBird API token and pass `location` as `"lat,lon"`.

Examples:

- `location="42.47,-76.45"`
- `location=" 42.47 , -76.45 "`

Latitude must be between -90 and 90, and longitude between -180 and 180.
