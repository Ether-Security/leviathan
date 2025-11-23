package cmd

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/Ether-Security/leviathan/core"
	"github.com/Ether-Security/leviathan/utils"
	"github.com/panjf2000/ants"
	"github.com/spf13/cobra"
)

func init() {
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Launch a workflow",
		Long:  "Launch a workflow",
		RunE:  runScan,
	}

	scanCmd.Flags().StringVarP(&options.Scan.Flow, "flow", "f", "sample", "Flow name for running")
	scanCmd.Flags().StringSliceVarP(&options.Scan.Targets, "targets", "t", []string{}, "Targets to use as input for workflow")
	scanCmd.Flags().StringVarP(&options.Scan.Output, "workspace", "w", "", "Force the workspace directory")
	scanCmd.Flags().StringToStringVarP(&options.Scan.Params, "params", "p", nil, "Custom params -p='foo=bar' (Multiple -p flags are accepted)")
	scanCmd.Flags().IntVarP(&options.Scan.Threads, "threads", "T", 1, "Define the number of concurrent jobs")
	scanCmd.Flags().BoolVar(&options.Scan.NoClean, "no-clean", false, "Disallow modules to execute cleaning scripts")
	scanCmd.Flags().BoolVar(&options.Scan.Resume, "resume", false, "Resume a previous scan")
	rootCmd.AddCommand(scanCmd)

	cobra.OnInitialize(initInput)
}

func initInput() {
	// Detect if an input come from STDIN
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			target := strings.TrimSpace(sc.Text())
			if err := sc.Err(); err == nil && target != "" {
				options.Scan.Targets = append(options.Scan.Targets, target)
			}
		}
	}
}

func runScan(_ *cobra.Command, _ []string) error {
	// Init logger
	utils.InitLog(&options)

	// Check if specific output is defined
	if options.Scan.Output != "" {
		options.Environment.Workspaces = options.Scan.Output
	}

	// Check Workflow
	if !utils.IsYamlValid(options.Scan.Flow, options.Environment.Workflows) {
		utils.Logger.Fatal().Str("workflow", options.Scan.Flow).Msg("Workflow invalid")
		return nil
	}

	// Foreach target launch a runner
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(options.Scan.Threads, func(i interface{}) {
		CreateRunner(i)
		wg.Done()
	})
	defer p.Release()

	for _, target := range options.Scan.Targets {
		wg.Add(1)
		_ = p.Invoke(strings.TrimSpace(target))
	}
	wg.Wait()
	return nil
}

func CreateRunner(j interface{}) {
	target := j.(string)

	runner, err := core.InitRunner(target, &options)
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
		utils.Logger.Error().Msgf("Unable to start the runner for target : %s", target)
		return
	}

	runner.Start()
}
