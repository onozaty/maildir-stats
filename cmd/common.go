package cmd

import (
	"fmt"

	"github.com/onozaty/maildir-stats/maildir"
	"github.com/spf13/pflag"
)

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
	case "", "name-asc":
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
