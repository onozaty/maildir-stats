package maildir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-imap/utf7"
)

func AggregateMailFolders(rootMailFolderPath string, aggregator Aggregator) error {

	// ルート(INBOX)
	aggregator.Start("")
	if err := aggregateMailFolder(rootMailFolderPath, aggregator); err != nil {
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

			aggregator.Start(mailFolderName)
			if err := aggregateMailFolder(filepath.Join(rootMailFolderPath, entry.Name()), aggregator); err != nil {
				return err
			}
		}
	}

	return nil
}

func aggregateMailFolder(mailFolderPath string, aggregator Aggregator) error {

	// tmpにあるのは配送中のものなので対象から除いておく
	for _, subName := range []string{"new", "cur"} {
		if err := aggregateMails(filepath.Join(mailFolderPath, subName), aggregator); err != nil {
			return err
		}
	}

	return nil
}

func aggregateMails(dirPath string, aggregator Aggregator) error {

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

		if err := aggregator.Aggregate(info); err != nil {
			return err
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
