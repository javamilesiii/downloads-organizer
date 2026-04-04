# downloads-organizer

A simple Go automation program to sort your downloaded files. It uses [rofi](https://github.com/davatorium/rofi) and is developed for Arch Linux.

When a file lands in your downloads folder, a rofi prompt appears asking you to pick a category. The file is then moved automatically.

## Requirements

- Go 1.26+
- `rofi`

## Installation

### 1. Build the binary

```bash
git clone https://github.com/javamilesiii/downloads-organizer
cd downloads-organizer
go build -o downloads-organizer
```

### 2. Install the binary

```bash
install -Dm755 downloads-organizer ~/.local/bin/downloads-organizer
```

### 3. Configure

Copy the example config to your config directory:

```bash
mkdir -p ~/.config/downloads-organizer
cp config.yaml ~/.config/downloads-organizer/config.yaml
```

Edit `~/.config/downloads-organizer/config.yaml` to match your directory layout:

```yaml
download_dir: Downloads/

categories:
  Documents:
    path: Documents/
    prompt_subfolder: true
  Music:
    path: Music/
    prompt_subfolder: false

ignore_files:
  - .crdownload
  - .part
  - .tmp
  - Unconfirmed
```

- `download_dir` — path relative to your home directory to watch
- `categories` — map of category names to target directories (relative to home)
- `prompt_subfolder` — if `true`, a second rofi prompt asks for a subdirectory
- `ignore_files` — entries starting with `.` are matched as suffixes; others as prefixes

### 4. Set up the systemd user service

Create `~/.config/systemd/user/downloads-organizer.service`:

```ini
[Unit]
Description=Downloads Organizer
After=graphical-session.target

[Service]
Type=simple
ExecStart=/home/username/.local/bin/downloads-organizer
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=default.target
```

User services inherit your session environment automatically, so no need to set `HOME`, `DISPLAY`, or `WAYLAND_DISPLAY` manually.

### 5. Enable and start the service

```bash
systemctl --user daemon-reload
systemctl --user enable --now downloads-organizer.service
```

### 6. Verify it is running

```bash
systemctl --user status downloads-organizer.service
```

## Stopping / restarting

```bash
systemctl --user stop downloads-organizer.service
systemctl --user restart downloads-organizer.service
```

## Logs

```bash
journalctl --user -u downloads-organizer.service -f
```
