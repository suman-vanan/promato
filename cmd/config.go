/*
Copyright © 2023 Suman Vanan suman.vanan@live.com
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	neturl "net/url"
)

var url string
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure promato to use a specific Prometheus instance",
	Run: func(cmd *cobra.Command, args []string) {
		handleConfigCmd()
	},
}

func handleConfigCmd() {
	fmt.Println("Validating URL...")
	_, err := neturl.ParseRequestURI(url)
	cobra.CheckErr(err)
	color.Green("URL is valid")
	viper.Set("url", url)
	err = viper.WriteConfig()
	cobra.CheckErr(err)
	color.Green("URL successfully set in config file: %s", viper.ConfigFileUsed())
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&url, "url", "", "Sets the URL in the config file")
	configCmd.MarkFlagRequired("url")
}
