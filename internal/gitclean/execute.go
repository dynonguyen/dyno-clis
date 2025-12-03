package gitclean

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/dynonguyen/dyno-clis/internal/utils"
)

type cliFlags struct {
	excludes             string // Ex: branch1,branch2
	yes                  bool
	fetchPrune           bool   // Run git fetch --prune --all
	keepNoMergedBranches string // Ex: origin/main,origin/develop
	keepRegex            string // Ex: ^(main|master|production|prod)$
}

const (
	cliName = "clean-git-branch"
)

var defaultFlags = &cliFlags{
	excludes:             "main,master,production,prod",
	yes:                  false,
	fetchPrune:           true,
	keepNoMergedBranches: "origin/master",
	keepRegex:            "",
}

func getLocalBranches() []string {
	local, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		fmt.Println("Failed to get local branches", err)
		return []string{}
	}

	return slices.DeleteFunc(strings.Split(string(local), "\n"), func(branch string) bool {
		return strings.TrimSpace(branch) == ""
	})
}

func getCurrentBranch() string {
	current, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		fmt.Println("Failed to get current branch", err)
		return ""
	}
	return strings.TrimSpace(string(current))
}

func getNoMergedBranches(keepNoMergedBranches string) []string {
	if keepNoMergedBranches == "" {
		return []string{}
	}

	noMergedBranches, err := exec.Command("git", "branch", "--format=%(refname:short)", "--no-merged", keepNoMergedBranches).Output()
	if err != nil {
		fmt.Println("Failed to get no merged branches", err)
		return []string{}
	}

	return slices.DeleteFunc(strings.Split(string(noMergedBranches), "\n"), func(branch string) bool {
		return strings.TrimSpace(branch) == ""
	})
}

func deleteBranches(branches []string) {
	for _, branch := range branches {
		if err := exec.Command("git", "branch", "-D", branch).Run(); err != nil {
			fmt.Println("Failed to delete branch", branch, err)
		}
	}
}

func fetchPrune() {
	fmt.Println("Fetching and pruning remote branches...")
	exec.Command("git", "fetch", "--prune", "--all").Run()
}

func parseFlags() *cliFlags {
	flags := defaultFlags

	flagItems := []utils.FlagItem{
		{
			Name:       "exclude",
			Desc:       "Comma-separated list of branches to exclude",
			Flags:      []string{"e", "exclude"},
			DefaultVal: defaultFlags.excludes,
			StrVal:     &flags.excludes,
		},
		{
			Name:       "yes",
			Desc:       "Skip confirmation",
			Flags:      []string{"y", "yes"},
			DefaultVal: defaultFlags.yes,
			BoolVal:    &flags.yes,
		},
		{
			Name:       "fetch prune",
			Desc:       "Run git fetch --prune --all before deleting branches",
			Flags:      []string{"f", "fetch", "prune"},
			DefaultVal: defaultFlags.fetchPrune,
			BoolVal:    &flags.fetchPrune,
		},
		{
			Name:       "keep no merged branches",
			Desc:       "Keep branches that are not merged into the specified branch, empty to not keep any",
			Flags:      []string{"k", "keep-no-merged"},
			DefaultVal: defaultFlags.keepNoMergedBranches,
			StrVal:     &flags.keepNoMergedBranches,
		},
		{
			Name:       "keep regex",
			Desc:       "Keep branches that match the regex",
			Flags:      []string{"r", "keep-regex"},
			DefaultVal: defaultFlags.keepRegex,
			StrVal:     &flags.keepRegex,
		},
	}

	utils.ParseFlags(flagItems, cliName+" -e branch1,branch2")

	return flags
}

func getKeepRegexPattern(keepRegex string) *regexp.Regexp {
	if keepRegex == "" {
		return nil
	}
	regex, err := regexp.Compile(keepRegex)
	if err != nil {
		fmt.Println("Failed to compile keep regex", err)
		return nil
	}
	return regex
}

func Execute() {
	flags := parseFlags()

	if flags.fetchPrune {
		fetchPrune()
	}

	localBranches := getLocalBranches()
	currentBranch := getCurrentBranch()
	noMergedBranches := getNoMergedBranches(flags.keepNoMergedBranches)

	excludedBranches := strings.Split(flags.excludes, ",")
	excludes := append(append(excludedBranches, noMergedBranches...), currentBranch)

	keepRegexPattern := getKeepRegexPattern(flags.keepRegex)
	deletedBranches := []string{}
	remainingBranches := []string{}

	for _, branch := range localBranches {
		if slices.Contains(excludes, branch) || (keepRegexPattern != nil && keepRegexPattern.MatchString(branch)) {
			remainingBranches = append(remainingBranches, branch)
		} else {
			deletedBranches = append(deletedBranches, branch)
		}
	}

	if len(deletedBranches) == 0 {
		fmt.Println("No branches to delete.")
		return
	}

	// Confirm action
	if !flags.yes {
		fmt.Printf("\n--- Summary ---\n")
		fmt.Printf("- Current branch: %s\n", currentBranch)

		if len(excludedBranches) > 0 {
			fmt.Printf("- Excludes (%d): %s\n", len(excludedBranches), strings.Join(excludedBranches, ", "))
		}

		if len(noMergedBranches) > 0 {
			fmt.Printf("- Keep no merged branches (%d): %s\n", len(noMergedBranches), strings.Join(noMergedBranches, ", "))
		}

		fmt.Printf("- ❌ Deleted branches (%d): %s\n", len(deletedBranches), strings.Join(deletedBranches, ", "))
		fmt.Printf("- ✅ Remaining branches (%d): %s\n", len(remainingBranches), strings.Join(remainingBranches, ", "))

		fmt.Printf("\n⚠️  WARNING: This will delete branches and may cause conflicts\n")

		if !utils.ConfirmAction("Do you want to continue? (Y/n): ", true) {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Delete branches
	deleteBranches(deletedBranches)

	fmt.Println("🍀 Done! 🍀")
}
