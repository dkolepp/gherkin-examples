package bddcommon

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

// LoadJSON loads JSON test data into v.
//
// A file at that path in the current folder wins, so a journey can bring its
// own data. Otherwise the file comes from bundled, the kit's embedded data/
// folder (like an Ansible role's defaults).
func LoadJSON(bundled fs.FS, fileName string, v any) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		content, err = fs.ReadFile(bundled, "data/"+fileName)
	}
	if err != nil {
		return fmt.Errorf("no test data %q in the current folder or bundled with the kit", fileName)
	}
	if err := json.Unmarshal(content, v); err != nil {
		return fmt.Errorf("test data %q: %w", fileName, err)
	}
	return nil
}
