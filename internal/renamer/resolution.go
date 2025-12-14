package renamer

import (
	"encoding/json"
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

var photoExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".heic": true,
	".webp": true,
	".heif": true,
}

func isHasFFProbe() bool {
	_, err := exec.Command("ffprobe", "-version").Output()
	return err == nil
}

func isMediaFile(file os.DirEntry) bool {
	ext := strings.ToLower(filepath.Ext(file.Name()))
	return photoExtensions[ext] || videoExtensions[ext]
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func getRatio(num, den int) (int, int) {
	if den == 0 {
		return 0, 0
	}
	g := gcd(num, den)
	return num / g, den / g
}

func getResolution(file os.DirEntry, path string) (int, int) {
	if !isMediaFile(file) {
		return 0, 0
	}

	filePath := filepath.Join(path, file.Name())
	output, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filePath).Output()
	if err != nil {
		return 0, 0
	}

	var ffProbeOutput ffprobeOutput
	err = json.Unmarshal(output, &ffProbeOutput)
	if err != nil {
		return 0, 0
	}

	if len(ffProbeOutput.Streams) == 0 {
		return 0, 0
	}

	return ffProbeOutput.Streams[0].Width, ffProbeOutput.Streams[0].Height
}
