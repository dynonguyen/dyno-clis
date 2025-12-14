![Logo](https://www.dropbox.com/scl/fi/3jf52ewmmzejc6cneum81/go-cli.jpeg?rlkey=54pm45ym5el4wxp228uyt46bn&st=jyxmy80l&raw=1)

**Some clis tools make your life more comfortable 🦖**

# new

`new` - Simplify the process of creating files or directories by combining the `mkdir` and `touch` commands.

🥹 Using `mkdir` & `touch`

```sh
# Netesd folder
mkdir -p folder/sub1/sub2
touch folder/sub1/sub2/file.go

# File
touch file2.go

# directory
mkdir dir
mkdir -p dir/sub1

# Multiple files
mkdir -p folder2
cd folder2
touch file2.go file3.go file4.txt
```

☕ Using `new`

```sh
# Netesd folder
new folder/sub1/sub2/file.go

# File
new file2.go

# directory
new dir/
new dir/sub1/

# Multiple files & directories
new "folder2/[file2.go,file3.go,file4.txt]"
```

### Installation

```sh
go install github.com/dynonguyen/dyno-clis/cmd/new@latest
```

### Usage

```sh
# Create file
new file.go
new file.go file2.go file3.js

# Create directory (end with /)
new dir/
new dir/sub1/sub2/ dir2/
new "dir/[sub1,sub2]/" # Equivalent: new dir/sub1/ dir/sub2/

# Create file in directory
new dir/sub1/file.go
new "dir/sub1/[file.go,file2.go]" # Equivalent: new dir/sub1/file1.go dir/sub2/file2.go

# Space character (surrounded by double quotes)
new "dir/orange cat/cat.go"

# Combination
new "dir/[file.go,file2.go]" dir/sub1/file.go file3.js
```

# envtoggle

`envtoggle` - Interactive environment variable toggle tool with multi-select interface for managing environment variables.

This tool provides an intuitive way to toggle multiple environment variables on/off using an interactive prompt. It supports vim navigation mode and can output results to stdout and/or save to file.

### Installation

```sh
go install github.com/dynonguyen/dyno-clis/cmd/envtoggle@latest
```

### Usage

```sh
# Basic usage - toggle environment variables
envtoggle -k KEY1,KEY2,KEY3

# Output to file
envtoggle -k KEY1,KEY2,KEY3 -f .env

# Disable printing to stdout (only write to file)
envtoggle -k KEY1,KEY2,KEY3 -f .env --print=false

# Use custom on/off values
envtoggle -k DEBUG,VERBOSE,LOGGING --on=true --off=false

# Disable vim mode navigation
envtoggle -k KEY1,KEY2 --vim=false
```

### Options

- `-k, --keys` (required): Comma-separated list of environment variable keys to toggle
- `-f, --file`: Output file path to write the export statements
- `-p, --print`: Print output to stdout (default: true)
- `-v, --vim`: Enable vim mode navigation in the prompt (default: true)
- `--on`: Value for enabled/selected keys (default: "1")
- `--off`: Value for disabled/unselected keys (default: "0")

### Examples

**Basic environment toggle:**

```sh
envtoggle -k NODE_ENV,DEBUG,VERBOSE
```

**Save to .env file:**

```sh
envtoggle -k API_ENABLED,CACHE_ENABLED,LOG_LEVEL -f .env --on=true --off=false
```

**Toggle and source the output:**

```sh
# Generate export statements and source them
envtoggle -k DEBUG,VERBOSE -f /tmp/env.sh && source /tmp/env.sh

# Or pipe directly to shell (be careful with this)
envtoggle -k DEBUG,VERBOSE | source /dev/stdin
```

### Interactive Interface

- Use arrow keys or vim keys (j/k) to navigate
- Space to toggle selection
- Enter to confirm selection
- Shows current state of variables based on environment
- Displays total number of keys being managed

# livephoto

`livephoto` - Move all videos in a directory to a "Live Photos" subdirectory and clean up empty directories.

### Installation

```sh
go install github.com/dynonguyen/dyno-clis/cmd/livephoto@latest
```

### Usage

```sh
livephoto -p /path/to/directory

# OR use current directory
livephoto
```

### Options

- `-p, --path`: Path to the directory to clean up (default: current directory)
- `-k, --keepEmptyDirs`: Keep empty directories after moving files (default: false)
- `-n, --livephotosName`: Name of the live photos directory (default: "Live Photos")

# gitclean

`gitclean` - Clean up git branches that are no longer needed.

### Installation

```sh
go install github.com/dynonguyen/dyno-clis/cmd/gitclean@latest
```

### Usage

```sh
# Clean up git branches with default settings
gitclean

# Exclude branches "branch1" and "branch2" and keep branches that match the regex pattern "^release.*$"
gitclean -e "branch1,branch2" -r "^release.*$"

# Clean up git branches and keep branches that have not been merged
gitclean -k "origin/main"
```

### Options

- `-e, --excludes`: Exclude branches from deletion (default: "main,master,production,prod")
- `-y, --yes`: Automatically answer yes to all prompts (default: false)
- `-f, --fetchPrune`: Run git fetch --prune --all before cleaning (default: true)
- `-k, --keepNoMergedBranches`: Keep branches that have not been merged (default: "origin/master")
- `-r, --keepRegex`: Keep branches that match the regex pattern (default: "")
- `-h`: Show help for the command

# renamer

`renamer` - Rename files in a directory with a given prefix, suffix, separator, created date, resolution, include, exclude, replace, and override.

### Installation

```sh
go install github.com/dynonguyen/dyno-clis/cmd/renamer@latest
```

### Usage

```sh
renamer -p /path/to/directory

# OR use current directory
renamer
```

### Options

- `-p, --path`: Path to the directory to rename files in, empty to use current directory (default: current directory)
- `--prefix`: Prefix to add to the file name
- `--suffix`: Suffix to add to the file name
- `--override`: Override the file name with the given name, empty to keep the original name
- `--separator`: Separator to use between the prefix, suffix and the original file name (default: "\_")
- `--allow-dir`: Allow renaming directories (default: false)
- `--created-date`: Add created date to the file name with the given format (example: YYYY-MM-DD or suffixYYYY-MM-DD)
- `--detect-resolution`: Auto detect resolution and add to the file name, only for photo & video files (example: prefix or suffix)
- `--include`: Only rename files that match the given regex
- `--exclude`: Exclude files that match the given regex
- `--replace`: Replace the given string or regex with the given replacement, format: old=new
- `--dry-run`: Display the files that will be renamed without actually renaming them (default: false)
- `-y, --yes`: Skip confirmation prompt and automatically proceed with renaming (default: false)

### Examples

**Add prefix to all files:**

```sh
renamer --prefix "IMG"
# Renames: photo.jpg → IMG_photo.jpg
```

**Add suffix to all files:**

```sh
renamer --suffix "backup"
# Renames: document.pdf → document_backup.pdf
```

**Use custom separator:**

```sh
renamer --prefix "2024" --separator "-"
# Renames: file.txt → 2024-file.txt
```

**Add created date to file names:**

```sh
renamer --created-date "YYYY-MM-DD"
# Renames: photo.jpg → 2024-01-15_photo.jpg

renamer --created-date "suffixYYYY-MM-DD"
# Renames: photo.jpg → photo_2024-01-15.jpg
```

**Auto detect and add resolution (requires ffprobe):**

```sh
renamer --detect-resolution "prefix"
# Renames: photo.jpg → 1920x1080_16x9_photo.jpg

renamer --detect-resolution "suffix"
# Renames: video.mp4 → video_1920x1080_16x9.mp4
```

**Rename only specific files with regex:**

```sh
renamer --include "\.(jpg|png)$"
# Only renames .jpg and .png files

renamer --exclude "backup"
# Excludes files containing "backup" in the name
```

**Replace text in file names:**

```sh
renamer --replace "old=new"
# Renames: oldfile.txt → newfile.txt

renamer --replace "IMG_(\d+)=Photo-$1"
# Renames: IMG_001.jpg → Photo-001.jpg (using regex)
```

**Override file names:**

```sh
renamer --override "renamed"
# Renames all files to: renamed.jpg, renamed.pdf, etc.
# (duplicate names will have random suffix added)
```

**Dry run to preview changes:**

```sh
renamer --prefix "NEW" --dry-run
# Shows what would be renamed without actually renaming
```

**Skip confirmation prompt:**

```sh
renamer --prefix "IMG" --yes
# Automatically proceeds without asking for confirmation
```

**Rename directories:**

```sh
renamer --prefix "folder" --allow-dir
# Also renames directories, not just files
```

**Combine multiple options:**

```sh
renamer --prefix "2024" --created-date "YYYY-MM-DD" --detect-resolution "suffix" --include "\.(jpg|mp4)$"
# Adds prefix, created date, and resolution to jpg and mp4 files only
# Example: photo.jpg → 2024_2024-01-15_photo_1920x1080_16x9.jpg
```
