package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/onozaty/maildir-stats/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserListCmd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user2:" + filepath.Join(temp, "user2", maildir) + "\n" +
		"user3:" + filepath.Join(temp, "user3", maildir) + "\n" +
		"user4:" + filepath.Join(temp, "user4", maildir) + "\n" +
		"user6:" + filepath.Join(temp, "user6", maildir) + "\n" +
		"user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_SizeLower(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--size-lower", "11",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user3:" + filepath.Join(temp, "user3", maildir) + "\n" +
		"user6:" + filepath.Join(temp, "user6", maildir) + "\n" +
		"user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_SizeUpper(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--size-upper", "11",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user2:" + filepath.Join(temp, "user2", maildir) + "\n" +
		"user4:" + filepath.Join(temp, "user4", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_SizeLowerUpper(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--size-lower", "11",
		"--size-upper", "12",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user3:" + filepath.Join(temp, "user3", maildir) + "\n" +
		"user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_CountLower(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--count-lower", "2",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user3:" + filepath.Join(temp, "user3", maildir) + "\n" +
		"user6:" + filepath.Join(temp, "user6", maildir) + "\n" +
		"user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_CountUpper(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--count-upper", "1",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user2:" + filepath.Join(temp, "user2", maildir) + "\n" +
		"user4:" + filepath.Join(temp, "user4", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_CountLowerUpper(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--count-lower", "2",
		"--count-upper", "3",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user1:" + filepath.Join(temp, "user1", maildir) + "\n" +
		"user3:" + filepath.Join(temp, "user3", maildir) + "\n" +
		"user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_SizeLowerUpper_CountLowerUpper(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

	// テスト用にメソッド差し替え
	loadPasswd = func(passwdPath string) ([]user.User, error) {
		return users, nil
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"user-list",
		"-d", maildir,
		"--size-lower", "12",
		"--size-upper", "12",
		"--count-lower", "2",
		"--count-upper", "2",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	expected := "user7:" + filepath.Join(temp, "user7", maildir) + "\n"
	assert.Equal(t, expected, result)
}

func TestUserListCmd_PasswdFileNotFound(t *testing.T) {

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
		"user-list",
		"-d", maildir,
	})

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + dummyPasswdPath
	assert.Contains(t, err.Error(), expect)
}

func TestUserListCmd_FolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	maildir := "Maildir"

	users := setupTestUserListMaildir(t, temp, maildir)

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
		"user-list",
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

func setupTestUserListMaildir(t *testing.T, temp string, maildir string) []user.User {

	// 名前でソートするので順番を適当に
	users := []user.User{}
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
		// user1
		userName := "user1"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"new/1", 10},
			{"cur/2", 1},
		})
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
			{"cur/1", 10},
			{"cur/2", 1},
			{"cur/3", 1},
		})
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
			{"cur/1", 10},
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
	{
		// user6
		userName := "user6"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"cur/1", 10},
			{"cur/2", 1},
			{"cur/3", 1},
			{"cur/4", 1},
		})
	}
	{
		// user7
		userName := "user7"
		homeDir := createDir(t, temp, userName)
		users = append(users, user.User{
			Name:    userName,
			HomeDir: homeDir,
		})

		mailDir := createDir(t, homeDir, maildir)
		createMailFolder(t, mailDir, []mail{
			{"cur/1", 10},
			{"cur/2", 2},
		})
	}

	return users
}
