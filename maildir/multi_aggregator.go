package maildir

type MultiAggregator struct {
	aggregators []Aggregator
}

func NewMultiAggregator(aggregators []Aggregator) *MultiAggregator {
	return &MultiAggregator{
		aggregators: aggregators,
	}
}

func (a *MultiAggregator) Start(mailFolderName string) {

	for _, aggregator := range a.aggregators {
		aggregator.Start(mailFolderName)
	}
}

func (a *MultiAggregator) Aggregate(mail mailInfo) error {

	for _, aggregator := range a.aggregators {
		aggregator.Aggregate(mail)
	}
	return nil
}
