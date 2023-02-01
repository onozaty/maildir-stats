package maildir

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateMailFolders(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	createMailFolder(t, temp, []mail{
		{"new/1", 1},
		{"new/2", 2},
		{"cur/3", 3},
		{"cur/4", 4},
		{"tmp/5", 5},
		{"tmp/6", 6},
	})

	// その他フォルダ
	{
		sub := createDir(t, temp, ".A")
		createMailFolder(t, sub, []mail{
			{"new/11", 11},
			{"cur/12", 12},
			{"tmp/13", 13},
		})
	}
	{
		sub := createDir(t, temp, ".B")
		createMailFolder(t, sub, []mail{
			{"new/21", 21},
		})
	}
	{
		sub := createDir(t, temp, ".C")
		createMailFolder(t, sub, []mail{
			{"cur/31", 31},
		})
	}
	{
		sub := createDir(t, temp, ".D")
		createMailFolder(t, sub, []mail{
			{"tmp/41", 41},
		})
	}
	{
		// マルチバイトが入ったフォルダ名(テスト)
		sub := createDir(t, temp, ".&MMYwuTDI-")
		createMailFolder(t, sub, []mail{
			{"cur/51", 51},
			{"cur/52", 52},
		})
	}
	{
		// メールフォルダ以外のフォルダ(先頭に"."無し)
		sub := createDir(t, temp, "a")
		createMailFolder(t, sub, []mail{
			{"new/61", 61},
			{"cur/62", 62},
			{"tmp/63", 63},
		})
	}

	// ACT
	aggregator := NewFolderAggregator()
	err := AggregateMailFolders(temp, aggregator)

	// ASSERT
	require.NoError(t, err)

	results := aggregator.Results
	SortByName(results)
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "", Count: 4, TotalSize: 10},
			{Name: "A", Count: 2, TotalSize: 23},
			{Name: "B", Count: 1, TotalSize: 21},
			{Name: "C", Count: 1, TotalSize: 31},
			{Name: "D", Count: 0, TotalSize: 0},
			{Name: "テスト", Count: 2, TotalSize: 103},
		},
		results,
	)
}

func TestAggregateMailFolders_RootFolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	// ACT
	aggregator := NewFolderAggregator()
	err := AggregateMailFolders(rootMailFolderPath, aggregator)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + rootMailFolderPath
	assert.Contains(t, err.Error(), expect)
}

func TestAggregateMailFolders_InvalidFolderName(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	createMailFolder(t, temp, []mail{
		{"new/1", 1},
	})

	// その他フォルダ
	{
		// フォルダ名としておかしなもの(修正UTF-7としてデコードできない者)
		sub := createDir(t, temp, ".&A")
		createMailFolder(t, sub, []mail{
			{"cur/2", 2},
		})
	}

	// ACT
	aggregator := NewFolderAggregator()
	err := AggregateMailFolders(temp, aggregator)

	// ASSERT
	assert.EqualError(t, err, "&A is invalid folder name: utf7: invalid UTF-7")
}

func TestAggregateMailFolders_SubFolderAggregateFailed(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	createMailFolder(t, temp, []mail{
		{"new/1", 1},
	})

	// その他フォルダ
	{
		// Maildirとあるべきフォルダ無し
		createDir(t, temp, ".A")
	}

	// ACT
	aggregator := NewFolderAggregator()
	err := AggregateMailFolders(temp, aggregator)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + filepath.Join(temp, ".A", "new")
	assert.Contains(t, err.Error(), expect)
}

func TestAggregateMailFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	createMailFolder(t, temp, []mail{
		{"new/1", 1},
		{"new/2", 10},
		{"cur/3", 3},
		{"cur/4", 4},
		// サブディレクトリは対象外
		{"cur/child/a", 1},
		// tmp配下は対象外
		{"tmp/x", 2},
		{"tmp/y", 3},
		{"z", 10},
	})

	// ACT
	aggregator := NewFolderAggregator()
	aggregator.Start("INBOX")
	err := aggregateMailFolder(temp, aggregator)

	// ASSERT
	require.NoError(t, err)

	results := aggregator.Results
	SortByName(results)
	assert.Equal(
		t,
		[]*AggregateResult{
			{Name: "INBOX", Count: 4, TotalSize: 18},
		},
		results,
	)
}

func TestAggregateMailFolder_SubFolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	{
		// tmpは対象外
		tmp := createDir(t, temp, "tmp")
		createFile(t, filepath.Join(tmp, "5"), "1")
	}

	createFile(t, filepath.Join(temp, "a"), "a")
	createFile(t, filepath.Join(temp, "xxx"), "xxx")

	// ACT
	aggregator := NewFolderAggregator()
	aggregator.Start("INBOX")
	err := aggregateMailFolder(temp, aggregator)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + filepath.Join(temp, "new")
	assert.Contains(t, err.Error(), expect)
}

func TestDecodeFolderName(t *testing.T) {

	{
		result, err := decodeFolderName("A")
		assert.NoError(t, err)
		assert.Equal(t, "A", result)
	}
	{
		result, err := decodeFolderName("a.b.c")
		assert.NoError(t, err)
		assert.Equal(t, "a.b.c", result)
	}
	{
		result, err := decodeFolderName("INBOX.&MMYwuTDI-.A-&MEI-&-1")
		assert.NoError(t, err)
		assert.Equal(t, "INBOX.テスト.A-あ&1", result)
	}
}

func TestDecodeFolderName_DecodeError(t *testing.T) {

	_, err := decodeFolderName("&AAA")
	assert.EqualError(t, err, "&AAA is invalid folder name: utf7: invalid UTF-7")
}

func TestMailInfoOf(t *testing.T) {

	temp := t.TempDir()

	{
		// ARRANGE
		fileInfo := createFile(t, filepath.Join(temp, "1491941793.10000000.example.com.XXXXX"), "1")

		// ACT
		mail := mailInfoOf(fileInfo)

		// ASSERT
		assert.Equal(t, int64(1), mail.size)
		assert.Equal(t, "2017-04-11T20:16:33Z", mail.time.UTC().Format(time.RFC3339))
	}

	{
		// ARRANGE
		// -> 区切り文字の"."なし
		fileInfo := createFile(t, filepath.Join(temp, "1491941793"), "123")

		// ACT
		mail := mailInfoOf(fileInfo)

		// ASSERT
		assert.Equal(t, int64(3), mail.size)
		assert.Equal(t, "2017-04-11T20:16:33Z", mail.time.UTC().Format(time.RFC3339))
	}

	{
		// ARRANGE
		// -> 数字ではない
		fileInfo := createFile(t, filepath.Join(temp, "xxxxxx.10000000.example.com.aaa"), "")

		// ACT
		mail := mailInfoOf(fileInfo)

		// ASSERT
		assert.Equal(t, int64(0), mail.size)
		assert.Equal(t, "1970-01-01T00:00:00Z", mail.time.UTC().Format(time.RFC3339))
	}
}

func createDir(t *testing.T, parent string, name string) string {

	dir := filepath.Join(parent, name)
	err := os.Mkdir(dir, 0777)
	require.NoError(t, err)

	return dir
}

func createFile(t *testing.T, path string, content string) fs.FileInfo {

	file, err := os.Create(path)
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	info, err := file.Stat()
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	return info
}

func createMailFolder(t *testing.T, baseDir string, mails []mail) {

	createDir(t, baseDir, "tmp")
	createDir(t, baseDir, "new")
	createDir(t, baseDir, "cur")

	for _, mail := range mails {

		mailPath := filepath.Join(baseDir, mail.name)
		err := os.MkdirAll(filepath.Dir(mailPath), 0777)
		require.NoError(t, err)

		// サイズだけあっていればよいので中身は適当に
		createFile(t, filepath.Join(baseDir, mail.name), strings.Repeat("x", mail.size))
	}
}

type mail struct {
	name string
	size int
}
