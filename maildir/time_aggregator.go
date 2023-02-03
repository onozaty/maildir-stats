package maildir

import (
	"time"
)

type TimeAggregator struct {
	resultByTime map[string]*AggregateResult
	timeToName   func(time time.Time) string
}

func NewYearAggregator() *TimeAggregator {
	return &TimeAggregator{
		resultByTime: map[string]*AggregateResult{},
		timeToName: func(time time.Time) string {
			if time.Unix() == 0 {
				return ""
			}
			return time.Format("2006")
		},
	}
}

func NewMonthAggregator() *TimeAggregator {
	return &TimeAggregator{
		resultByTime: map[string]*AggregateResult{},
		timeToName: func(time time.Time) string {
			if time.Unix() == 0 {
				return ""
			}
			return time.Format("2006-01")
		},
	}
}

func (a *TimeAggregator) StartUser(userName string) {
	// 何もしない
}

func (a *TimeAggregator) StartMailFolder(mailFolderName string) {
	// 何もしない
}

func (a *TimeAggregator) Aggregate(mail mailInfo) {

	name := a.timeToName(mail.time)

	result, ok := a.resultByTime[name]
	if !ok {
		result = &AggregateResult{
			Name:      name,
			Count:     0,
			TotalSize: 0,
		}
		a.resultByTime[name] = result
	}

	result.Count++
	result.TotalSize += mail.size
}

func (a *TimeAggregator) Results() []*AggregateResult {

	results := []*AggregateResult{}

	for _, result := range a.resultByTime {
		results = append(results, result)
	}

	return results
}
