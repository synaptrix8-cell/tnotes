# Tnotes

Tnotes is a fast terminal Markdown note-taking app built with Bubble Tea v2 and Lip Gloss v2.

It stores your notes in `~/.tnotes` as `.md` files.

## Project Structure

- `cmd/tnotes`: minimal app entrypoint.
- `internal/storage`: vault path setup and note file operations.
- `internal/styles`: reusable Lip Gloss and widget styles.
- `internal/tui`: Bubble Tea model, update flow, and views.

## Install and Run

### Run directly

```bash
go run ./cmd/tnotes
```

### Install binary

```bash
go install ./cmd/tnotes
```

### Make targets

From the project root:

```bash
make run
make build
```

## Keybindings

| Key | Action | Context |
| --- | --- | --- |
| `Ctrl+N` | Create/open note prompt | Main view |
| `Ctrl+S` | Save current note and close editor | Note editor |
| `Ctrl+L` | List notes | Main view |
| `Enter` | Confirm create/open, open selected note, or confirm rename | Input/list/rename contexts |
| `r` | Rename selected note | Notes list |
| `d` | Delete selected note (with confirmation) | Notes list |
| `/` | Start list filtering/search mode | Notes list |
| `Esc` | Cancel/back (or cancel delete/rename) | Context-sensitive |
| `Ctrl+Q` / `q` | Quit app (`q` in top-level view) | Global / top-level |

## Roadmap

- Unsaved changes warning before leaving editor.
- Split-pane Markdown preview.
- Daily quick notes.
- Full-text content search.
