package maildir

type MultiAggregator struct {
	aggregators []Aggregator
}

func NewMultiAggregator(aggregators []Aggregator) *MultiAggregator {
	return &MultiAggregator{
		aggregators: aggregators,
	}
}

func (a *MultiAggregator) StartUser(userName string) {
	for _, aggregator := range a.aggregators {
		aggregator.StartUser(userName)
	}
}

func (a *MultiAggregator) StartMailFolder(mailFolderName string) {

	for _, aggregator := range a.aggregators {
		aggregator.StartMailFolder(mailFolderName)
	}
}

func (a *MultiAggregator) Aggregate(mail mailInfo) {

	for _, aggregator := range a.aggregators {
		aggregator.Aggregate(mail)
	}
}
