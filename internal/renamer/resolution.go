package renamer

import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

var videoExtensions = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".webm": true,
	".m4v":  true,
	".3gp":  true,
	".flv":  true,
	".wmv":  true,
	".mpg":  true,
	".mpeg": true,
	".m2v":  true,
	".mts":  true,
	".m2ts": true,
}

// Photos supported by Go standard library (fast, no ffprobe needed)
var goSupportedPhotos = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Photos that require ffprobe (not supported by Go standard library)
var ffprobeOnlyPhotos = map[string]bool{
	".heic": true,
	".heif": true,
	".webp": true,
}

func isHasFFProbe() bool {
	_, err := exec.Command("ffprobe", "-version").Output()
	return err == nil
}

func isMediaFile(file os.DirEntry) bool {
	ext := strings.ToLower(filepath.Ext(file.Name()))
	return goSupportedPhotos[ext] || ffprobeOnlyPhotos[ext] || videoExtensions[ext]
}

func isGoSupportedPhoto(ext string) bool {
	return goSupportedPhotos[strings.ToLower(ext)]
}

// getImageResolution reads image dimensions using Go standard library (very fast, only reads header)
func getImageResolution(filePath string) (int, int) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}

	return config.Width, config.Height
}

// getResolutionFFProbe uses ffprobe for videos and unsupported image formats
func getResolutionFFProbe(filePath string) (int, int) {
	output, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filePath).Output()
	if err != nil {
		return 0, 0
	}

	var ffProbeOutput ffprobeOutput
	if err := json.Unmarshal(output, &ffProbeOutput); err != nil {
		return 0, 0
	}

	if len(ffProbeOutput.Streams) == 0 {
		return 0, 0
	}

	return ffProbeOutput.Streams[0].Width, ffProbeOutput.Streams[0].Height
}

func getResolution(file os.DirEntry, path string) (int, int) {
	if !isMediaFile(file) {
		return 0, 0
	}

	filePath := filepath.Join(path, file.Name())
	ext := filepath.Ext(file.Name())

	// Use Go standard library for supported images (much faster)
	if isGoSupportedPhoto(ext) {
		return getImageResolution(filePath)
	}

	// Use ffprobe for videos and HEIC/HEIF
	return getResolutionFFProbe(filePath)
}
