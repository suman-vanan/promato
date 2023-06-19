/*
Copyright © 2023 Suman Vanan suman.vanan@live.com
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"time"
)

// explorerCmd represents the explorer command
var explorerCmd = &cobra.Command{
	Use:   "explorer",
	Short: "Metrics Explorer",
	Long:  `Metrics Explorer`,
	Run: func(cmd *cobra.Command, args []string) {
		getMetadata()
	},
}

func getMetadata() {
	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
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
	rootCmd.AddCommand(explorerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explorerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explorerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
