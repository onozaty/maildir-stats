package maildir

import (
	"io/fs"
	"sort"
)

type FolderAggregator struct {
	results map[string]*AggregateResult
}

func NewFolderAggregator() *FolderAggregator {
	return &FolderAggregator{
		results: map[string]*AggregateResult{},
	}
}

func (a *FolderAggregator) Start(mailFolderName string) {
	a.results[mailFolderName] = &AggregateResult{
		Name:      mailFolderName,
		Count:     0,
		TotalSize: 0,
	}
}

func (a *FolderAggregator) Aggregate(mailFolderName string, fileInfo fs.FileInfo) error {
	// Startで格納しているので必ず存在する
	result := a.results[mailFolderName]

	result.Count++
	result.TotalSize += fileInfo.Size()

	return nil
}

func (a *FolderAggregator) Results() []*AggregateResult {

	results := []*AggregateResult{}
	for _, result := range a.results {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}
