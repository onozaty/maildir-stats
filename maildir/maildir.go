package maildir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/emersion/go-imap/utf7"
)

type AggregateResult struct {
	Name      string
	Count     int32
	TotalSize int64
}

func AggregateMailFolders(rootMailFolderPath string) ([]*AggregateResult, error) {

	aggregateResults := []*AggregateResult{}

	// ルート(INBOX)
	rootAggregateResult, err := aggregateMailFolder(rootMailFolderPath, "")
	if err != nil {
		return nil, err
	}
	aggregateResults = append(aggregateResults, rootAggregateResult)

	// その他メールフォルダ
	entries, err := os.ReadDir(rootMailFolderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// ディレクトリの先頭が"."になっているものがメールフォルダ
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
			mailFolderName, err := decodeFolderName(entry.Name()[1:])
			if err != nil {
				return nil, err
			}

			subAggregateResult, err := aggregateMailFolder(filepath.Join(rootMailFolderPath, entry.Name()), mailFolderName)
			if err != nil {
				return nil, err
			}

			aggregateResults = append(aggregateResults, subAggregateResult)
		}
	}

	sort.Slice(aggregateResults, func(i, j int) bool {
		return aggregateResults[i].Name < aggregateResults[j].Name
	})

	return aggregateResults, nil
}

func aggregateMailFolder(mailFolderPath string, mailFolderName string) (*AggregateResult, error) {

	count := int32(0)
	totalSize := int64(0)

	// tmpにあるのは配送中のものなので対象から除いておく
	for _, subName := range []string{"new", "cur"} {
		subCount, subSize, err := aggregateMails(filepath.Join(mailFolderPath, subName))
		if err != nil {
			return nil, err
		}
		count += subCount
		totalSize += subSize
	}

	return &AggregateResult{
		Name:      mailFolderName,
		Count:     count,
		TotalSize: totalSize,
	}, nil
}

func aggregateMails(dirPath string) (int32, int64, error) {

	count := int32(0)
	totalSize := int64(0)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, 0, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return 0, 0, err
		}

		count++
		totalSize += info.Size()
	}

	return count, totalSize, nil
}

func decodeFolderName(encodedName string) (string, error) {
	decoder := utf7.Encoding.NewDecoder()
	decodedName, err := decoder.String(encodedName)

	if err != nil {
		return "", fmt.Errorf("%s is invalid folder name: %w", encodedName, err)
	}
	return decodedName, nil
}
