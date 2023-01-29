package maildir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emersion/go-imap/utf7"
)

type AggregateResult struct {
	Name      string
	Count     int32
	TotalSize int64
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
