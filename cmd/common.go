package cmd

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/onozaty/maildir-stats/maildir"
	"github.com/onozaty/maildir-stats/user"
	"github.com/spf13/pflag"
)

const passwdPath = "/etc/passwd"

type SortCondition int

const (
	NameAsc SortCondition = iota
	NameDesc
	CountAsc
	CountDesc
	SizeAsc
	SizeDesc
)

func sortResults(results []*maildir.AggregateResult, sortCondition SortCondition) {

	switch sortCondition {
	case NameAsc:
		maildir.SortByName(results)
	case NameDesc:
		maildir.SortByName(results)
		reverse(results)
	case CountAsc:
		maildir.SortByCount(results)
	case CountDesc:
		maildir.SortByCount(results)
		reverse(results)
	case SizeAsc:
		maildir.SortByTotalSize(results)
	case SizeDesc:
		maildir.SortByTotalSize(results)
		reverse(results)
	}
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func getSortCondition(f *pflag.FlagSet, name string) (SortCondition, error) {

	str, _ := f.GetString(name)

	switch str {
	case "name-asc":
		return NameAsc, nil
	case "name-desc":
		return NameDesc, nil
	case "count-asc":
		return CountAsc, nil
	case "count-desc":
		return CountDesc, nil
	case "size-asc":
		return SizeAsc, nil
	case "size-desc":
		return SizeDesc, nil
	default:
		return -1, fmt.Errorf("invalid sort condition '%s'", str)
	}
}

func printSummaryReport(writer io.Writer, results []*maildir.AggregateResult) {

	summaryCount := int64(0)
	summaryTotalSize := int64(0)

	for _, result := range results {
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

func printUserReport(writer io.Writer, userAggregator *maildir.UserAggregator, sortCondition SortCondition) {

	results := userAggregator.Results()
	sortResults(results, sortCondition)

	fmt.Fprintf(writer, "[User]\n")
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

// テスト用に差し替え可能にしておく
var loadPasswd = loadPasswdReal

func loadPasswdReal(passwdPath string) ([]user.User, error) {
	return user.UsersFromPasswd(passwdPath)
}
