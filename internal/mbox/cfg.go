package mbox

import (
	"os"
	"path/filepath"

	"github.com/la5nta/pat/app"
	"github.com/la5nta/pat/cfg"
	"github.com/la5nta/wl2k-go/mailbox"
)

func GetDefaultInbox() string {
	configDir, err := patConfigDir()
	if err != nil {
		return ""
	}

	mboxPath, err := patMailboxDir()
	if err != nil {
		return ""
	}

	patCfg, err := app.LoadConfig(filepath.Join(configDir, "config.json"), cfg.DefaultConfig)
	if err != nil {
		return ""
	}

	if len(patCfg.MyCall) < 1 {
		return ""
	}

	return filepath.Join(mboxPath, patCfg.MyCall, mailbox.DIR_INBOX)
}

func patMailboxDir() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")

	if base == "" || !filepath.IsAbs(base) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		base = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(base, "pat", "mailbox"), nil
}

func patConfigDir() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")

	if base == "" || !filepath.IsAbs(base) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		base = filepath.Join(home, ".config")
	}

	return filepath.Join(base, "pat"), nil
}
