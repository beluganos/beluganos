// -*- coding: utf-8 -*-

package main

import (
	"os"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

const (
	cfgPathDefault = "/etc/beluganos/fibc.conf"
)

func rootCmd(name string) *cobra.Command {
	var verbose bool
	var showCompletion bool

	rootCmd := &cobra.Command{
		Use:   name,
		Short: "Beluganos control command.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showCompletion {
				cmd.GenBashCompletion(os.Stdout)
			} else {
				cmd.Usage()
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show detail messages.")
	rootCmd.PersistentFlags().BoolVar(&showCompletion, "show-completion", false, "Show bash-comnpletion")

	rootCmd.AddCommand(
		containerCmd(),
		statusCmd(),
		fibcCmd(),
		ovsCmd(),
		playbookCmd(),
	)

	return rootCmd
}

func main() {
	if err := rootCmd("ffctl").Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
