package livephoto

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dynonguyen/dyno-clis/internal/utils"
)

const (
	cliName               = "livephotos"
	defaultLivephotosName = "Live Photos"
)

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
	".vob":  true,
	".asf":  true,
	".rm":   true,
	".rmvb": true,
	".divx": true,
	".xvid": true,
	".f4v":  true,
	".ogv":  true,
	".ogm":  true,
	".mxf":  true,
	".dv":   true,
	".mod":  true,
	".tod":  true,
	".ts":   true,
}

type cliFlags struct {
	path           string
	keepEmptyDirs  bool
	livephotosName string
	yes            bool
}

// Check if file is a video based on extension
func isVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return videoExtensions[ext]
}

// Ensure directory exists, create if it doesn't
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}
	return nil
}

// Avoid duplicate file names by adding _1, _2,...
func uniqueName(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(path)
	name := path[:len(path)-len(ext)]
	i := 1
	for {
		newPath := fmt.Sprintf("%s_%d%s", name, i, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		i++
	}
}

// Remove empty directories (process from deepest to shallowest)
func removeEmptyDirectories(directories []string, excludeDir string) {
	for i := len(directories) - 1; i >= 0; i-- {
		dir := directories[i]

		// Skip excluded directory (e.g., livephotos folder)
		if dir == excludeDir {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", dir, err)
			continue
		}

		// If directory is empty, remove it
		if len(entries) == 0 {
			if err := os.Remove(dir); err != nil {
				fmt.Printf("Error removing empty directory %s: %v\n", dir, err)
			}
		}
	}
}

func parseFlags() *cliFlags {
	flags := &cliFlags{
		path:          "",
		keepEmptyDirs: false,
		yes:           false,
	}

	flagItems := []utils.FlagItem{
		{
			Name:   "path",
			Desc:   "Path to the directory to clean up (default: current directory)",
			Flags:  []string{"p", "path"},
			StrVal: &flags.path,
		},
		{
			Name:    "Keep empty directories",
			Desc:    "Keep empty directories after moving files (default: false)",
			Flags:   []string{"k", "keepEmptyDirs"},
			BoolVal: &flags.keepEmptyDirs,
		},
		{
			Name:       "Live Photos Name",
			Desc:       "Name of the live photos directory (default: " + defaultLivephotosName + ")",
			Flags:      []string{"n", "livephotosName"},
			StrVal:     &flags.livephotosName,
			DefaultVal: defaultLivephotosName,
		},
		{
			Name:    "Yes",
			Desc:    "Skip confirmation prompt (default: false)",
			Flags:   []string{"y", "yes"},
			BoolVal: &flags.yes,
		},
	}

	utils.ParseFlags(flagItems, cliName+" -p /path/to/directory")

	return flags
}

func confirmAction(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	if message != "" {
		fmt.Print(message)
	}
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	// Default to yes if empty (just pressed Enter)
	if response == "" {
		return true
	}
	return response == "y" || response == "yes"
}

func Execute() {
	flags := parseFlags()
	rootPath := flags.path

	// Get the current directory if no path is provided
	if rootPath == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current directory", err)
			return
		}
		rootPath = currentPath
	}

	livephotosDir := filepath.Join(rootPath, flags.livephotosName)
	cleanupDirs := []string{}

	// First pass: collect information about what will be moved
	videoCount := 0
	otherFileCount := 0
	dirCount := 0

	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == rootPath {
			return nil
		}

		if path == livephotosDir {
			return nil
		}

		if info.IsDir() {
			dirCount++
			return nil
		}

		if isVideoFile(info.Name()) {
			videoCount++
		} else {
			otherFileCount++
		}

		return nil
	}); err != nil {
		log.Fatal("Failed to scan directory:", err)
	}

	// Show summary and ask for confirmation
	if !flags.yes {
		fmt.Printf("\n=== Summary ===\n")
		fmt.Printf("Root path: %s\n", rootPath)
		fmt.Printf("Live Photos directory: %s\n", livephotosDir)
		fmt.Printf("Video files to move: %d\n", videoCount)
		fmt.Printf("Other files to move: %d\n", otherFileCount)
		if !flags.keepEmptyDirs {
			fmt.Printf("Empty directories to remove: %d\n", dirCount)
		}
		fmt.Printf("\n⚠️  WARNING: This will move files and may remove empty directories!\n")
		if !confirmAction("Do you want to continue? (Y/n): ") {
			fmt.Println("Operation cancelled.")
			return
		}
		fmt.Println()
	}

	// Second pass: perform the actual operations
	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			return err
		}

		// Skip the root directory
		if path == rootPath {
			return nil
		}

		// Skip the livephotos directory itself to avoid moving it into itself
		if path == livephotosDir {
			return nil
		}

		// If it's a directory, collect it for later cleanup
		if info.IsDir() {
			if !flags.keepEmptyDirs {
				cleanupDirs = append(cleanupDirs, path)
			}
			return nil
		}

		// Determine destination based on file type
		var dest string
		if isVideoFile(info.Name()) {
			// Ensure livephotos directory exists
			if err := ensureDir(livephotosDir); err != nil {
				return err
			}
			dest = filepath.Join(livephotosDir, info.Name())
		} else {
			// Move non-video files to root
			dest = filepath.Join(rootPath, info.Name())
		}

		dest = uniqueName(dest)

		if err := os.Rename(path, dest); err != nil {
			return fmt.Errorf("rename %s -> %s: %w", path, dest, err)
		}

		fmt.Printf("Moved %s -> %s\n", path, dest)

		return nil
	}); err != nil {
		log.Fatal("Failed to walk through the directory", err)
	}

	// Remove empty directories (process from deepest to shallowest)
	if !flags.keepEmptyDirs && len(cleanupDirs) > 0 {
		fmt.Printf("Removing empty directories (%d) after moving files...\n", len(cleanupDirs))
		removeEmptyDirectories(cleanupDirs, livephotosDir)
	}
}
