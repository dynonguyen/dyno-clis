![Logo](https://www.dropbox.com/scl/fi/3jf52ewmmzejc6cneum81/go-cli.jpeg?rlkey=54pm45ym5el4wxp228uyt46bn&st=jyxmy80l&raw=1)

**Some CLI utilities are written in Go**

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
