package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gokrazy/internal/config"
	"github.com/gokrazy/internal/instanceflag"
	"github.com/gokrazy/tools/internal/packer"
	"github.com/spf13/cobra"
)

// updateCmd is gok update.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Build a gokrazy instance and update over the network",
	Long: `Build a gokrazy instance and update over the network.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NArg() > 0 {
			fmt.Fprint(os.Stderr, `positional arguments are not supported

`)
			return cmd.Usage()
		}

		return updateImpl.run(cmd.Context(), args, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

type updateImplConfig struct{}

var updateImpl updateImplConfig

func init() {
	instanceflag.RegisterPflags(updateCmd.Flags())
}

func (r *updateImplConfig) run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	cfg, err := config.ReadFromFile()
	if err != nil {
		return err
	}
	log.Printf("cfg: %+v", cfg)

	// gok update is mutually exclusive with gok overwrite
	cfg.InternalCompatibilityFlags.Overwrite = ""
	cfg.InternalCompatibilityFlags.OverwriteBoot = ""
	cfg.InternalCompatibilityFlags.OverwriteRoot = ""
	cfg.InternalCompatibilityFlags.OverwriteMBR = ""

	if cfg.InternalCompatibilityFlags.Update == "" {
		cfg.InternalCompatibilityFlags.Update = "yes"
	}

	if err := os.Chdir(config.InstancePath()); err != nil {
		return err
	}

	packer.Main(cfg)

	return nil
}
