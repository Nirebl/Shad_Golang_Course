//go:build !solution

package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"gitlab.com/slon/shad-go/gitfame/gitstats"
	"gitlab.com/slon/shad-go/gitfame/gitutility"
	"gitlab.com/slon/shad-go/gitfame/printer"
	"gitlab.com/slon/shad-go/gitfame/utility"
)

func main() {
	var extensions []string
	var languages []string
	var exclude []string
	var restrictTo []string

	repository := pflag.String("repository", ".", "путь до Git репозитория")
	revision := pflag.String("revision", "HEAD", "указатель на коммит")
	orderBy := pflag.String("order-by", "lines", "ключ сортировки результатов")
	useCommitter := pflag.Bool("use-committer", false, "заменять автора на коммиттера")
	format := pflag.String("format", "tabular", "формат вывода")
	pflag.StringSliceVar(&extensions, "extensions", []string{}, "список расширений файлов")
	pflag.StringSliceVar(&languages, "languages", []string{}, "список языков программирования")
	pflag.StringSliceVar(&exclude, "exclude", []string{}, "набор Glob паттернов для исключения файлов")
	pflag.StringSliceVar(&restrictTo, "restrict-to", []string{}, "набор Glob паттернов для ограничения файлов")

	pflag.Parse()

	var repoInfo gitstats.RepoRequestData
	repoInfo.RepoHandler = gitutility.RepoHandler(*repository)
	repoInfo.Revision = *revision
	repoInfo.UseCommitter = *useCommitter
	repoInfo.Extensions = extensions
	repoInfo.Languages = languages
	repoInfo.Exclude = exclude
	repoInfo.RestrictTo = restrictTo

	stats, err := repoInfo.GetRepoStats()
	if err != nil {
		fmt.Println("Can't get statistics: ", err)
		os.Exit(1)
	}

	printer, err := printer.Select(*format)
	if err != nil {
		fmt.Println("Can't create printer: ", err)
		os.Exit(1)
	}

	stats = utility.SortUserInfos(stats, *orderBy)

	err = printer.Print(os.Stdout, stats)
	if err != nil {
		fmt.Println("Print error: ", err)
		os.Exit(1)
	}
}
