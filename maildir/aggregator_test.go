package maildir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByName(t *testing.T) {

	// ARRANGE
	results := []*AggregateResult{
		{Name: "b", Count: 1, TotalSize: 10},
		{Name: "a", Count: 20, TotalSize: 30},
		{Name: "aa", Count: 0, TotalSize: 0},
		{Name: "c", Count: 2, TotalSize: 3},
	}

	// ACT
	SortByName(results)

	// ASSERT
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "a", Count: 20, TotalSize: 30},
			{Name: "aa", Count: 0, TotalSize: 0},
			{Name: "b", Count: 1, TotalSize: 10},
			{Name: "c", Count: 2, TotalSize: 3},
		},
		results,
	)
}

func TestSortByCount(t *testing.T) {

	// ARRANGE
	results := []*AggregateResult{
		{Name: "b", Count: 10, TotalSize: 10},
		{Name: "a", Count: 10, TotalSize: 30},
		{Name: "aa", Count: 0, TotalSize: 0},
		{Name: "c", Count: 2, TotalSize: 3},
	}

	// ACT
	SortByCount(results)

	// ASSERT
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "aa", Count: 0, TotalSize: 0},
			{Name: "c", Count: 2, TotalSize: 3},
			{Name: "a", Count: 10, TotalSize: 30},
			{Name: "b", Count: 10, TotalSize: 10},
		},
		results,
	)
}

func TestSortByTotalSize(t *testing.T) {

	// ARRANGE
	results := []*AggregateResult{
		{Name: "b", Count: 1, TotalSize: 10},
		{Name: "a", Count: 2, TotalSize: 2},
		{Name: "aa", Count: 3, TotalSize: 1},
		{Name: "c", Count: 4, TotalSize: 2},
	}

	// ACT
	SortByTotalSize(results)

	// ASSERT
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "aa", Count: 3, TotalSize: 1},
			{Name: "a", Count: 2, TotalSize: 2},
			{Name: "c", Count: 4, TotalSize: 2},
			{Name: "b", Count: 1, TotalSize: 10},
		},
		results,
	)
}
