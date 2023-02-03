package maildir

type UserAggregator struct {
	results []*AggregateResult
	current *AggregateResult
}

func NewUserAggregator() *UserAggregator {
	return &UserAggregator{
		results: []*AggregateResult{},
	}
}

func (a *UserAggregator) StartUser(userName string) {
	a.current = &AggregateResult{
		Name:      userName,
		Count:     0,
		TotalSize: 0,
	}
	a.results = append(a.results, a.current)
}

func (a *UserAggregator) StartMailFolder(mailFolderName string) {
	// 何もしない
}

func (a *UserAggregator) Aggregate(mail mailInfo) {
	a.current.Count++
	a.current.TotalSize += mail.size
}

func (a *UserAggregator) Results() []*AggregateResult {
	return a.results
}
