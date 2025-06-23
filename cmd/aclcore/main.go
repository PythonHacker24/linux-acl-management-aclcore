package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-aclcore/config"
	"github.com/PythonHacker24/linux-acl-management-aclcore/internal/manager"
	"github.com/PythonHacker24/linux-acl-management-aclcore/internal/utils"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)
	}
}

func exec() error {

	/* config must load here in exec() if needed to */

	/* setting up cobra for cli interactions */
	var (
		configPath string
		rootCmd    = &cobra.Command{
			Use:   "aclcore <command> <subcommand>",
			Short: "Core Daemon for linux acl management",
			Example: heredoc.Doc(`
				$ aclcore --config /path/to/aclcore.yaml
			`),
			Run: func(cmd *cobra.Command, args []string) {
				if configPath != "" {
					fmt.Printf("Using config file: %s\n\n", configPath)
				} else {
					fmt.Printf("No config file provided.\n\n")
				}
			},
		}
	)

	/* adding --config argument */
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")

	/* Execute the command */
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("arguments error: %s", err.Error())
		os.Exit(1)
	}

	/*
		load config file
		if there is an error in loading the config file, then it will exit with code 1
	*/
	if err := config.LoadConfig(configPath); err != nil {
		fmt.Printf("Configuration Error in %s: %s",
			configPath,
			err.Error(),
		)
		/* since the configuration is invalid, don't proceed */
		os.Exit(1)
	}

	/*
		true for production, false for development mode
		logger is only for core components (after this step)
	*/
	utils.InitLogger(!config.COREDConfig.DConfig.DebugMode)

	/* zap.L() can be used all over the code for global level logging */
	zap.L().Info("Logger Initiated ...")

	/* preparing graceful shutdown */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupt
		cancel()
	}()

	return run(ctx)
}

func run(ctx context.Context) error {

	/* error channel for error propogation */
	errCh := make(chan error, 1)

	/* run the connection pool with manager */
	manager.ConnPool(errCh)

	select {
	case <-ctx.Done():
		zap.L().Info("Shutdown process initiated")
	case err := <-errCh:

		/* context done can be called here (optional for now) */

		zap.L().Error("Fatal Error from core",
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("Shutting down core daemon")

	return nil
}
