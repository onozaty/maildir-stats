package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/onozaty/maildir-stats/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllCmd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_Maildir(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "xxx" // デフォルトとは異なる名前で

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_User(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-u",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[User]
  Name  | Number of mails | Total size(byte)  
--------+-----------------+-------------------
  user1 |               6 |               21  
  user2 |               2 |              300  
  user3 |               3 |            6,000  
  user4 |               0 |                0  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_User_Sort(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-u",
		"--sort-user", "name-desc", // ソートは1パターン試す
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[User]
  Name  | Number of mails | Total size(byte)  
--------+-----------------+-------------------
  user4 |               0 |                0  
  user3 |               3 |            6,000  
  user2 |               2 |              300  
  user1 |               6 |               21  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_Year(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-y",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2021 |               2 |              300  
  2022 |               6 |            3,014  
  2023 |               3 |            3,007  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_Year_Sort(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-y",
		"--sort-year", "size-asc", // ソートは1パターン試す
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2021 |               2 |              300  
  2023 |               3 |            3,007  
  2022 |               6 |            3,014  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_Month(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-m",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2021-12 |               2 |              300  
  2022-11 |               3 |            1,006  
  2022-12 |               3 |            2,008  
  2023-01 |               2 |            3,003  
  2023-02 |               1 |                4  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_Month_Sort(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-m",
		"--sort-month", "count-desc", // ソートは1パターン試す
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2022-12 |               3 |            2,008  
  2022-11 |               3 |            1,006  
  2023-01 |               2 |            3,003  
  2021-12 |               2 |              300  
  2023-02 |               1 |                4  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_User_Year_Month(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-u", "-y", "-m",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := `[Summary]
Number of mails : 11
Total size      : 6,321 byte

[User]
  Name  | Number of mails | Total size(byte)  
--------+-----------------+-------------------
  user1 |               6 |               21  
  user2 |               2 |              300  
  user3 |               3 |            6,000  
  user4 |               0 |                0  

[Year]
  Year | Number of mails | Total size(byte)  
-------+-----------------+-------------------
  2021 |               2 |              300  
  2022 |               6 |            3,014  
  2023 |               3 |            3,007  

[Month]
  Month   | Number of mails | Total size(byte)  
----------+-----------------+-------------------
  2021-12 |               2 |              300  
  2022-11 |               3 |            1,006  
  2022-12 |               3 |            2,008  
  2023-01 |               2 |            3,003  
  2023-02 |               1 |                4  

`
	assert.Equal(t, expected, result)
}

func TestAllCmd_PasswdFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	dummyPasswdPath := filepath.Join(temp, "passwd")

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return user.UsersFromPasswd(dummyPasswdPath)
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + dummyPasswdPath
	assert.Contains(t, err.Error(), expect)
}

func TestAllCmd_FolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// Maildirはあるが、その配下にnew/cur/tmpが無い
	userName := "user9"
	homeDir := createDir(t, temp, userName)
	users = append(users, user.User{
		Name:    userName,
		HomeDir: homeDir,
	})
	mailDir := createDir(t, homeDir, "Maildir")

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + filepath.Join(mailDir, "new")
	assert.Contains(t, err.Error(), expect)
}

func TestAllCmd_InvalidSortUser(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-u",
		"--sort-user", "xxx",
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func TestAllCmd_InvalidSortYear(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-y",
		"--sort-year", "xxx",
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func TestAllCmd_InvalidSortMonth(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestAllMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"all",
		"-d", maildir,
		"-y",
		"--sort-month", "xxx",
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.EqualError(t, err, "invalid sort condition 'xxx'")
}

func setupTestAllMaildir(t *testing.T, temp string, maildir string) []user.User {

	users := []user.User{}
	{
		// user1
		userName := "user1"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"new/1667260800", 1}, // 2022-11-01
			{"cur/1669852800", 2}, // 2022-12-01
		})
		{
			sub := createDir(t, mailDir, ".A")
			createMailFolder(t, sub, []mail{
				{"new/1672531200", 3}, // 2023-01-01
				{"cur/1675209600", 4}, // 2023-02-01
			})
		}
		{
			sub := createDir(t, mailDir, ".B")
			createMailFolder(t, sub, []mail{
				{"new/1669766400", 5}, // 2022-11-30
				{"cur/1672444800", 6}, // 2022-12-31
			})
		}
	}
	{
		// user2
		userName := "user2"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"cur/1640908800", 100}, // 2021-12-31
		})
		{
			sub := createDir(t, mailDir, ".Z")
			createMailFolder(t, sub, []mail{
				{"new/1638316800", 200}, // 2021-12-01
			})
		}
	}
	{
		// user3
		userName := "user3"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"cur/1667260800", 1000}, // 2022-11-01
			{"cur/1669852800", 2000}, // 2022-12-01
			{"cur/1672531200", 3000}, // 2023-01-01
		})
	}
	{
		// user4
		userName := "user4"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			// メール無し
		})
	}
	{
		// user5
		userName := "user5"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		// Maildirなし
	}

	return users
}
