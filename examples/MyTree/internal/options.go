package tree

import "flag"

type Options struct {
	ShowFiles   bool
	ShowHidden  bool
	Color       bool
	ShowDate    bool
	ShowTime    bool
	Human       bool
	SortBy      string
	Reverse     bool
	MaxDepth    int
	FollowLinks bool
	Parallel    int
}

func ParseFlags() Options {
	showFiles := flag.Bool("f", false, "show files")
	showHidden := flag.Bool("a", false, "show hidden")
	color := flag.Bool("c", true, "enable ANSI colors")

	// 时间相关（独立）：
	// -d 显示日期（YYYY-MM-DD）
	// -t 显示时间（HH:MM:SS）
	showDate := flag.Bool("d", false, "show date (YYYY-MM-DD)")
	showTime := flag.Bool("t", false, "show time (HH:MM:SS)")

	human := flag.Bool("H", true, "human-readable sizes")
	sortBy := flag.String("S", "name", "sort by: name|size|time")
	reverse := flag.Bool("r", false, "reverse sort order")
	maxDepth := flag.Int("L", -1, "max depth")
	follow := flag.Bool("l", false, "follow symlinks")
	parallel := flag.Int("p", 8, "parallelism")

	flag.Parse()

	return Options{
		ShowFiles:   *showFiles,
		ShowHidden:  *showHidden,
		Color:       *color,
		ShowDate:    *showDate,
		ShowTime:    *showTime,
		Human:       *human,
		SortBy:      *sortBy,
		Reverse:     *reverse,
		MaxDepth:    *maxDepth,
		FollowLinks: *follow,
		Parallel:    *parallel,
	}
}
