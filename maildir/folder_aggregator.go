package maildir

type FolderAggregator struct {
	Results []*AggregateResult
	current *AggregateResult
}

func NewFolderAggregator() *FolderAggregator {
	return &FolderAggregator{
		Results: []*AggregateResult{},
	}
}

func (a *FolderAggregator) Start(mailFolderName string) {

	a.current = &AggregateResult{
		Name:      mailFolderName,
		Count:     0,
		TotalSize: 0,
	}
	a.Results = append(a.Results, a.current)
}

func (a *FolderAggregator) Aggregate(mail mailInfo) error {
	a.current.Count++
	a.current.TotalSize += mail.size

	return nil
}
