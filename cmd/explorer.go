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
	Args:  cobra.MaximumNArgs(1),
	Run:   handleExploreCmd(),
}

func handleExploreCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		checkConfig()
		if len(args) == 0 {
			exploreAllMetricSeries()
		}
		exploreSpecificMetricSeries(args[0])
	}
}

func exploreSpecificMetricSeries(seriesName string) {
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
	// fixme: can the hard-coded time range specified below be configurable by the user?
	seriesQueryResponse, warnings, err := v1api.Series(ctx, []string{seriesName}, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	labelValuesByName := map[string]map[string]int{}
	for _, series := range seriesQueryResponse {
		for labelName, labelValue := range series {
			labelName := string(labelName)
			labelValue := string(labelValue)
			if labelName != "__name__" {
				_, ok := labelValuesByName[labelName]
				if !ok {
					labelValuesByName[labelName] = map[string]int{labelValue: 1}
				} else {
					_, ok := labelValuesByName[labelName][labelValue]
					if !ok {
						labelValuesByName[labelName][labelValue] = 1
					} else {
						labelValuesByName[labelName][labelValue]++
					}
				}
			}
		}
	}

	labelCardinalities := map[string]int{}
	for labelName, labelValues := range labelValuesByName {
		labelCardinalities[labelName] = len(labelValues)
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "Label Name"},
			{Align: simpletable.AlignLeft, Text: "Values"},
		},
	}

	for labelName, labelValues := range labelCardinalities {
		row := []*simpletable.Cell{
			{Text: labelName},
			{Text: fmt.Sprintf("%d values", labelValues)},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	fmt.Println(table.String())
}

func exploreAllMetricSeries() {
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

func checkConfig() {
	if !viper.IsSet("url") {
		color.Red("Error: Prometheus API URL is not set. Please use 'config' command to set URL.")
		cobra.CheckErr(errors.New("prometheus url not found in config file"))
	}
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
