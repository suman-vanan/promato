/*
Copyright © 2023 Suman Vanan suman.vanan@live.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sort"
	"time"
)

var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Metrics Explorer",
	Long:  `Metrics Explorer`,
	Run: func(cmd *cobra.Command, args []string) {
		handleExploreCmd()
	},
}

func handleExploreCmd() {
	if !viper.IsSet("url") {
		color.Red("Error: Prometheus API URL is not set. Please use 'config' command to set URL.")
		cobra.CheckErr(errors.New("prometheus url not found in config file"))
	}

	client, err := api.NewClient(api.Config{
		Address: viper.GetString("url"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	metadataResponse, err := v1api.Metadata(ctx, "", "")
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Help"},
			{Align: simpletable.AlignCenter, Text: "Unit"},
		},
	}

	metricNames := make([]string, len(metadataResponse))
	i := 0
	for name := range metadataResponse {
		metricNames[i] = name
		i++
	}
	sort.Strings(metricNames)
	for _, name := range metricNames {
		metadata := metadataResponse[name][0]
		// fixme: some "help" strings are too long to fit neatly into a table row
		if len(metadata.Help) > 80 {
			metadata.Help = "<Help text is too long to fit within table>"
		}
		row := []*simpletable.Cell{
			{Text: name},
			{Text: string(metadata.Type)},
			{Text: metadata.Help},
			{Text: metadata.Unit},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	fmt.Println(table.String())
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
