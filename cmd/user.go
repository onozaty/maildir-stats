package cmd

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/onozaty/maildir-stats/maildir"
	"github.com/spf13/cobra"
)

func newUserCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "user",
		Short: "Report user statistics",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirPath, _ := cmd.Flags().GetString("dir")

			reportFolder, _ := cmd.Flags().GetBool("folder")
			reportFolderSortCondition, err := getSortCondition(cmd.Flags(), "sort-folder")
			if err != nil { // 許可されていなパラメータの可能性あり
				return err
			}

			reportYear, _ := cmd.Flags().GetBool("year")
			reportYearSortCondition, err := getSortCondition(cmd.Flags(), "sort-year")
			if err != nil { // 許可されていなパラメータの可能性あり
				return err
			}

			reportMonth, _ := cmd.Flags().GetBool("month")
			reportMonthSortCondition, err := getSortCondition(cmd.Flags(), "sort-month")
			if err != nil { // 許可されていなパラメータの可能性あり
				return err
			}

			inboxFolderName, _ := cmd.Flags().GetString("inbox-name")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runUserReport(
				maildirPath,
				userReportCondition{
					reportFolder:              reportFolder,
					reportFolderSortCondition: reportFolderSortCondition,
					reportYear:                reportYear,
					reportYearSortCondition:   reportYearSortCondition,
					reportMonth:               reportMonth,
					reportMonthSortCondition:  reportMonthSortCondition,
				},
				inboxFolderName,
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("dir", "d", "", "User maildir path.")
	subCmd.MarkFlagRequired("dir")

	subCmd.Flags().BoolP("folder", "f", false, "Report by folder.")
	subCmd.Flags().StringP("sort-folder", "", "name-asc", "Sorting condition for report by folder.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")
	subCmd.Flags().BoolP("year", "y", false, "Report by year.")
	subCmd.Flags().StringP("sort-year", "", "name-asc", "Sorting condition for report by year.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")
	subCmd.Flags().BoolP("month", "m", false, "Report by month.")
	subCmd.Flags().StringP("sort-month", "", "name-asc", "Sorting condition for report by month.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")

	subCmd.Flags().StringP("inbox-name", "", "", "The name of the inbox folder. (default \"\")")

	return subCmd
}

type userReportCondition struct {
	reportFolder              bool
	reportFolderSortCondition SortCondition
	reportYear                bool
	reportYearSortCondition   SortCondition
	reportMonth               bool
	reportMonthSortCondition  SortCondition
}

func runUserReport(maildirPath string, condition userReportCondition, inboxFolderName string, writer io.Writer) error {

	// Summaryを集計するためにもFolderAggregatorはデフォルトで用意する
	folderAggregator := maildir.NewFolderAggregator()
	aggregators := []maildir.Aggregator{folderAggregator}

	var yearAggregator *maildir.TimeAggregator
	var monthAggregator *maildir.TimeAggregator

	if condition.reportYear {
		yearAggregator = maildir.NewYearAggregator()
		aggregators = append(aggregators, yearAggregator)
	}
	if condition.reportMonth {
		monthAggregator = maildir.NewMonthAggregator()
		aggregators = append(aggregators, monthAggregator)
	}

	if err := maildir.AggregateMailFolders(maildirPath, inboxFolderName, maildir.NewMultiAggregator(aggregators)); err != nil {
		return err
	}

	// Summary
	printSummaryReport(writer, folderAggregator)
	fmt.Fprintf(writer, "\n")

	// Folder
	if condition.reportFolder {
		printFolderReport(writer, folderAggregator, condition.reportFolderSortCondition)
		fmt.Fprintf(writer, "\n")
	}

	// Year
	if condition.reportYear {
		printYearReport(writer, yearAggregator, condition.reportYearSortCondition)
		fmt.Fprintf(writer, "\n")
	}

	// Month
	if condition.reportMonth {
		printMonthReport(writer, monthAggregator, condition.reportMonthSortCondition)
		fmt.Fprintf(writer, "\n")
	}

	return nil
}

func printSummaryReport(writer io.Writer, folderAggregator *maildir.FolderAggregator) {

	summaryCount := int32(0)
	summaryTotalSize := int64(0)

	for _, result := range folderAggregator.Results() {
		summaryCount += result.Count
		summaryTotalSize += result.TotalSize
	}

	fmt.Fprintf(writer, "[Summary]\n")
	fmt.Fprintf(writer, "Number of mails : %s\n", humanize.Comma(int64(summaryCount)))
	fmt.Fprintf(writer, "Total size      : %s byte\n", humanize.Comma(int64(summaryTotalSize)))
}

func printFolderReport(writer io.Writer, folderAggregator *maildir.FolderAggregator, sortCondition SortCondition) {

	results := folderAggregator.Results()
	sortResults(results, sortCondition)

	fmt.Fprintf(writer, "[Folder]\n")
	renderTableLayout(writer, results, "Name")
}

func printYearReport(writer io.Writer, yearAggregator *maildir.TimeAggregator, sortCondition SortCondition) {

	results := yearAggregator.Results()
	sortResults(results, sortCondition)

	fmt.Fprintf(writer, "[Year]\n")
	renderTableLayout(writer, results, "Year")
}

func printMonthReport(writer io.Writer, monthAggregator *maildir.TimeAggregator, sortCondition SortCondition) {

	results := monthAggregator.Results()
	sortResults(results, sortCondition)

	fmt.Fprintf(writer, "[Month]\n")
	renderTableLayout(writer, results, "Month")
}

func renderTableLayout(writer io.Writer, results []*maildir.AggregateResult, nameTitle string) {

	table := tablewriter.NewWriter(writer)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	table.SetBorder(false)
	table.SetHeader([]string{nameTitle, "Number of mails", "Total size(byte)"})

	for _, result := range results {
		table.Append(
			[]string{result.Name, humanize.Comma(int64(result.Count)), humanize.Comma(result.TotalSize)})
	}

	table.Render()
}
