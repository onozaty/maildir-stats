package cmd

import (
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"

	"github.com/onozaty/maildir-stats/maildir"
	"github.com/onozaty/maildir-stats/user"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func newUserListCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "user-list",
		Short: "Output user list",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirName, _ := cmd.Flags().GetString("mail-dir")

			sizeLower, _ := cmd.Flags().GetInt64("size-lower")
			sizeUpper, _ := cmd.Flags().GetInt64("size-upper")
			if !cmd.Flags().Changed("size-upper") {
				// 未設定の場合は、int64の最大値いれておく
				sizeUpper = math.MaxInt64
			}

			countLower, _ := cmd.Flags().GetInt64("count-lower")
			countUpper, _ := cmd.Flags().GetInt64("count-upper")
			if !cmd.Flags().Changed("count-upper") {
				// 未設定の場合は、int64の最大値いれておく
				countUpper = math.MaxInt64
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runUserList(
				maildirName,
				userListCondition{
					sizeLower:  sizeLower,
					sizeUpper:  sizeUpper,
					countLower: countLower,
					countUpper: countUpper,
				},
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("mail-dir", "d", "", "User maildir name.")
	subCmd.MarkFlagRequired("mail-dir")

	subCmd.Flags().Int64P("size-lower", "", 0, "Size lower limit.")
	subCmd.Flags().Int64P("size-upper", "", 0, "Size upper limit.")
	subCmd.Flags().Int64P("count-lower", "", 0, "Count lower limit.")
	subCmd.Flags().Int64P("count-upper", "", 0, "Count upper limit.")

	return subCmd
}

type userListCondition struct {
	sizeLower  int64
	sizeUpper  int64
	countLower int64
	countUpper int64
}

func runUserList(maildirName string, condition userListCondition, writer io.Writer) error {

	allUsers, err := loadPasswd(passwdPath)
	if err != nil {
		return err
	}

	// ユーザ毎に集計
	userAggregator := maildir.NewUserAggregator()

	// フォルダ毎での集計はしないので、INBOXは空文字固定で
	if err := maildir.AggregateUsers(allUsers, maildirName, "", userAggregator); err != nil {
		return err
	}

	matchUsers := []user.User{}
	for _, result := range userAggregator.Results() {
		if condition.within(result) {
			index := slices.IndexFunc(allUsers, func(u user.User) bool {
				return u.Name == result.Name
			})
			matchUsers = append(matchUsers, allUsers[index])
		}
	}

	// 一致したユーザは名前でソートして出力
	sort.Slice(matchUsers, func(i, j int) bool {
		return matchUsers[i].Name < matchUsers[j].Name
	})

	for _, user := range matchUsers {
		fmt.Fprintf(writer, "%s:%s\n", user.Name, filepath.Join(user.HomeDir, maildirName))
	}

	return nil
}

func (c *userListCondition) within(result *maildir.AggregateResult) bool {
	return result.Count >= c.countLower && result.Count <= c.countUpper &&
		result.TotalSize >= c.sizeLower && result.TotalSize <= c.sizeUpper
}
