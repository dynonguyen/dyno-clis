package envtoggle

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dynonguyen/dyno-clis/internal/utils"
)

const (
	cliName = "envtoggle"
)

type cliFlags struct {
	keys              string
	filePath          string
	isPrint           bool
	vimMode           bool
	onValue, offValue string
}

func runCommand(flags *cliFlags) error {
	keys := slices.DeleteFunc(strings.Split(flags.keys, ","), func(s string) bool {
		return strings.TrimSpace(s) == ""
	})

	if len(keys) == 0 {
		return errors.New("please provide --keys")
	}

	// Prompt multi-select
	selected := []string{}
	prompt := &survey.MultiSelect{
		PageSize: 10,
		Message:  fmt.Sprintf("%s (%d keys)", cliName, len(keys)),
		Options:  keys,
		Default: func() []string {
			var d []string
			for _, k := range keys {
				if os.Getenv(k) == flags.onValue {
					d = append(d, k)
				}
			}
			return d
		}(),
		VimMode: flags.vimMode,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return err
	}

	selectedSet := map[string]bool{}
	for _, s := range selected {
		selectedSet[s] = true
	}

	// Print exports
	exportedEnvs := ""
	for _, k := range keys {
		val := flags.offValue
		if selectedSet[k] {
			val = flags.onValue
		}
		exportedEnvs += fmt.Sprintf("export %s=%s\n", k, val)
	}

	if flags.isPrint {
		fmt.Println(exportedEnvs)
	}

	if flags.filePath != "" {
		file, err := os.Create(flags.filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = file.WriteString(exportedEnvs); err != nil {
			return err
		}
	}

	return nil
}

func parseFlags() *cliFlags {
	flags := &cliFlags{
		isPrint: true,
		vimMode: true,
	}

	flagItems := []utils.FlagItem{
		{
			Name:     "keys",
			Flags:    []string{"k", "keys"},
			Desc:     "List of env keys to toggle",
			StrVal:   &flags.keys,
			Required: true,
		},
		{
			Name:   "file",
			Desc:   "Output file to write the exports to",
			Flags:  []string{"f", "file"},
			StrVal: &flags.filePath,
		},
		{
			Name:       "print",
			Desc:       "Print the output to stdout",
			Flags:      []string{"p", "print"},
			BoolVal:    &flags.isPrint,
			DefaultVal: true,
		},
		{
			Name:       "vimMode",
			Desc:       "Enable vim mode in the prompt",
			Flags:      []string{"v", "vim"},
			BoolVal:    &flags.vimMode,
			DefaultVal: true,
		},
		{
			Name:       "onValue",
			Desc:       "Value to set for enabled keys",
			Flags:      []string{"on"},
			StrVal:     &flags.onValue,
			DefaultVal: "1",
		},
		{
			Name:       "offValue",
			Desc:       "Value to set for disabled keys",
			Flags:      []string{"off"},
			StrVal:     &flags.offValue,
			DefaultVal: "0",
		},
	}

	utils.ParseFlags(flagItems, cliName+" -k KEY1,KEY2,KEY3")

	return flags
}

func Execute() {
	flags := parseFlags()
	if err := runCommand(flags); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
