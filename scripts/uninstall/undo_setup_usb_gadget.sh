#!/bin/bash
set -e

if [[ -f /boot/firmware/config.txt ]]; then
  CONFIG_TXT=/boot/firmware/config.txt
else
  CONFIG_TXT=/boot/config.txt
fi

if [[ -f "${CONFIG_TXT}.bak" ]]; then
  sudo cp "${CONFIG_TXT}.bak" "${CONFIG_TXT}"
fi

sudo sed -i '/dtoverlay=dwc2/d' "${CONFIG_TXT}"
