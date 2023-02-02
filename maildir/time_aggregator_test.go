package maildir

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateMailFolders_MonthAggregator(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	createMailFolder(t, temp, []mail{
		{"new/1675209600", 1},      // 2023-02-01
		{"new/1677628800.xxxx", 2}, // 2023-03-01
		{"cur/1669852800.1.2", 3},  // 2022-12-01
		{"cur/1677542400", 4},      // 2023-02-28
		{"cur/xxxxxxxxxx", 5},      // 日付無し
	})

	// その他フォルダ
	{
		sub := createDir(t, temp, ".A")
		createMailFolder(t, sub, []mail{
			{"new/1677715200", 11}, // 2023-03-02
			{"cur/1672531200", 12}, // 2023-01-01
		})
	}
	{
		sub := createDir(t, temp, ".B")
		createMailFolder(t, sub, []mail{
			{"new/1672617600", 21}, // 2023-01-02
			{"new/1672617601", 22}, // 2023-01-02
		})
	}
	{
		sub := createDir(t, temp, ".C")
		createMailFolder(t, sub, []mail{
			// メール無し
		})
	}

	// ACT
	aggregator := NewMonthAggregator()
	err := AggregateMailFolders(temp, aggregator)

	// ASSERT
	require.NoError(t, err)

	results := aggregator.Results()
	SortByName(results)
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "", Count: 1, TotalSize: 5},
			{Name: "2022-12", Count: 1, TotalSize: 3},
			{Name: "2023-01", Count: 3, TotalSize: 55},
			{Name: "2023-02", Count: 2, TotalSize: 5},
			{Name: "2023-03", Count: 2, TotalSize: 13},
		},
		results,
	)
}

func TestAggregateMailFolders_YearAggregator(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	createMailFolder(t, temp, []mail{
		{"new/1675209600", 1},      // 2023-02-01
		{"new/1677628800.xxxx", 2}, // 2023-03-01
		{"cur/1669852800.1.2", 3},  // 2022-12-01
		{"cur/1609459200", 4},      // 2021-01-01
		{"cur/xxxxxxxxxx", 5},      // 日付無し
	})

	// その他フォルダ
	{
		sub := createDir(t, temp, ".A")
		createMailFolder(t, sub, []mail{
			{"new/1640908800", 11}, // 2021-12-31
			{"cur/1672531200", 12}, // 2023-01-01
		})
	}
	{
		sub := createDir(t, temp, ".B")
		createMailFolder(t, sub, []mail{
			{"new/1672617600", 21}, // 2023-01-02
		})
	}
	{
		sub := createDir(t, temp, ".C")
		createMailFolder(t, sub, []mail{
			// メール無し
		})
	}

	// ACT
	aggregator := NewYearAggregator()
	err := AggregateMailFolders(temp, aggregator)

	// ASSERT
	require.NoError(t, err)

	results := aggregator.Results()
	SortByName(results)
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "", Count: 1, TotalSize: 5},
			{Name: "2021", Count: 2, TotalSize: 15},
			{Name: "2022", Count: 1, TotalSize: 3},
			{Name: "2023", Count: 4, TotalSize: 36},
		},
		results,
	)
}
