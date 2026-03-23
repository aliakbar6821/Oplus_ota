package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ─── Styles ───────────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("212")).
			Padding(0, 2).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("82"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Underline(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 2).
			MarginTop(1).
			MarginBottom(1)
)

// ─── Regions ──────────────────────────────────────────────────────────────────

type Region struct {
	Name  string
	NvID  string
	Server string
}

var regions = map[string]Region{
	"CN": {Name: "China",        NvID: "10010111", Server: "CN"},
	"EU": {Name: "Europe",       NvID: "01000100", Server: "EU"},
	"IN": {Name: "India",        NvID: "00011011", Server: "IN"},
	"EX": {Name: "Export",       NvID: "00000000", Server: "EU"},
	"RU": {Name: "Russia",       NvID: "00110111", Server: "EU"},
	"TW": {Name: "Taiwan",       NvID: "00011010", Server: "EU"},
	"ID": {Name: "Indonesia",    NvID: "00110011", Server: "EU"},
	"MY": {Name: "Malaysia",     NvID: "00111000", Server: "EU"},
	"TH": {Name: "Thailand",     NvID: "00111001", Server: "EU"},
	"VN": {Name: "Vietnam",      NvID: "00111100", Server: "EU"},
	"PH": {Name: "Philippines",  NvID: "00111110", Server: "EU"},
	"SG": {Name: "Singapore",    NvID: "00101100", Server: "EU"},
	"TR": {Name: "Turkey",       NvID: "01010001", Server: "EU"},
	"SA": {Name: "Saudi Arabia", NvID: "10000011", Server: "EU"},
	"BR": {Name: "Brazil",       NvID: "10011110", Server: "EU"},
	"MX": {Name: "Mexico",       NvID: "01111011", Server: "EU"},
	"EG": {Name: "Egypt",        NvID: "01110101", Server: "EU"},
}

// ─── OTA Response ─────────────────────────────────────────────────────────────

type OTAResponse struct {
	Body struct {
		RealOtaVersion  string `json:"realOtaVersion"`
		RealVersionName string `json:"realVersionName"`
		RealAndroidVer  string `json:"realAndroidVersion"`
		RealOsVersion   string `json:"realOsVersion"`
		SecurityPatch   string `json:"securityPatch"`
		Components []struct {
			ComponentPackets struct {
				URL       string `json:"url"`
				ManualURL string `json:"manualUrl"`
				Size      string `json:"size"`
			} `json:"componentPackets"`
		} `json:"components"`
		Description struct {
			PanelURL string `json:"panelUrl"`
		} `json:"description"`
	} `json:"body"`
	ErrMsg       string `json:"errMsg"`
	ResponseCode int    `json:"responseCode"`
}

// ─── Query OTA ────────────────────────────────────────────────────────────────

func queryOTA(model, otaVer, nvID, region string, mode int) (*OTAResponse, error) {
	args := []string{
		otaVer,
		"--model", model,
		"--carrier", nvID,
		"--region", region,
		"--mode", strconv.Itoa(mode),
	}

	cmd := exec.Command("updater", args...)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return nil, fmt.Errorf("updater command failed: %v", err)
	}

	// Strip ANSI codes
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	clean := ansi.ReplaceAllString(string(out), "")

	var resp OTAResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(clean)), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &resp, nil
}

// ─── Save URL ─────────────────────────────────────────────────────────────────

func saveURL(model, region, otaVer, url string) string {
	filename := fmt.Sprintf("ota_%s_%s.txt", model, region)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ""
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "[%s]\nOTA: %s\nURL: %s\n\n", timestamp, otaVer, url)
	return filename
}

// ─── Show Result ──────────────────────────────────────────────────────────────

func showResult(resp *OTAResponse, model, region string) {
	body := resp.Body
	url := ""
	manualURL := ""
	size := ""
	if len(body.Components) > 0 {
		url = body.Components[0].ComponentPackets.URL
		manualURL = body.Components[0].ComponentPackets.ManualURL
		size = body.Components[0].ComponentPackets.Size
	}

	// Convert size
	sizeDisplay := ""
	if size != "" {
		if bytes, err := strconv.ParseInt(size, 10, 64); err == nil {
			sizeDisplay = fmt.Sprintf("%.2f GB", float64(bytes)/1073741824)
		}
	}

	// Result box
	result := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		labelStyle.Render("📱 Version: "), valueStyle.Render(body.RealVersionName),
		labelStyle.Render("🤖 Android: "), valueStyle.Render(body.RealAndroidVer),
		labelStyle.Render("🎨 OS:      "), valueStyle.Render(body.RealOsVersion),
		labelStyle.Render("🔒 Patch:   "), valueStyle.Render(body.SecurityPatch),
		labelStyle.Render("📦 OTA:     "), valueStyle.Render(body.RealOtaVersion),
		labelStyle.Render("💾 Size:    "), valueStyle.Render(sizeDisplay),
	)

	fmt.Println(successStyle.Render("✅ Update found!"))
	fmt.Println(boxStyle.Render(result))

	if url != "" {
		fmt.Println(labelStyle.Render("📥 Download URL:"))
		fmt.Println(urlStyle.Render(url))
		fmt.Println()
		if manualURL != "" && manualURL != url {
			fmt.Println(labelStyle.Render("📥 Manual URL (try if slow):"))
			fmt.Println(urlStyle.Render(manualURL))
			fmt.Println()
		} else if manualURL != "" {
			fmt.Println(warnStyle.Render("ℹ️  Manual URL is same as download URL"))
			fmt.Println()
		}

		// Save both URLs to file
		filename := saveURL(model, region, body.RealOtaVersion, url)
		if filename != "" {
			// Also save manualURL
			if manualURL != "" {
				f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
				if f != nil {
					fmt.Fprintf(f, "ManualURL: %s\n\n", manualURL)
					f.Close()
				}
			}
			fmt.Println(successStyle.Render("✅ Saved to " + filename))
		}

		// Post options
		var action string
		huh.NewSelect[string]().
			Title("What do you want to do?").
			Options(
				huh.NewOption("Print Download URL", "print"),
				huh.NewOption("Print Manual URL (try if slow)", "printmanual"),
				huh.NewOption("Show changelog URL", "changelog"),
				huh.NewOption("Continue", "skip"),
			).
			Value(&action).
			Run()

		switch action {
		case "print":
			fmt.Println()
			fmt.Println(labelStyle.Render("=== DOWNLOAD URL ==="))
			fmt.Println(url)
			fmt.Println(labelStyle.Render("===================="))
		case "printmanual":
			fmt.Println()
			fmt.Println(labelStyle.Render("=== MANUAL URL ==="))
			fmt.Println(manualURL)
			fmt.Println(labelStyle.Render("=================="))
		case "changelog":
			if body.Description.PanelURL != "" {
				fmt.Println(urlStyle.Render(body.Description.PanelURL))
			} else {
				fmt.Println(warnStyle.Render("No changelog URL available."))
			}
		}
	} else {
		fmt.Println(errorStyle.Render("❌ No download URL in response."))
	}
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	fmt.Println(titleStyle.Render("  OPlus OTA Finder  "))

	for {
		// ── Step 1: Model input ──
		var model string
		if err := huh.NewInput().
			Title("Device Model").
			Placeholder("e.g. PLJ110, RMX3820").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("model cannot be empty")
				}
				return nil
			}).
			Value(&model).
			Run(); err != nil {
			fmt.Println(errorStyle.Render("Cancelled."))
			return
		}
		model = strings.TrimSpace(strings.ToUpper(model))

		// ── Step 2: OTA version ──
		var otaVer string
		if err := huh.NewInput().
			Title("OTA Version").
			Placeholder("e.g. PLJ110_11.A").
			Description("Format: MODEL_11.X").
			Validate(func(s string) error {
				matched, _ := regexp.MatchString(`^[A-Z0-9]+_11\.[A-Z]`, strings.ToUpper(s))
				if !matched {
					return fmt.Errorf("invalid format, use MODEL_11.X (e.g. PLJ110_11.A)")
				}
				return nil
			}).
			Value(&otaVer).
			Run(); err != nil {
			return
		}
		otaVer = strings.TrimSpace(strings.ToUpper(otaVer))

		// ── Step 3: Region select ──
		regionOptions := []huh.Option[string]{}
		for code, r := range regions {
			regionOptions = append(regionOptions, huh.NewOption(
				fmt.Sprintf("%s — %s", code, r.Name), code,
			))
		}

		var regionCode string
		if err := huh.NewSelect[string]().
			Title("Select Region").
			Options(regionOptions...).
			Value(&regionCode).
			Run(); err != nil {
			return
		}

		region := regions[regionCode]
		nvID := region.NvID

		// ── Step 4: Custom NV ID? ──
		var useCustomNV bool
		huh.NewConfirm().
			Title(fmt.Sprintf("NV ID: %s — use custom?", nvID)).
			Value(&useCustomNV).
			Run()

		if useCustomNV {
			huh.NewInput().
				Title("Custom NV ID").
				Value(&nvID).
				Run()
		}

		// ── Step 5: Mode ──
		var modeStr string
		huh.NewSelect[string]().
			Title("Update Mode").
			Options(
				huh.NewOption("Stable", "0"),
				huh.NewOption("Beta / Testing", "1"),
			).
			Value(&modeStr).
			Run()
		mode, _ := strconv.Atoi(modeStr)

		// ── Query ──
		fmt.Println()
		fmt.Println(warnStyle.Render("🔍 Querying OTA server..."))
		fmt.Printf("%s %s\n", labelStyle.Render("Model:  "), valueStyle.Render(model))
		fmt.Printf("%s %s\n", labelStyle.Render("OTA:    "), valueStyle.Render(otaVer))
		fmt.Printf("%s %s\n", labelStyle.Render("Region: "), valueStyle.Render(regionCode+" — "+region.Name))
		fmt.Printf("%s %s\n", labelStyle.Render("NV ID:  "), valueStyle.Render(nvID))
		fmt.Println()

		resp, err := queryOTA(model, otaVer, nvID, regionCode, mode)
		if err != nil {
			fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))
		} else if resp.ResponseCode != 200 || resp.Body.RealOtaVersion == "" {
			msg := resp.ErrMsg
			if msg == "" {
				msg = "No update found"
			}
			fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Server error (%d): %s", resp.ResponseCode, msg)))
			if resp.ResponseCode == 2004 {
				fmt.Println(warnStyle.Render("💡 Already on latest version or OTA version not found."))
			}
			if resp.ResponseCode == 500 {
				fmt.Println(warnStyle.Render("💡 Wrong OTA version or device not found in this region."))
			}
		} else {
			showResult(resp, model, regionCode)
		}

		// ── Next action ──
		fmt.Println()
		var next string
		if err := huh.NewSelect[string]().
			Title("What next?").
			Options(
				huh.NewOption("Query another version/region (same model)", "same"),
				huh.NewOption("Query different model", "new"),
				huh.NewOption("Exit", "exit"),
			).
			Value(&next).
			Run(); err != nil || next == "exit" {
			fmt.Println(successStyle.Render("👋 Goodbye."))
			return
		}

		if next == "new" {
			fmt.Println(titleStyle.Render("  OPlus OTA Finder  "))
			continue
		}

		// same model — loop with new version/region
		for {
			var newOtaVer string
			if err := huh.NewInput().
				Title("OTA Version").
				Placeholder("e.g. PLJ110_11.A").
				Description("Format: MODEL_11.X").
				Validate(func(s string) error {
					matched, _ := regexp.MatchString(`^[A-Z0-9]+_11\.[A-Z]`, strings.ToUpper(s))
					if !matched {
						return fmt.Errorf("invalid format")
					}
					return nil
				}).
				Value(&newOtaVer).
				Run(); err != nil {
				return
			}

			var newRegionCode string
			huh.NewSelect[string]().
				Title("Select Region").
				Options(regionOptions...).
				Value(&newRegionCode).
				Run()

			newRegion := regions[newRegionCode]
			newNvID := newRegion.NvID

			fmt.Println()
			fmt.Println(warnStyle.Render("🔍 Querying OTA server..."))

			resp, err := queryOTA(model, strings.ToUpper(newOtaVer), newNvID, newRegionCode, 0)
			if err != nil {
				fmt.Println(errorStyle.Render("❌ " + err.Error()))
			} else if resp.ResponseCode != 200 || resp.Body.RealOtaVersion == "" {
				msg := resp.ErrMsg
				if msg == "" { msg = "No update found" }
				fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Server error (%d): %s", resp.ResponseCode, msg)))
			} else {
				showResult(resp, model, newRegionCode)
			}

			var again string
			huh.NewSelect[string]().
				Title("What next?").
				Options(
					huh.NewOption("Query another version/region (same model)", "same"),
					huh.NewOption("Back to model selection", "back"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&again).
				Run()

			if again == "exit" {
				fmt.Println(successStyle.Render("👋 Goodbye."))
				return
			}
			if again == "back" {
				break
			}
		}
	}
}
