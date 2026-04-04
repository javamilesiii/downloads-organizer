# downloads-organizer

A simple Go automation program to sort your downloaded files. It uses [rofi](https://github.com/davatorium/rofi) and is developed for Arch Linux.

When a file lands in your downloads folder, a rofi prompt appears asking you to pick a category. The file is then moved automatically.

## Requirements

- Go 1.26+
- `rofi`

## Installation

### 1. Build the binary

```bash
git clone (https://github.com/javamilesiii/downloads-organizer)
cd downloads-organizer
go build -o downloads-organizer
```

### 2. Install the binary

```bashconfig
sudo install -Dm755 downloads-organizer /usr/local/bin/downloads-organizer
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

### 4. Set up the systemd service

Paste the following into `/etc/systemd/system/downloads-organizer.service`:

```ini
[Unit]
Description=Downloads Organizer
After=graphical-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/downloads-organizer
Restart=on-failure
RestartSec=5s
Environment=HOME=/home/your-username
Environment=DISPLAY=:0
Environment=WAYLAND_DISPLAY=wayland-0
Environment=XDG_RUNTIME_DIR=/run/user/1000

[Install]
WantedBy=graphical-session.target
```

> Replace `your-username` with your actual username and `1000` with your actual user ID (check with `id -u`). All four environment variables are required: `HOME` for config lookup, and the rest so rofi can open on your screen.

### 5. Enable and start the service

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now downloads-organizer.service
```

### 6. Verify it is running

```bash
systemctl status downloads-organizer.service
```

## Stopping / restarting

```bash
sudo systemctl stop downloads-organizer.service
sudo systemctl restart downloads-organizer.service
```

## Logs

```bash
journalctl -u downloads-organizer.service -f
```
