package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/ljgago/api-viewer/pkg/log"
	"github.com/ljgago/api-viewer/pkg/server"
)

var (
	version   = "(untracked)"
	commit    = "(untracked)"
	buildDate = "(untracked)"
	system    = "(untracked)"
	arch      = "(untracked)"
)

type Flags struct {
	Port    string
	Proxy   string
	File    string
	Theme   string
	JSON    bool
	Version bool
}

type TemplateData struct {
	Url   string
	Proxy string
}

// rootCmd represents the root
var rootCmd = &cobra.Command{
	Use:   "api-viewer",
	Short: "An OpenAPI Specification (OAS) viewer",
	Long:  `An OpenAPI Specification (OAS) viewer`,
	Run: func(cmd *cobra.Command, args []string) {
		rootRun(cmd, args)
	},
}

func init() {
	rootCmd.Flags().StringP("port", "p", "8090", "server port")
	rootCmd.Flags().StringP("proxy", "x", "8091", "proxy port")
	rootCmd.Flags().StringP("file", "f", "", "spec file")
	rootCmd.Flags().BoolP("json", "j", false, "show log format in json")
	rootCmd.Flags().StringP("theme", "t", "default", "color theme")
	rootCmd.Flags().BoolP("version", "v", false, "version of api-viewer")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseFlags(cmd *cobra.Command) (*Flags, error) {
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return nil, fmt.Errorf("Error getting server port flag -> %w", err)
	}

	proxy, err := cmd.Flags().GetString("proxy")
	if err != nil {
		return nil, fmt.Errorf("Error getting proxy port flag -> %w", err)
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, fmt.Errorf("Error getting file flag -> %w", err)
	}

	theme, err := cmd.Flags().GetString("theme")
	if err != nil {
		return nil, fmt.Errorf("Error getting theme flag -> %w", err)
	}

	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return nil, fmt.Errorf("Error getting log flag -> %w", err)
	}

	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		return nil, fmt.Errorf("Error getting version flag -> %w", err)
	}

	return &Flags{
		Port:    port,
		Proxy:   proxy,
		File:    file,
		Theme:   theme,
		JSON:    json,
		Version: version,
	}, nil
}

func rootRun(cmd *cobra.Command, _ []string) {
	flags, err := parseFlags(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if flags.Version {
		fmt.Printf("api-viewer version %s\nCommit: %s\nBuild Date: %s\nOS/Arch: %s/%s\n", version, commit, buildDate, system, arch)
		return
	}

	if flags.File == "" {
		cmd.Help()
		return
	}

	log.Setup(flags.JSON, "debug")

	var wg sync.WaitGroup

	opts := server.Options{
		Port:  flags.Port,
		Proxy: flags.Proxy,
		File:  flags.File,
		Theme: flags.Theme,
	}

	go server.NewServer(&wg, opts)
	go server.NewProxy(&wg, opts)

	wg.Add(2)
	wg.Wait()
}
