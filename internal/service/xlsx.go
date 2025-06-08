package service

import (
	"fmt"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"net/http"
	"os"
	"path/filepath"
)

func (h *PillsHandler) ExportPillsXLSX(w http.ResponseWriter, r *http.Request) {
	dirPath := filepath.Join("result")

	latestFile, err := utils.FindLatestParsingFile(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Export file not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed finding export file", http.StatusInternalServerError)
		}
		return
	}

	file, err := os.Open(latestFile)
	if err != nil {
		http.Error(w, "Failed to open export file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	fileName := filepath.Base(latestFile)
	if err != nil {
		http.Error(w, "Failed to load file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
}
