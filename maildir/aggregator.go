package maildir

import (
	"io/fs"
	"sort"
)

type AggregateResult struct {
	Name      string
	Count     int32
	TotalSize int64
}

func SortByName(results []*AggregateResult) {

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
}

func SortByCount(results []*AggregateResult) {

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count == results[j].Count {
			return results[i].Name < results[j].Name
		}
		return results[i].Count < results[j].Count
	})
}

func SortByTotalSize(results []*AggregateResult) {

	sort.Slice(results, func(i, j int) bool {
		if results[i].TotalSize == results[j].TotalSize {
			return results[i].Name < results[j].Name
		}
		return results[i].TotalSize < results[j].TotalSize
	})
}

type Aggregator interface {
	Start(mailFolderName string)
	Aggregate(fileInfo fs.FileInfo) error
}
