package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCmd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
         |               4 |               10  
  A      |               2 |               30  
  B      |               2 |              300  
  C      |               0 |                0  
  テスト |               2 |            3,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_SortNameAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "name-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
         |               4 |               10  
  A      |               2 |               30  
  B      |               2 |              300  
  C      |               0 |                0  
  テスト |               2 |            3,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_SortNameDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "name-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
  テスト |               2 |            3,000  
  C      |               0 |                0  
  B      |               2 |              300  
  A      |               2 |               30  
         |               4 |               10  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_CountAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "count-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
  C      |               0 |                0  
  A      |               2 |               30  
  B      |               2 |              300  
  テスト |               2 |            3,000  
         |               4 |               10  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_CountDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "count-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
         |               4 |               10  
  テスト |               2 |            3,000  
  B      |               2 |              300  
  A      |               2 |               30  
  C      |               0 |                0  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_SizeAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "size-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
  C      |               0 |                0  
         |               4 |               10  
  A      |               2 |               30  
  B      |               2 |              300  
  テスト |               2 |            3,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_SizeDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "size-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
  テスト |               2 |            3,000  
  B      |               2 |              300  
  A      |               2 |               30  
         |               4 |               10  
  C      |               0 |                0  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2022 |               3 |            3,003  
  2023 |               7 |              337  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_NameAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "name-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2022 |               3 |            3,003  
  2023 |               7 |              337  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_NameDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "name-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2023 |               7 |              337  
  2022 |               3 |            3,003  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_CountAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "count-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2022 |               3 |            3,003  
  2023 |               7 |              337  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_CountDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "count-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2023 |               7 |              337  
  2022 |               3 |            3,003  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_SizeAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "size-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2023 |               7 |              337  
  2022 |               3 |            3,003  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Year_SizeDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "size-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2022 |               3 |            3,003  
  2023 |               7 |              337  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-11 |               1 |            2,000  
  2022-12 |               2 |            1,003  
  2023-01 |               3 |              320  
  2023-02 |               2 |                5  
  2023-03 |               2 |               12  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_NameAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "name-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-11 |               1 |            2,000  
  2022-12 |               2 |            1,003  
  2023-01 |               3 |              320  
  2023-02 |               2 |                5  
  2023-03 |               2 |               12  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_NameDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "name-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2023-03 |               2 |               12  
  2023-02 |               2 |                5  
  2023-01 |               3 |              320  
  2022-12 |               2 |            1,003  
  2022-11 |               1 |            2,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_CountAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "count-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-11 |               1 |            2,000  
  2022-12 |               2 |            1,003  
  2023-02 |               2 |                5  
  2023-03 |               2 |               12  
  2023-01 |               3 |              320  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_CountDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "count-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2023-01 |               3 |              320  
  2023-03 |               2 |               12  
  2023-02 |               2 |                5  
  2022-12 |               2 |            1,003  
  2022-11 |               1 |            2,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_SizeAsc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "size-asc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2023-02 |               2 |                5  
  2023-03 |               2 |               12  
  2023-01 |               3 |              320  
  2022-12 |               2 |            1,003  
  2022-11 |               1 |            2,000  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Month_SizeDesc(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "size-desc",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-11 |               1 |            2,000  
  2022-12 |               2 |            1,003  
  2023-01 |               3 |              320  
  2023-03 |               2 |               12  
  2023-02 |               2 |                5  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_Folder_Year_Month(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f", "-y", "-m",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
         |               4 |               10  
  A      |               2 |               30  
  B      |               2 |              300  
  C      |               0 |                0  
  テスト |               2 |            3,000  

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2022 |               3 |            3,003  
  2023 |               7 |              337  

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-11 |               1 |            2,000  
  2022-12 |               2 |            1,003  
  2023-01 |               3 |              320  
  2023-02 |               2 |                5  
  2023-03 |               2 |               12  

`
	assert.Equal(t, expected, result)
}

func TestUserCmd_MaildirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	maildir := filepath.Join(temp, "xxx") // 存在しないディレクトリ

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", maildir,
	})

	err := rootCmd.Execute()
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + maildir
	assert.Contains(t, err.Error(), expect)
}

func TestUserCmd_InvalidSortFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-f",
		"--sort-folder", "xxx",
	})

	err := rootCmd.Execute()
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func TestUserCmd_InvalidSortYear(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-y",
		"--sort-year", "xxx",
	})

	err := rootCmd.Execute()
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func TestUserCmd_InvalidSortMonth(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	setupTestUserMaildir(t, temp)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user",
		"-d", temp,
		"-m",
		"--sort-month", "xxx",
	})

	err := rootCmd.Execute()
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func setupTestUserMaildir(t *testing.T, rootMailFolderPath string) {

	// INBOX
	createMailFolder(t, rootMailFolderPath, []mail{
		{"new/1675209600", 1}, // 2023-02-01
		{"new/1677628800", 2}, // 2023-03-01
		{"cur/1669852800", 3}, // 2022-12-01
		{"cur/1677542400", 4}, // 2023-02-28
		{"tmp/1677542400", 5}, // tmp配下なので対象外
	})

	// その他フォルダ
	{
		sub := createDir(t, rootMailFolderPath, ".A")
		createMailFolder(t, sub, []mail{
			{"new/1677715200", 10}, // 2023-03-02
			{"cur/1672531200", 20}, // 2023-01-01
		})
	}
	{
		sub := createDir(t, rootMailFolderPath, ".B")
		createMailFolder(t, sub, []mail{
			{"new/1672617600", 100}, // 2023-01-02
			{"new/1672617601", 200}, // 2023-01-02
		})
	}
	{
		sub := createDir(t, rootMailFolderPath, ".C")
		createMailFolder(t, sub, []mail{
			// メール無し
		})
	}
	{
		// マルチバイトが入ったフォルダ名(テスト)
		sub := createDir(t, rootMailFolderPath, ".&MMYwuTDI-")
		createMailFolder(t, sub, []mail{
			{"cur/1672444800", 1000}, // 2022-12-31
			{"cur/1669766400", 2000}, // 2022-11-30
		})
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
