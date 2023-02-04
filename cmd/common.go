package cmd

import (
	"fmt"

	"github.com/onozaty/maildir-stats/maildir"
	"github.com/spf13/pflag"
)

type SortCondition int

const (
	Name SortCondition = iota
	Count
	Size
)

func sortResults(results []*maildir.AggregateResult, sortCondition SortCondition) {

	switch sortCondition {
	case Name:
		maildir.SortByName(results)
	case Count:
		maildir.SortByCount(results)
	case Size:
		maildir.SortByTotalSize(results)
	}
}

func getSortCondition(f *pflag.FlagSet, name string) (SortCondition, error) {

	str, _ := f.GetString(name)

	switch str {
	case "", "name":
		return Name, nil
	case "count":
		return Count, nil
	case "size":
		return Size, nil
	default:
		return -1, fmt.Errorf("invalid sort condition '%s'", str)
	}
}
