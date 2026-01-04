package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "deck",
	Short: "Deck - Magento 2 Docker Development Environment",
	Long:  `A CLI tool to manage multiple Magento 2 projects using Docker with Traefik reverse proxy.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(binMagentoCmd)
}
