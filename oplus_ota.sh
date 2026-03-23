#!/bin/bash

# Colors
RED="\e[31m"; GREEN="\e[32m"; PURPLE="\e[35m"
YELLOW="\e[33m"; BLUE="\e[34m"; RESET="\e[0m"

# Check updater is installed
if ! command -v updater &>/dev/null; then
    echo -e "${RED}❌ updater not found. Installing...${RESET}"
    export PATH=$PATH:$(go env GOPATH)/bin
    if ! command -v updater &>/dev/null; then
        echo -e "${RED}❌ Please run: go install github.com/Houvven/OplusUpdater/cmd/updater@latest${RESET}"
        exit 1
    fi
fi

export PATH=$PATH:$(go env GOPATH)/bin

# Regions
declare -A REGIONS=(
    [CN]="China 10010111"
    [EU]="Europe 01000100"
    [IN]="India 00011011"
    [EX]="Export 00000000"
    [RU]="Russia 00110111"
    [TW]="Taiwan 00011010"
    [ID]="Indonesia 00110011"
    [MY]="Malaysia 00111000"
    [TH]="Thailand 00111001"
    [VN]="Vietnam 00111100"
    [PH]="Philippines 00111110"
    [SG]="Singapore 00101100"
    [TR]="Turkey 01010001"
    [SA]="Saudi_Arabia 10000011"
    [BR]="Brazil 10011110"
    [MX]="Mexico 01111011"
    [EG]="Egypt 01110101"
)

# Print banner
print_banner() {
    clear
    echo -e "${GREEN}+============================================+${RESET}"
    echo -e "${GREEN}|==  ${RESET}    OPlus OTA Finder (OplusUpdater)    ${GREEN}==|${RESET}"
    echo -e "${GREEN}+============================================+${RESET}"
    echo ""
}

# Print regions table
print_regions() {
    echo -e "${YELLOW}Available Regions:${RESET}"
    echo -e "+--------+----------------------+------------------+"
    printf "| %-6s | %-20s | %-16s |\n" "Code" "Region" "NV Identifier"
    echo -e "+--------+----------------------+------------------+"
    for key in $(echo "${!REGIONS[@]}" | tr ' ' '\n' | sort); do
        data=(${REGIONS[$key]})
        printf "| ${YELLOW}%-6s${RESET} | %-20s | %-16s |\n" "$key" "${data[0]}" "${data[1]}"
    done
    echo -e "+--------+----------------------+------------------+"
    echo ""
}

# Run updater query
run_query() {
    local model="$1"
    local ota_ver="$2"
    local nv_id="$3"
    local region="$4"
    local mode="${5:-0}"

    echo -e "\n${BLUE}🔍 Querying OTA server...${RESET}"
    echo -e "🛠  Model:   ${GREEN}$model${RESET}"
    echo -e "🛠  Version: ${GREEN}$ota_ver${RESET}"
    echo -e "🛠  Region:  ${GREEN}$region${RESET}"
    echo -e "🛠  NV ID:   ${GREEN}$nv_id${RESET}"
    echo ""

    output=$(updater "$ota_ver" \
        --model "$model" \
        --carrier "$nv_id" \
        --region "$region" \
        --mode "$mode" 2>&1)

    # Check for error - success if body has androidVersion
    if echo "$output" | grep -q '"androidVersion"'; then
        # Use Python to parse JSON properly
        parsed=$(echo "$output" | python3 -c "
import sys, json, re
try:
    raw = sys.stdin.read()
    clean = re.sub(r'\x1b\[[0-9;]*m', '', raw)
    data = json.loads(clean)
    body = data.get('body', {})
    components = body.get('components', [])
    dl_url = ''
    size = ''
    if components:
        packets = components[0].get('componentPackets', {})
        dl_url = packets.get('url', '')
        size = packets.get('size', '')
    desc = body.get('description', {})
    panel_url = desc.get('panelUrl', '')
    size_gb = round(int(size) / 1073741824, 2) if size else 0
    print('REAL_OTA=' + body.get('realOtaVersion', ''))
    print('VERSION_NAME=' + body.get('realVersionName', ''))
    print('ANDROID_VER=' + body.get('realAndroidVersion', ''))
    print('OS_VER=' + body.get('realOsVersion', ''))
    print('SECURITY_PATCH=' + body.get('securityPatch', ''))
    print('SIZE=' + str(size_gb) + ' GB')
    manual_url = ''
    if components:
        manual_url = components[0].get('componentPackets', {}).get('manualUrl', '')
    print('PANEL_URL=' + panel_url)
    print('DL_URL=' + dl_url)
    print('MANUAL_URL=' + manual_url)
except Exception as e:
    print('ERROR=' + str(e))
")

        # Load parsed values
        real_ota=$(echo "$parsed" | grep "^REAL_OTA=" | cut -d= -f2-)
        version_name=$(echo "$parsed" | grep "^VERSION_NAME=" | cut -d= -f2-)
        android_ver=$(echo "$parsed" | grep "^ANDROID_VER=" | cut -d= -f2-)
        os_ver=$(echo "$parsed" | grep "^OS_VER=" | cut -d= -f2-)
        security_patch=$(echo "$parsed" | grep "^SECURITY_PATCH=" | cut -d= -f2-)
        size_display=$(echo "$parsed" | grep "^SIZE=" | cut -d= -f2-)
        panel_url=$(echo "$parsed" | grep "^PANEL_URL=" | cut -d= -f2-)
        dl_url=$(echo "$parsed" | grep "^DL_URL=" | cut -d= -f2-)
        manual_url=$(echo "$parsed" | grep "^MANUAL_URL=" | cut -d= -f2-)

        echo -e "${GREEN}✅ Update found!${RESET}"
        echo -e "+============================================+"
        echo -e "  📱 ${YELLOW}$version_name${RESET}"
        echo -e "  🤖 Android:  ${GREEN}$android_ver${RESET}"
        echo -e "  🎨 OS:       ${GREEN}$os_ver${RESET}"
        echo -e "  🔒 Patch:    ${GREEN}$security_patch${RESET}"
        echo -e "  📦 OTA:      ${BLUE}$real_ota${RESET}"
        echo -e "  💾 Size:     ${YELLOW}$size_display${RESET}"
        echo -e "+============================================+"
        echo ""

        if [[ -n "$dl_url" ]]; then
            echo -e "📥 ${GREEN}Download URL:${RESET}"
            echo -e "${BLUE}$dl_url${RESET}"
            echo ""

            # Clean URLs
            clean_url=$(echo "$dl_url" | tr -d '\n\r ' )
            clean_manual=$(echo "$manual_url" | tr -d '\n\r ' )

            # Auto save to file
            save_file="ota_${model}_${region}.txt"
            echo "OTA: $real_ota" >> "$save_file"
            echo "URL: $clean_url" >> "$save_file"
            [[ -n "$clean_manual" ]] && echo "ManualURL: $clean_manual" >> "$save_file"
            echo "" >> "$save_file"
            echo -e "${GREEN}✅ Saved to $save_file${RESET}"

            echo ""
            echo -e "Options:"
            echo -e "  ${YELLOW}1${RESET} - Print Download URL"
            echo -e "  ${GREEN}2${RESET} - Print Manual URL (try if slow)"
            echo -e "  ${BLUE}3${RESET} - Show changelog URL"
            echo -e "  ${RED}4${RESET} - Continue"
            read -p "Select: " post_action
            case "$post_action" in
                1)
                    echo ""
                    echo -e "${GREEN}=== DOWNLOAD URL ===${RESET}"
                    echo "$clean_url"
                    echo -e "${GREEN}===================${RESET}"
                    ;;
                2)
                    echo ""
                    echo -e "${GREEN}=== MANUAL URL ===${RESET}"
                    echo "$clean_manual"
                    echo -e "${GREEN}=================${RESET}"
                    ;;
                3)
                    if [[ -n "$panel_url" ]]; then
                        echo -e "${BLUE}$panel_url${RESET}"
                    fi
                    ;;
                4) ;;
            esac
        else
            echo -e "${RED}❌ No download URL found in response.${RESET}"
            echo -e "${YELLOW}Raw output:${RESET}"
            echo "$output"
        fi
    else
        err=$(echo "$output" | grep -o '"errMsg": *"[^"]*"' | cut -d'"' -f4)
        code=$(echo "$output" | grep -o '"responseCode": *[0-9]*' | head -1 | grep -o '[0-9]*$')
        if [[ -n "$err" || -n "$code" ]]; then
            echo -e "${RED}❌ Error: $err (code: $code)${RESET}"
            if [[ "$code" == "2004" ]]; then
                echo -e "${YELLOW}💡 Tip: No update found. Try a different OTA version.${RESET}"
            fi
        else
            resp_code=$(echo "$output" | python3 -c "
import sys,json,re
raw=sys.stdin.read()
clean=re.sub(r'\x1b\[[0-9;]*m','',raw)
try:
    d=json.loads(clean)
    print(d.get('responseCode','?'))
except:
    print('?')
" 2>/dev/null)
            err_msg=$(echo "$output" | python3 -c "
import sys,json,re
raw=sys.stdin.read()
clean=re.sub(r'\x1b\[[0-9;]*m','',raw)
try:
    d=json.loads(clean)
    print(d.get('errMsg','Unknown error'))
except:
    print('Unknown error')
" 2>/dev/null)
            echo -e "${RED}❌ Server error (code: $resp_code): $err_msg${RESET}"
            if [[ "$resp_code" == "500" ]]; then
                echo -e "${YELLOW}💡 Tip: Wrong OTA version or device not found in this region.${RESET}"
            elif [[ "$resp_code" == "2004" ]]; then
                echo -e "${YELLOW}💡 Tip: No update found. You may already be on the latest version.${RESET}"
            fi
        fi
    fi
}

# Main loop
print_banner

while true; do
    print_banner
    print_regions

    # Input model
    read -p "📱 Enter device model (e.g. PLJ110, RMX3820): " model
    if [[ -z "$model" ]]; then
        echo -e "${RED}❌ Model cannot be empty.${RESET}"
        continue
    fi

    # Input OTA version
    echo -e "${YELLOW}💡 OTA version format: MODEL_11.X (e.g. PLJ110_11.A)${RESET}"
    while true; do
        read -p "📌 Enter OTA version: " ota_ver
        if [[ -z "$ota_ver" ]]; then
            echo -e "${RED}❌ OTA version cannot be empty.${RESET}"
            continue
        fi
        if [[ ! "$ota_ver" =~ ^[A-Z0-9]+_11\.[A-Z]$ ]] && [[ ! "$ota_ver" =~ ^[A-Z0-9]+_11\.[A-Z]\. ]]; then
            echo -e "${RED}❌ Invalid format. Use MODEL_11.X (e.g. PLJ110_11.A)${RESET}"
            continue
        fi
        break
    done

    # Input region
    read -p "🌍 Enter region code (e.g. CN, EU, IN): " region
    region=$(echo "$region" | tr '[:lower:]' '[:upper:]')
    if [[ -z "${REGIONS[$region]}" ]]; then
        echo -e "${RED}❌ Invalid region. Use codes from the table above.${RESET}"
        continue
    fi

    # Get NV ID from region
    region_data=(${REGIONS[$region]})
    nv_id="${region_data[1]}"

    # Override NV ID if needed
    echo -e "${YELLOW}💡 Auto NV ID for $region: $nv_id${RESET}"
    read -p "🔢 Press Enter to use auto NV ID or enter custom: " custom_nv
    if [[ -n "$custom_nv" ]]; then
        nv_id="$custom_nv"
    fi

    # Mode
    echo -e "${BLUE}Mode: 0=Stable (default), 1=Testing/Beta${RESET}"
    read -p "⚙️  Mode (0/1) [default 0]: " mode
    mode=${mode:-0}

    run_query "$model" "$ota_ver" "$nv_id" "$region" "$mode"

    # Continue options
    echo ""
    echo -e "🔄 ${YELLOW}1${RESET} - Query another version/region for same model"
    echo -e "🔄 ${GREEN}2${RESET} - Query different model"
    echo -e "❌ ${RED}3${RESET} - Exit"
    read -p "Select: " opt

    case "$opt" in
        1)
            echo -e "${YELLOW}💡 Format: MODEL_11.X (e.g. PLJ110_11.A)${RESET}"
            while true; do
                read -p "📌 Enter new OTA version: " ota_ver
                if [[ -z "$ota_ver" ]]; then
                    echo -e "${RED}❌ OTA version cannot be empty.${RESET}"
                    continue
                fi
                if [[ ! "$ota_ver" =~ ^[A-Z0-9]+_11\.[A-Z]$ ]] && [[ ! "$ota_ver" =~ ^[A-Z0-9]+_11\.[A-Z]\. ]]; then
                    echo -e "${RED}❌ Invalid format. Use MODEL_11.X (e.g. PLJ110_11.A)${RESET}"
                    continue
                fi
                break
            done
            while true; do
                read -p "🌍 Enter region code (e.g. CN, EU, IN): " region
                region=$(echo "$region" | tr '[:lower:]' '[:upper:]')
                if [[ -z "$region" ]]; then
                    echo -e "${RED}❌ Region cannot be empty.${RESET}"
                    continue
                fi
                if [[ -z "${REGIONS[$region]}" ]]; then
                    echo -e "${RED}❌ Invalid region. Try again.${RESET}"
                    continue
                fi
                break
            done
            region_data=(${REGIONS[$region]})
            nv_id="${region_data[1]}"
            echo -e "${YELLOW}💡 Auto NV ID for $region: $nv_id${RESET}"
            read -p "🔢 Press Enter to use auto NV ID or enter custom: " custom_nv
            [[ -n "$custom_nv" ]] && nv_id="$custom_nv"
            read -p "⚙️  Mode (0/1) [default 0]: " mode
            mode=${mode:-0}
            run_query "$model" "$ota_ver" "$nv_id" "$region" "$mode"
            ;;
        2) continue ;;
        3) echo -e "${GREEN}👋 Goodbye.${RESET}"; exit 0 ;;
        *) echo -e "${RED}❌ Invalid option.${RESET}" ;;
    esac
done
