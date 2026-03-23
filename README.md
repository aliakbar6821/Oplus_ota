# OPlus OTA Finder

An interactive Termux script to find and retrieve OTA download links for **OPPO / OnePlus / Realme** devices directly from official OPlus servers.

> Powered by [Houvven/OplusUpdater](https://pkg.go.dev/github.com/Houvven/OplusUpdater)

---

## ✨ Features

- Interactive region/model selection with NV identifier table
- Auto NV ID detection per region
- Parses full OTA response — shows Android version, OS version, security patch, file size
- Clean single-line download URL output
- Saves results to `ota_MODEL_REGION.txt`
- Query different versions/regions without restarting

---

## 📋 Requirements

- **Termux** — [Download from F-Droid](https://f-droid.org/packages/com.termux/)
- **Internet connection**

---

## 🚀 Installation

### 1. Update and install dependencies
```bash
apt update && apt upgrade -y
pkg install golang python git -y
```

### 2. Install OplusUpdater
```bash
go install github.com/Houvven/OplusUpdater/cmd/updater@latest
```

### 3. Add Go binaries to PATH
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### 4. Clone this repo
```bash
git clone https://github.com/YOUR_USERNAME/oplus-ota-finder ~/ota
cd ~/ota
chmod +x oplus_ota.sh
```

### 5. Add alias so `updater` runs the script
```bash
echo 'alias updater="bash ~/ota/oplus_ota.sh"' >> ~/.bashrc
source ~/.bashrc
```

Now just type `updater` anywhere in Termux to launch the tool.

---

## 📖 Usage

```bash
updater
```

### Steps
1. Enter device model (e.g. `PLJ110`, `RMX3820`)
2. Enter OTA version (e.g. `PLJ110_11.A`)
3. Enter region code from the table (e.g. `CN`, `EU`, `IN`)
4. Press Enter to use auto NV ID or enter custom
5. Select mode: `0` stable, `1` beta

### After finding update
- Download URL is printed and saved automatically to `ota_MODEL_REGION.txt`
- Option `1` — Print clean URL for manual copy
- Option `2` — Show changelog URL
- Option `3` — Continue to next query

---

## 🌍 Supported Regions

| Code | Region | NV Identifier |
|------|--------|---------------|
| CN | China | 10010111 |
| EU | Europe | 01000100 |
| IN | India | 00011011 |
| EX | Export | 00000000 |
| RU | Russia | 00110111 |
| TW | Taiwan | 00011010 |
| ID | Indonesia | 00110011 |
| MY | Malaysia | 00111000 |
| TH | Thailand | 00111001 |
| VN | Vietnam | 00111100 |
| PH | Philippines | 00111110 |
| SG | Singapore | 00101100 |
| TR | Turkey | 01010001 |
| SA | Saudi Arabia | 10000011 |
| BR | Brazil | 10011110 |
| MX | Mexico | 01111011 |
| EG | Egypt | 01110101 |

---

## 📄 Credits

- **[Houvven/OplusUpdater](https://pkg.go.dev/github.com/Houvven/OplusUpdater)** — OPlus OTA query backend that powers this tool
- **[XDA - OPPO OTA Downloader](https://xdaforums.com/t/oppo-ota-downloader.4728411/)** — Interface design inspiration

---

## ⚠️ Disclaimer

This tool uses official OPlus OTA servers for informational purposes only.
Download links are temporary and expire after a short time.
