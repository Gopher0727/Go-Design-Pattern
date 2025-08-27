package tree

import (
	"fmt"
	"io"
	"strings"
)

func Render(root *DirComposite, w io.Writer, opts Options) error {
	root.Render(w, "", true, opts)
	return nil
}

func renderMeta(n FSNode, opts Options) string {
	name := n.Name()
	parts := []string{name}

	// 单独处理 -d (ShowDate) 与 -t (ShowTime)，互不耦合
	if opts.ShowDate || opts.ShowTime {
		t := n.ModTime()
		var ts string
		switch {
		case opts.ShowDate && opts.ShowTime:
			ts = t.Format("2006-01-02 15:04:05")
		case opts.ShowDate:
			ts = t.Format("2006-01-02")
		case opts.ShowTime:
			ts = t.Format("15:04:05")
		}
		parts = append(parts, ts)
	}

	if !n.IsDir() {
		if opts.Human {
			parts = append(parts, humanSize(n.Size()))
		} else {
			parts = append(parts, fmt.Sprintf("%dB", n.Size()))
		}
	}

	// 统一颜色处理（目录/文件由 ColorizeName 决定）
	parts[0] = ColorizeName(name, n.IsDir(), opts)

	return strings.Join(parts, "  ")
}

func humanSize(s int64) string {
	const unit = 1024
	if s < unit {
		return fmt.Sprintf("%dB", s)
	}
	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	val := float64(s) / float64(div)
	return fmt.Sprintf("%.1f%cB", val, "KMGTPE"[exp])
}
