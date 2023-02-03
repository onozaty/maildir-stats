package user

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersFromPasswd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	passwdContents := `root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
nobody:x:99:99:Nobody:/:/sbin/nologin
hanako:x:1001:1001:hanako:/home/hanako:/bin/bash
taro:x:1002:1002:taro:/home/taro:/bin/bash
apache:x:48:48:Apache:/usr/share/httpd:/sbin/nologin`

	passwdPath := filepath.Join(temp, "passwd")
	createFile(t, passwdPath, passwdContents)

	// ACT
	users, err := UsersFromPasswd(passwdPath)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(
		t,
		[]User{
			{Name: "root", HomeDir: "/root"},
			{Name: "bin", HomeDir: "/bin"},
			{Name: "nobody", HomeDir: "/"},
			{Name: "hanako", HomeDir: "/home/hanako"},
			{Name: "taro", HomeDir: "/home/taro"},
			{Name: "apache", HomeDir: "/usr/share/httpd"},
		},
		users)
}

func TestUsersFromPasswd_Empty_Comment(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	passwdContents := `root:x:0:0:root:/root:/bin/bash

hanako:x:1001:1001:hanako:/home/hanako:/bin/bash
#taro:x:1002:1002:taro:/home/taro:/bin/bash
apache:x:48:48:Apache:/usr/share/httpd:/sbin/nologin`

	passwdPath := filepath.Join(temp, "passwd")
	createFile(t, passwdPath, passwdContents)

	// ACT
	users, err := UsersFromPasswd(passwdPath)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(
		t,
		[]User{
			{Name: "root", HomeDir: "/root"},
			{Name: "hanako", HomeDir: "/home/hanako"},
			{Name: "apache", HomeDir: "/usr/share/httpd"},
		},
		users)
}

func TestUsersFromPasswd_PasswdFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	passwdPath := filepath.Join(temp, "passwd") // 作成なし

	// ACT
	_, err := UsersFromPasswd(passwdPath)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + passwdPath
	assert.Contains(t, err.Error(), expect)
}

func TestUsersFromPasswd_IllegalFormat(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	passwdContents := "root:x:0:0:root:/root" // おかしなフォーマット
	passwdPath := filepath.Join(temp, "passwd")
	createFile(t, passwdPath, passwdContents)

	// ACT
	_, err := UsersFromPasswd(passwdPath)

	// ASSERT
	assert.EqualError(t, err, "illegal format of passwd file")
}

func createFile(t *testing.T, path string, content string) {

	file, err := os.Create(path)
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)
}
