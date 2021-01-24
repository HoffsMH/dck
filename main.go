package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"os/exec"
	"regexp"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "rck",
	Short: "repository check",
	Long: `check all repos recursively from a given directory
					for their up-to-dateness with origin`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := filepath.Abs(args[0])

		caps, err := CapList(dir)
		if err != nil {
			fmt.Println(err)
		}
		for _, cap := range caps {
			o, _ := exec.Command("ls", cap).Output()
			count := len(strings.Split(string(o), "\n"))
			fmt.Println(cap, ": ", count)
		}
	},
}

func CapList(dir string) ([]string, error) {
	var list []string
	permissionDenied := regexp.MustCompile(`permission denied`)

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if !permissionDenied.MatchString(fmt.Sprintf("%+v", err)) {
					return err
				}
				return nil
			}

			if info.IsDir() {
				isCap, _ := IsCapture(path)
				if isCap {
					list = append(list, path)
				}
			}
			return nil
		})
		return list, err
}

func IsCapture(dir string) (bool, error) {
	isCapture := regexp.MustCompile(`capture`)
	return isCapture.MatchString(filepath.Base(dir)), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
