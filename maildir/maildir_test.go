package maildir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateMailFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	{
		new := createDir(t, temp, "new")
		createFile(t, new, "1", "1")
		createFile(t, new, "2", "1234567890")
	}
	{
		cur := createDir(t, temp, "cur")
		createFile(t, cur, "3", "123")
		createFile(t, cur, "4", "1234")

		// ディレクトリがあっても無視されること
		curChild := createDir(t, cur, "child")
		createFile(t, curChild, "a", "a")
	}
	{
		// tmpは対象外
		tmp := createDir(t, temp, "tmp")
		createFile(t, tmp, "5", "1")
	}

	createFile(t, temp, "a", "a")
	createFile(t, temp, "xxx", "xxx")

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
		createFile(t, tmp, "5", "1")
	}

	createFile(t, temp, "a", "a")
	createFile(t, temp, "xxx", "xxx")

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

func createFile(t *testing.T, dir string, name string, content string) {

	file, err := os.Create(filepath.Join(dir, name))
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)
}
