package maildir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap/utf7"
	"github.com/onozaty/maildir-stats/user"
)

type mailInfo struct {
	size int64
	time time.Time
}

func AggregateUsers(users []user.User, maildirName string, inboxFolderName string, aggregator Aggregator) error {

	for _, user := range users {

		// ユーザのhomeディレクトリにメールディレクトリがあった場合のみ対象に
		userMailFolderPath := filepath.Join(user.HomeDir, maildirName)
		if file, err := os.Stat(userMailFolderPath); err != nil || !file.IsDir() {
			continue
		}

		aggregator.StartUser(user.Name)
		if err := AggregateMailFolders(userMailFolderPath, inboxFolderName, aggregator); err != nil {
			return err
		}
	}
	return nil
}

func AggregateMailFolders(rootMailFolderPath string, inboxFolderName string, aggregator Aggregator) error {

	// ルート(INBOX)
	aggregator.StartMailFolder(inboxFolderName)
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

			aggregator.StartMailFolder(mailFolderName)
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

		aggregator.Aggregate(mailInfoOf(info))
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

func mailInfoOf(fileInfo fs.FileInfo) mailInfo {
	// ファイル名の先頭部分がUnix時間
	// 例: 1674617693.M958571P8888.localhost.localdomain,S=545,W=562:2,S
	//     -> 1674617693 がUnix時間
	unixtimePart := strings.Split(fileInfo.Name(), ".")[0]
	unixtime, err := strconv.ParseInt(unixtimePart, 10, 64)
	if err != nil {
		unixtime = 0
	}

	time := time.Unix(unixtime, 0).UTC() // UTCで扱う

	return mailInfo{
		time: time,
		size: fileInfo.Size(),
	}
}
