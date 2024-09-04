package main

import (
	"fmt"
	"os/exec"
)

func Alert(title, message, appIcon string) error {
	cmd := exec.Command("kdialog", "--title", title, "--passivepopup", message, "30", "--icon", appIcon)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to display notification: %v", err)
	}

	beepCmd := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/message.oga")
	if err := beepCmd.Run(); err != nil {
		return fmt.Errorf("failed to play beep sound: %v", err)
	}

	return nil
}
