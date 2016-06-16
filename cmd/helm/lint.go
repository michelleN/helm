package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/lint"
)

var isRepo bool

var longLintHelp = `
This command takes a path to a chart and runs a series of tests to verify that
the chart is well-formed.

If the linter encounters things that will cause the chart to fail installation,
it will emit [ERROR] messages. If it encounters issues that break with convention
or recommendation, it will emit [WARNING] messages.
`

var lintCommand = &cobra.Command{
	Use:   "lint [flags] PATH",
	Short: "Examines a chart for possible issues",
	Long:  longLintHelp,
	RunE:  lintCmd,
}

func init() {
	lintCommand.Flags().BoolVar(&isRepo, "repo", false, "If set, will lint all charts in a given directory (ignoring index.yaml)")
	RootCommand.AddCommand(lintCommand)
}

var errLintNoChart = errors.New("no chart found for linting (missing Chart.yaml)")

func lintCmd(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	if isRepo {
		if err := lintRepo(path); err != nil {
			return err
		}
	} else {
		// Guard: Error out of this is not a chart.
		if _, err := os.Stat(filepath.Join(path, "Chart.yaml")); err != nil {
			return errLintNoChart
		}

		issues := lint.All(path)
		for _, i := range issues {
			fmt.Printf("%s\n", i)
		}
	}
	return nil
}

func lintRepo(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New(path + " is not a directory")
	}

	chartList := []string{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		name := f.Name()
		if strings.HasSuffix(name, ".tgz") && !f.IsDir() {
			chartList = append(chartList, f.Name())
		}
		return nil
	})

	// chartutil.LoadFile(chartPath)
	//unpack and lint each chart
	return nil
}
