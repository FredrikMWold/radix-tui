# radix-tui

[![Go Reference](https://pkg.go.dev/badge/github.com/FredrikMWold/radix-tui.svg)](https://pkg.go.dev/github.com/FredrikMWold/radix-tui)
[![Release](https://img.shields.io/github/v/release/FredrikMWold/radix-tui?sort=semver)](https://github.com/FredrikMWold/radix-tui/releases)

**A keyboard-first TUI for Radix applications and pipelines**, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Browse applications, inspect environments and recent jobs, open job details in your browser, and trigger Build & Deploy or Apply Config right from your terminal.

![Demo](./radix-tui-demo.gif)

<details>
	<summary><strong>Quick keys</strong></summary>

| Context | Key | Action |
|---|---|---|
| Applications | `↑`/`↓` | Move selection |
| Applications | `Enter` | Select application |
| Pipeline | `↑`/`↓` | Move selection |
| Pipeline | `Enter` | Open job in browser |
| Pipeline | `Ctrl+n` | Build & Deploy form |
| Pipeline | `Ctrl+a` | Apply Config form |
| Pipeline | `Ctrl+r` | Refresh application |
| Anywhere | `Esc` | Back (context aware) |
| Anywhere | `q`/`Ctrl+C` | Quit |

> Tip: The help footer updates based on what you can do in the current view.

</details>

## Features

- 🧭 Browse Radix applications and environments
- 📊 See recent pipeline jobs with status and timestamps
- ⚙️ Trigger Build & Deploy or Apply Config from dedicated forms
- 🌐 Open the selected job in your default browser

## Install

Install with Go:

```sh
go install github.com/FredrikMWold/radix-tui@latest
```

Or download a prebuilt binary from Releases and place it on your PATH:

- https://github.com/FredrikMWold/radix-tui/releases