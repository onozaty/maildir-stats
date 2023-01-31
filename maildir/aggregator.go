package maildir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-imap/utf7"
)

type AggregateResult struct {
	Name      string
	Count     int32
	TotalSize int64
}

type Aggregator interface {
	Start(mailFolderName string)
	Aggregate(mailFolderName string, fileInfo fs.FileInfo) error
}

func AggregateMailFolders(rootMailFolderPath string, aggregators []Aggregator) error {

	// ルート(INBOX)
	if err := aggregateMailFolder(rootMailFolderPath, "", aggregators); err != nil {
		return err
	}

	// その他メールフォルダ
	entries, err := os.ReadDir(rootMailFolderPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// ディレクトリの先頭が"."になっているものがメールフォルダ
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
			mailFolderName, err := decodeFolderName(entry.Name()[1:])
			if err != nil {
				return err
			}

			if err := aggregateMailFolder(filepath.Join(rootMailFolderPath, entry.Name()), mailFolderName, aggregators); err != nil {
				return err
			}
		}
	}

	return nil
}

func aggregateMailFolder(mailFolderPath string, mailFolderName string, aggregators []Aggregator) error {

	for _, aggregator := range aggregators {
		aggregator.Start(mailFolderName)
	}

	// tmpにあるのは配送中のものなので対象から除いておく
	for _, subName := range []string{"new", "cur"} {
		if err := aggregateMails(filepath.Join(mailFolderPath, subName), mailFolderName, aggregators); err != nil {
			return err
		}
	}

	return nil
}

func aggregateMails(dirPath string, mailFolderName string, aggregators []Aggregator) error {

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		for _, aggregator := range aggregators {
			if err := aggregator.Aggregate(mailFolderName, info); err != nil {
				return err
			}
		}
	}

	return nil
}

func decodeFolderName(encodedName string) (string, error) {
	decoder := utf7.Encoding.NewDecoder()
	decodedName, err := decoder.String(encodedName)

	if err != nil {
		return "", fmt.Errorf("%s is invalid folder name: %w", encodedName, err)
	}
	return decodedName, nil
}
