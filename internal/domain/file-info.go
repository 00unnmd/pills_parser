package domain

import "time"

type FileInfo struct {
	Path    string
	ModTime time.Time
	Date    time.Time
}
