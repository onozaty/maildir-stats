package maildir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	results, err := AggregateMailFolders(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(
		t,
		[]*AggregateResult{
			{
				Name:      "",
				Count:     4,
				TotalSize: 10,
			},
			{
				Name:      "A",
				Count:     2,
				TotalSize: 23,
			},
			{
				Name:      "B",
				Count:     1,
				TotalSize: 21,
			},
			{
				Name:      "C",
				Count:     1,
				TotalSize: 31,
			},
			{
				Name:      "D",
				Count:     0,
				TotalSize: 0,
			},
			{
				Name:      "テスト",
				Count:     2,
				TotalSize: 103,
			},
		},
		results,
	)
}

func TestAggregateMailFolders_RootFolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	// ACT
	_, err := AggregateMailFolders(rootMailFolderPath)

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
	_, err := AggregateMailFolders(temp)

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
	_, err := AggregateMailFolders(temp)

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
	result, err := aggregateMailFolder(temp, "INBOX")

	// ASSERT
	require.NoError(t, err)
	assert.Equal(
		t,
		&AggregateResult{
			Name:      "INBOX",
			Count:     4,
			TotalSize: 18,
		},
		result,
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
	_, err := aggregateMailFolder(temp, "INBOX")

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

func createDir(t *testing.T, parent string, name string) string {

	dir := filepath.Join(parent, name)
	err := os.Mkdir(dir, 0777)
	require.NoError(t, err)

	return dir
}

func createFile(t *testing.T, path string, content string) {

	file, err := os.Create(path)
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)
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
