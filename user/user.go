package user

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Name    string
	HomeDir string
}

func UsersFromPasswd(passwdPath string) ([]User, error) {

	file, err := os.Open(passwdPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	users := []User{}
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if line == "" || line[0:1] == "#" {
			// 空行とコメントは無視
			continue
		}

		parts := strings.Split(fileScanner.Text(), ":")

		// 0: 名前
		// 5: homeディレクトリ
		if len(parts) < 7 {
			return nil, fmt.Errorf("illegal format of passwd file")
		}

		users = append(users, User{
			Name:    parts[0],
			HomeDir: parts[5],
		})
	}

	return users, nil
}
