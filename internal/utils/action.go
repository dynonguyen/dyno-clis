package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmAction(message string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	if message != "" {
		fmt.Print(message)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	lowerResp := strings.ToLower(strings.TrimSpace(response))

	if defaultYes && lowerResp == "" {
		return true
	}

	return lowerResp == "y" || lowerResp == "yes"
}
