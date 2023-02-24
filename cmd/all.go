package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/maildir-stats/maildir"
	"github.com/spf13/cobra"
)

func newAllCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "all",
		Short: "Report all users statistics",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirName, _ := cmd.Flags().GetString("mail-dir")

			reportUser, _ := cmd.Flags().GetBool("user")
			reportUserSortCondition, err := getSortCondition(cmd.Flags(), "sort-user")
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

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runAllReport(
				maildirName,
				allReportCondition{
					reportUser:               reportUser,
					reportUserSortCondition:  reportUserSortCondition,
					reportYear:               reportYear,
					reportYearSortCondition:  reportYearSortCondition,
					reportMonth:              reportMonth,
					reportMonthSortCondition: reportMonthSortCondition,
				},
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("mail-dir", "d", "", "User maildir name.")
	subCmd.MarkFlagRequired("mail-dir")

	subCmd.Flags().BoolP("user", "u", false, "Report by user.")
	subCmd.Flags().StringP("sort-user", "", "name-asc", "Sorting condition for report by user.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")
	subCmd.Flags().BoolP("year", "y", false, "Report by year.")
	subCmd.Flags().StringP("sort-year", "", "name-asc", "Sorting condition for report by year.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")
	subCmd.Flags().BoolP("month", "m", false, "Report by month.")
	subCmd.Flags().StringP("sort-month", "", "name-asc", "Sorting condition for report by month.\ncan be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc")

	return subCmd
}

type allReportCondition struct {
	reportUser               bool
	reportUserSortCondition  SortCondition
	reportYear               bool
	reportYearSortCondition  SortCondition
	reportMonth              bool
	reportMonthSortCondition SortCondition
}

func runAllReport(maildirName string, condition allReportCondition, writer io.Writer) error {

	users, err := loadPasswd(passwdPath)
	if err != nil {
		return err
	}

	// Summaryを集計するためにもUserAggregatorはデフォルトで用意する
	userAggregator := maildir.NewUserAggregator()
	aggregators := []maildir.Aggregator{userAggregator}

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

	// フォルダ毎での集計はしないので、INBOXは空文字固定で
	if err := maildir.AggregateUsers(users, maildirName, "", maildir.NewMultiAggregator(aggregators)); err != nil {
		return err
	}

	// Summary
	printSummaryReport(writer, userAggregator.Results())
	fmt.Fprintf(writer, "\n")

	// User
	if condition.reportUser {
		printUserReport(writer, userAggregator, condition.reportUserSortCondition)
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
