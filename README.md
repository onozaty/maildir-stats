# maildir-stats

[![GitHub license](https://img.shields.io/github/license/onozaty/maildir-stats)](https://github.com/onozaty/maildir-stats/blob/main/LICENSE)
[![Test](https://github.com/onozaty/maildir-stats/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/maildir-stats/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/maildir-stats/branch/main/graph/badge.svg?token=J7YRQFI233)](https://codecov.io/gh/onozaty/maildir-stats)

`maildir-stats` is a tool that reports maildir statistics.

`maildir-stats` has the following subcommands

* [user](#user) Report user statistics.
* [all](#all) Report all users statistics.
* [user-list](#user-list) Output user list.

## user

Report user statistics.  

### Usage

```
maildir-stats user -d MAIL_DIR_PATH [-f] [--sort-folder SORT_COND] [-y] [--sort-year SORT_COND] [-m] [--sort-month SORT_COND] [--inbox-name INBOX_NAME]
```

```
Usage:
  maildir-stats user [flags]

Flags:
  -d, --dir string           User maildir path.
  -f, --folder               Report by folder.
      --sort-folder string   Sorting condition for report by folder.
                             can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
  -y, --year                 Report by year.
      --sort-year string     Sorting condition for report by year.
                             can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
  -m, --month                Report by month.
      --sort-month string    Sorting condition for report by month.
                             can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
      --inbox-name string    The name of the inbox folder. (default "")
  -h, --help                 help for user
```

### Example

This is the case when all statistics are reported by specifying the maildir of `user1`.

```
$ maildir-stats user -d /home/user1/Maildir -f -y -m

[Summary]
Number of mails : 10
Total size      : 3,340 byte

[Folder]
  Name   | Number of mails | Total size(byte)  
---------+-----------------+-------------------
         |               4 |               10  
  A      |               2 |               30  
  B      |               2 |              300  
  C      |               0 |                0  
  XXXXXX |               2 |            3,000  

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
```

## all

Report all users statistics.  
Target user information is obtained from `/etc/passwd`.

### Usage

```
maildir-stats all -d MAIL_DIR_NAME [-u] [--sort-user SORT_COND] [-y] [--sort-year SORT_COND] [-m] [--sort-month SORT_COND]
```

```
Usage:
  maildir-stats users [flags]

Flags:
  -d, --mail-dir string     User maildir name.
  -u, --user                Report by user.
      --sort-user string    Sorting condition for report by user.
                            can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
  -y, --year                Report by year.
      --sort-year string    Sorting condition for report by year.
                            can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
  -m, --month               Report by month.
      --sort-month string   Sorting condition for report by month.
                            can be specified: name-asc, name-desc, count-asc, count-desc, size-asc, size-desc (default "name-asc")
  -h, --help                help for users
```

### Example

This is the case when all user statistics are reported.  

```
$ maildir-stats users -d Maildir -u -y -m

[Summary]
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

```

## user-list

Output user list.  
Target user information is obtained from `/etc/passwd`.

### Usage

```
maildir-stats user-list -d MAIL_DIR_NAME [--size-lower SIZE] [--size-upper SIZE] [--count-lower COUNT] [--count-upper COUNT]
```

```
Usage:
  maildir-stats user-list [flags]

Flags:
  -d, --mail-dir string   User maildir name.
      --size-lower int    Size lower limit.
      --size-upper int    Size upper limit.
      --count-lower int   Count lower limit.
      --count-upper int   Count upper limit.
  -h, --help              help for user-list
```

### Example

Outputs a list of users whose size is 1000 bytes or larger.

```
$ maildir-stats user-list -d Maildir --size-lower 1000

user1:/home/user1/Maildir
user2:/home/user2/Maildir
user4:/home/user4/Maildir
```

## Install

`maildir-stats` is implemented in golang and runs on all major platforms such as Windows, Mac OS, and Linux.  
You can download the binaries for each OS from the links below.

You can download the binary from the following.

* https://github.com/onozaty/maildir-stats/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
