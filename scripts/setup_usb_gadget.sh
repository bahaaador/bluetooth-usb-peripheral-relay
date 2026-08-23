#!/bin/bash
set -e

# Raspberry Pi OS Bookworm+ and DietPi use /boot/firmware/config.txt; older images use /boot/config.txt
if [[ -f /boot/firmware/config.txt ]]; then
  CONFIG_TXT=/boot/firmware/config.txt
else
  CONFIG_TXT=/boot/config.txt
fi

echo "Using boot config: ${CONFIG_TXT}"

sudo cp "${CONFIG_TXT}" "${CONFIG_TXT}.bak"

append_line_if_absent() {
  local line="$1"
  if ! sudo grep -Fxq "$line" "${CONFIG_TXT}" 2>/dev/null; then
    echo "$line" | sudo tee -a "${CONFIG_TXT}" > /dev/null
  fi
}

append_line_if_absent "dtoverlay=dwc2"
append_line_if_absent "dtoverlay=dwc2,dr_mode=peripheral"

sudo cp /etc/modules /etc/modules.bak

append_module_if_absent() {
  local mod="$1"
  if [[ -f /etc/modules ]] && grep -Fix "^${mod}$" /etc/modules >/dev/null 2>&1; then
    return 0
  fi
  echo "$mod" | sudo tee -a /etc/modules > /dev/null
}

append_module_if_absent "dwc2"
append_module_if_absent "libcomposite"
 
echo "USB gadget (peripheral) mode enabled: updated ${CONFIG_TXT} and /etc/modules. Reboot for changes to take effect."
