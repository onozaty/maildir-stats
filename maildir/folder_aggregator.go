package maildir

type FolderAggregator struct {
	results []*AggregateResult
	current *AggregateResult
}

func NewFolderAggregator() *FolderAggregator {
	return &FolderAggregator{
		results: []*AggregateResult{},
	}
}

func (a *FolderAggregator) StartMailFolder(mailFolderName string) {

	a.current = &AggregateResult{
		Name:      mailFolderName,
		Count:     0,
		TotalSize: 0,
	}
	a.results = append(a.results, a.current)
}

func (a *FolderAggregator) Aggregate(mail mailInfo) {
	a.current.Count++
	a.current.TotalSize += mail.size
}

func (a *FolderAggregator) Results() []*AggregateResult {
	return a.results
}
