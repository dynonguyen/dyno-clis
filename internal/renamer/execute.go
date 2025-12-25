package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dynonguyen/dyno-clis/internal/utils"
)

const (
	cliName           = "renamer"
	suffixFlag        = "suffix"
	emptyOverrideFlag = "<empty>"
)

type cliFlags struct {
	yes, allowDir, dryRun bool
	path, prefix, suffix, override, separator,
	include, exclude, detectResolution, createdDate, replace string
	uniqueSuffix bool
}

type replacer struct {
	regex       *regexp.Regexp
	replacement string
}

var defaultFlags = cliFlags{
	path:             "",
	prefix:           "",
	suffix:           "",
	override:         "",
	separator:        "_",
	include:          "",
	exclude:          "",
	createdDate:      "",
	detectResolution: "",
	replace:          "",
	yes:              false,
	allowDir:         false,
	dryRun:           false,
	uniqueSuffix:     false,
}

func parseFlags() *cliFlags {
	flags := defaultFlags

	flagItems := []utils.FlagItem{
		{
			Name:   "path",
			Desc:   "Path to the directory to rename files in, empty to use current directory",
			Flags:  []string{"p", "path"},
			StrVal: &flags.path,
		},
		{
			Name:   "prefix",
			Desc:   "Prefix to add to the file name",
			Flags:  []string{"prefix"},
			StrVal: &flags.prefix,
		},
		{
			Name:   "suffix",
			Desc:   "Suffix to add to the file name",
			Flags:  []string{"suffix"},
			StrVal: &flags.suffix,
		},
		{
			Name:   "override",
			Desc:   "Override the file name with the given name, empty to keep the original name",
			Flags:  []string{"override"},
			StrVal: &flags.override,
		},
		{
			Name:       "separator",
			Desc:       "Separator to use between the prefix, suffix and the original file name",
			Flags:      []string{"separator"},
			DefaultVal: defaultFlags.separator,
			StrVal:     &flags.separator,
		},
		{
			Name:    "allow directories",
			Desc:    "Allow renaming directories",
			Flags:   []string{"allow-dir"},
			BoolVal: &flags.allowDir,
		},
		{
			Name:    "created date",
			Desc:    "Add created date to the file name with the given format",
			Example: "YYYY-MM-DD",
			Flags:   []string{"created-date"},
			StrVal:  &flags.createdDate,
		},
		{
			Name:    "detect resolution",
			Desc:    "Auto detect resolution and add to the file name, only for photo & video files",
			Flags:   []string{"detect-resolution"},
			Example: "prefix or suffix",
			StrVal:  &flags.detectResolution,
		},
		{
			Name:   "include",
			Desc:   "Only rename files that match the given regex",
			Flags:  []string{"include"},
			StrVal: &flags.include,
		},
		{
			Name:   "exclude",
			Desc:   "Exclude files that match the given regex",
			Flags:  []string{"exclude"},
			StrVal: &flags.exclude,
		},
		{
			Name:   "replace",
			Desc:   "Replace the given string or regex with the given replacement, format: old=new",
			Flags:  []string{"replace"},
			StrVal: &flags.replace,
		},
		{
			Name:    "unique suffix",
			Desc:    "Add a unique suffix to the file name to avoid duplicate file names",
			Flags:   []string{"unique-suffix"},
			BoolVal: &flags.uniqueSuffix,
		},
		{
			Name:    "dry run",
			Desc:    "Display the files that will be renamed without actually renaming them",
			Flags:   []string{"dry-run"},
			BoolVal: &flags.dryRun,
		},
		{
			Name:    "yes",
			Desc:    "Skip confirmation prompt and automatically proceed with renaming",
			Flags:   []string{"y", "yes"},
			BoolVal: &flags.yes,
		},
	}

	utils.ParseFlags(flagItems, cliName+" -p /path/to/directory")

	return &flags
}

func getItemInDir(path string) []os.DirEntry {
	items, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Failed to read directory", err)
		return []os.DirEntry{}
	}
	return items
}

func getFileCreatedTime(file os.DirEntry) time.Time {
	info, _ := file.Info()

	stat := info.Sys().(*syscall.Stat_t)
	var sec, nsec int64 = 0, 0

	switch {
	case stat.Birthtimespec.Sec > 0:
		sec, nsec = stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec
	case stat.Ctimespec.Sec > 0:
		sec, nsec = stat.Ctimespec.Sec, stat.Ctimespec.Nsec
	case stat.Atimespec.Sec > 0:
		sec, nsec = stat.Atimespec.Sec, stat.Atimespec.Nsec
	}

	if sec != 0 && nsec != 0 {
		return time.Unix(sec, nsec)
	}

	return info.ModTime()
}

func isHiddenFile(name string) bool {
	return strings.HasPrefix(name, ".")
}

func getRenamedName(f os.DirEntry, opts *cliFlags, replacer *replacer) (ignored bool, newName string) {
	oldName := f.Name()

	// Ignore hidden files
	if isHiddenFile(oldName) {
		return true, ""
	}

	ext := filepath.Ext(oldName)
	nameWoutExt := strings.TrimSuffix(oldName, ext)

	if opts.exclude != "" {
		if match := regexp.MustCompile(opts.exclude).MatchString(oldName); match {
			return true, ""
		}
	}

	if opts.include != "" {
		if match := regexp.MustCompile(opts.include).MatchString(oldName); !match {
			return true, ""
		}
	}

	if opts.override != "" {
		shouldEmpty := opts.override == emptyOverrideFlag && (len(opts.prefix) > 0 || len(opts.suffix) > 0 || opts.uniqueSuffix)
		if shouldEmpty {
			nameWoutExt = ""
		} else {
			nameWoutExt = opts.override
		}
	} else if replacer != nil {
		nameWoutExt = replacer.regex.ReplaceAllString(nameWoutExt, replacer.replacement)
	}

	var withSeparator = func(a, b string) string {
		if a == "" {
			return b
		}
		if b == "" {
			return a
		}
		return a + opts.separator + b
	}
	if opts.createdDate != "" {
		createdTime := getFileCreatedTime(f)

		if strings.HasPrefix(opts.createdDate, suffixFlag) {
			nameWoutExt = withSeparator(nameWoutExt, createdTime.Format(utils.ConvertDateLayout(opts.createdDate[len(suffixFlag):])))
		} else {
			nameWoutExt = withSeparator(createdTime.Format(utils.ConvertDateLayout(opts.createdDate)), nameWoutExt)
		}
	}

	if opts.detectResolution != "" {
		w, h := getResolution(f, opts.path)
		if w > 0 && h > 0 {
			resolution := fmt.Sprintf("%dx%d", w, h)
			if opts.detectResolution == suffixFlag {
				nameWoutExt = withSeparator(nameWoutExt, resolution)
			} else {
				nameWoutExt = withSeparator(resolution, nameWoutExt)
			}
		}
	}

	if opts.prefix != "" {
		nameWoutExt = withSeparator(opts.prefix, nameWoutExt)
	}

	if opts.suffix != "" {
		nameWoutExt = withSeparator(nameWoutExt, opts.suffix)
	}

	if opts.uniqueSuffix {
		nameWoutExt = withSeparator(nameWoutExt, utils.GenUniqueStr())
	}

	newName = nameWoutExt + ext
	return oldName == newName, newName
}

func getReplaceRegex(replace string) (*replacer, error) {
	if replace == "" {
		return nil, nil
	}

	parts := strings.Split(replace, "=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid replace format: %s, expected: old=new", replace)
	}

	regex, err := regexp.Compile(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to compile replace regex: %s, error: %v", parts[0], err)
	}

	return &replacer{regex: regex, replacement: parts[1]}, nil
}

func processRename(path string, flags *cliFlags, replacer *replacer) map[string]string {
	fmt.Println("Processing...")

	items := getItemInDir(path)
	renamed := make(map[string]string, len(items))
	shouldSyncProcess := flags.detectResolution == ""

	if shouldSyncProcess {
		for _, item := range items {
			if !flags.allowDir && item.IsDir() {
				continue
			}

			ignored, newName := getRenamedName(item, flags, replacer)
			if ignored {
				continue
			}

			// Avoid duplicate file names by adding a unique string
			if _, exists := renamed[newName]; exists {
				ext := filepath.Ext(newName)
				nameWoutExt := strings.TrimSuffix(newName, ext)
				newName = fmt.Sprintf("%s%s%s%s", nameWoutExt, flags.separator, utils.GenUniqueStr(), ext)
			}

			renamed[newName] = item.Name()
		}
		return renamed
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	const maxWorkers = 100
	semaphore := make(chan struct{}, maxWorkers)

	for _, item := range items {
		if !flags.allowDir && item.IsDir() {
			continue
		}

		wg.Add(1)
		go func(item os.DirEntry) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ignored, newName := getRenamedName(item, flags, replacer)
			if ignored {
				return
			}

			// Avoid duplicate file names by adding a unique string
			mu.Lock()
			if _, exists := renamed[newName]; exists {
				ext := filepath.Ext(newName)
				nameWoutExt := strings.TrimSuffix(newName, ext)
				newName = fmt.Sprintf("%s%s%s%s", nameWoutExt, flags.separator, utils.GenUniqueStr(), ext)
			}
			renamed[newName] = item.Name()
			mu.Unlock()
		}(item)
	}
	wg.Wait()

	return renamed
}

func Execute() {
	flags := parseFlags()

	if flags.detectResolution != "" && !isHasFFProbe() {
		fmt.Println("detect resolution flag is set, but ffprobe is not installed. Please install it to use this feature")
		os.Exit(1)
	}

	path := flags.path
	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory", err)
			os.Exit(1)
		}
		path = currentPath
	}

	replacer, err := getReplaceRegex(flags.replace)
	if err != nil {
		fmt.Println("Failed to get replace regex", err)
		os.Exit(1)
	}

	renamed := processRename(path, flags, replacer)

	if len(renamed) == 0 {
		fmt.Println("No files to rename!")
		return
	}

	var displaySummary = func() {
		fmt.Printf("\n--- Summary ---\n")
		fmt.Printf("Path: %s\n", flags.path)
		fmt.Printf("Number of files to rename: %d\n", len(renamed))

	}

	// Run in dry run mode
	if flags.dryRun {
		displaySummary()
		fmt.Println("--- Dry run mode, will not rename the files ---")
		fmt.Println("------------------------------------------------")

		for newName, oldName := range renamed {
			oldPath, newPath := filepath.Join(flags.path, oldName), filepath.Join(flags.path, newName)
			fmt.Printf("%s ➡️  %s\n", oldPath, newPath)
		}

		return
	}

	// Show summary and ask for confirmation
	if !flags.yes {
		displaySummary()
		if !utils.ConfirmAction("Do you want to continue? (Y/n): ", true) {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	successCount := 0
	for newName, oldName := range renamed {
		oldPath, newPath := filepath.Join(flags.path, oldName), filepath.Join(flags.path, newName)
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Failed to rename %s ➡️  %s: %v\n", oldName, newName, err)
			continue
		}
		successCount++
	}

	fmt.Printf("🍀 Successfully renamed %d files\n", successCount)
}
