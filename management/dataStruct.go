package management

import (
	"BRGS/models"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

var EXCEL_HEAD_TRANSLATE_DIC = map[string]string{
	"name":            "name",
	"watchDir":        "watchDir",
	"tempDir":         "tempDir",
	"archiveDir":      "archiveDir",
	"archiveInterval": "archiveInterval",
	"syncInterval":    "syncInterval",
}

var EXCEL_HEAD_ORDER = []string{
	"name",
	"watchDir",
	"tempDir",
	"archiveDir",
	"archiveInterval",
	"syncInterval",
}

type BackupArchive struct {
	root            models.FSTreeRoot
	archiveDir      string
	name            string
	tempDir         string
	watchDir        string
	archiveInterval int
	syncInterval    int
}

func init() {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	for k := range EXCEL_HEAD_TRANSLATE_DIC {
		EXCEL_HEAD_TRANSLATE_DIC[k] = cfg.Section("excel_head_ch").Key(k).String()
	}

}
