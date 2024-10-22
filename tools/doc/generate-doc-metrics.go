package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"

	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

func generateDocMetrics() {
	tmplFuncs := template.FuncMap{
		"tableSettings": func() string {
			sSlice := [][]string{}

			sSlice = append(sSlice, []string{fmt.Sprintf("--%s", models.MetricsFlagName), "false", "Enable metrics collection"})
			sSlice = append(sSlice, []string{fmt.Sprintf("--%s", models.MetricsPortFlagName), models.MetricsDefaultAddr, "Port to expose metrics on"})
			sSlice = append(sSlice, []string{fmt.Sprintf("--%s", models.MetricsPathFlagName), models.MetricsDefaultPath, "Path to expose metrics on"})

			prettyPrintedTable, err := markdown.NewTableFormatterBuilder().
				WithPrettyPrint().
				Build("Flag", "Default", "Description").
				Format(sSlice)
			if err != nil {
				panic(err)
			}

			return prettyPrintedTable
		},
		"tableMetrics": func() string {
			metrics.InitAll()

			mMap := map[string]string{}

			for _, mm := range metrics.Metrics {
				for name, m := range mm {
					mMap[name] = m.(metrics.Metric).GetHelp()
				}
			}

			// Extract keys from map
			keys := make([]string, 0, len(mMap))
			for k := range mMap {
				keys = append(keys, k)
			}

			// sort the map
			sort.Strings(keys)

			mSlice := [][]string{}
			for _, k := range keys {
				mSlice = append(mSlice, []string{k, mMap[k]})
			}

			prettyPrintedTable, err := markdown.NewTableFormatterBuilder().
				WithPrettyPrint().
				Build("Metrics", "Description").
				Format(mSlice)
			if err != nil {
				panic(err)
			}

			return prettyPrintedTable
		},
	}

	// os read file
	file, err := os.ReadFile("docs/advanced/metrics.md.tmpl")
	if err != nil {
		log.Default().Printf("Failed to open file: %v", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("metrics").Funcs(tmplFuncs).Parse(string(file)))

	// write template to file
	f, err := os.Create("docs/advanced/metrics.md")
	defer f.Close()
	if err != nil {
		log.Default().Printf("Failed to create file: %v", err)
		f.Close()
		os.Exit(1) //nolint: gocritic
	}
	if err := tmpl.Execute(f, nil); err != nil {
		log.Default().Printf("Failed to execute template: %v", err)
		os.Exit(1)
	}

	// os.Exit(0)
}
