# OPlus OTA Finder

An interactive terminal app to find and retrieve OTA download links for **OPPO / OnePlus / Realme** devices directly from official OPlus servers.

Built with [huh](https://github.com/charmbracelet/huh) for a beautiful TUI experience.

> Powered by [Houvven/OplusUpdater](https://pkg.go.dev/github.com/Houvven/OplusUpdater)

---

## ✨ Features

- Beautiful interactive forms with keyboard navigation
- Scrollable region selector — no need to type codes
- Inline input validation
- Shows Android version, OS version, security patch, file size
- Saves download URL automatically to `ota_MODEL_REGION.txt`
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
pkg install golang git -y
```

### 2. Install OplusUpdater backend
```bash
go install github.com/Houvven/OplusUpdater/cmd/updater@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3. Clone and build
```bash
git clone https://github.com/aliakbar6821/Oplus_ota ~/ota-go
cd ~/ota-go
go mod tidy
go build -o oplus-ota .
```

### 4. Add alias
```bash
echo 'alias updater-go="~/ota-go/oplus-ota"' >> ~/.bashrc
source ~/.bashrc
```

Now just type `updater-go` anywhere in Termux to launch.

---

## 📖 Usage

```bash
updater-go
```

### Steps
1. Enter device model (e.g. `PLJ110`, `RMX3820`)
2. Enter OTA version (e.g. `PLJ110_11.A`)
3. Select region from scrollable list
4. Confirm or enter custom NV ID
5. Select mode: Stable or Beta

### After finding update
- Download URL printed and saved to `ota_MODEL_REGION.txt`
- Option to print clean URL for manual copy
- Option to show changelog URL

---

## 🌍 Supported Regions

| Code | Region |
|------|--------|
| CN | China |
| EU | Europe |
| IN | India |
| EX | Export |
| RU | Russia |
| TW | Taiwan |
| ID | Indonesia |
| MY | Malaysia |
| TH | Thailand |
| VN | Vietnam |
| PH | Philippines |
| SG | Singapore |
| TR | Turkey |
| SA | Saudi Arabia |
| BR | Brazil |
| MX | Mexico |
| EG | Egypt |

---

## 📄 Credits

- **[Houvven/OplusUpdater](https://pkg.go.dev/github.com/Houvven/OplusUpdater)** — OPlus OTA query backend
- **[charmbracelet/huh](https://github.com/charmbracelet/huh)** — TUI forms library
- **[charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)** — Terminal styling

---

## ⚠️ Disclaimer

This tool uses official OPlus OTA servers for informational purposes only.
Download links are temporary and expire after a short time.
