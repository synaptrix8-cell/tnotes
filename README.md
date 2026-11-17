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

## Download Binary (No Go Required)

1. Open the GitHub Releases page for this repository.
2. Download the archive for your platform:
	- Linux amd64 or arm64: `tnotes_*_linux_*.tar.gz`
	- macOS amd64 or arm64: `tnotes_*_darwin_*.tar.gz`
	- Windows amd64 or arm64: `tnotes_*_windows_*.zip`
3. Extract the archive.

Linux or macOS:

```bash
tar -xzf tnotes_*_linux_*.tar.gz
sudo mv tnotes /usr/local/bin/tnotes
```

For macOS, use the darwin archive name in the same command pattern:

```bash
tar -xzf tnotes_*_darwin_*.tar.gz
sudo mv tnotes /usr/local/bin/tnotes
```

Windows:

1. Right-click the `.zip` and extract it.
2. Move `tnotes.exe` to a folder you keep for CLI tools.
3. Add that folder to your PATH environment variable.
4. Open a new terminal and run `tnotes.exe`.

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
