package tree

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

type FSNode interface {
	Name() string
	IsDir() bool
	ModTime() time.Time
	Size() int64
	Path() string
	Render(w io.Writer, prefix string, isLast bool, opts Options)
}

type FileLeaf struct {
	path string
	info fs.FileInfo
}

func NewFileLeaf(path string, info fs.FileInfo) *FileLeaf {
	return &FileLeaf{
		path: path,
		info: info,
	}
}

func (f *FileLeaf) Name() string { return f.info.Name() }

func (f *FileLeaf) IsDir() bool { return false }

func (f *FileLeaf) ModTime() time.Time { return f.info.ModTime() }

func (f *FileLeaf) Size() int64 { return f.info.Size() }

func (f *FileLeaf) Path() string { return f.path }

func (f *FileLeaf) Render(w io.Writer, prefix string, isLast bool, opts Options) {
	conn := "├── "
	if isLast {
		conn = "└── "
	}
	fmt.Fprintf(w, "%s%s%s\n", prefix, conn, renderMeta(f, opts))
}

type DirComposite struct {
	path     string
	info     fs.FileInfo
	children []FSNode
}

func NewDirComposite(path string, info fs.FileInfo) *DirComposite {
	return &DirComposite{
		path:     path,
		info:     info,
		children: nil,
	}
}

func (d *DirComposite) Name() string { return d.info.Name() }

func (d *DirComposite) IsDir() bool { return true }

func (d *DirComposite) ModTime() time.Time { return d.info.ModTime() }

func (d *DirComposite) Size() (ans int64) {
	// 递归计算目录大小（子节点构建完成后是只读的）
	for _, c := range d.children {
		ans += c.Size()
	}
	return ans
}

func (d *DirComposite) Path() string { return d.path }

func (d *DirComposite) Add(child FSNode) { d.children = append(d.children, child) }

// Render 方法用于以树状结构的形式渲染目录及其子节点。
// 它会根据当前节点是否为最后一个节点，动态调整前缀和连接符号。
//
// 参数:
//   - w: io.Writer 接口，用于输出渲染结果。
//   - prefix: 当前节点的前缀字符串，用于控制缩进。
//   - isLast: 布尔值，表示当前节点是否为其父节点的最后一个子节点。
//   - opts: Options 结构体，包含渲染选项（如是否启用颜色）。
//
// 渲染逻辑:
//   - 如果 prefix 为空，直接输出当前目录的路径。
//   - 否则，输出当前目录的名称（可选颜色）。
//   - 遍历子节点，递归调用 Render 方法，调整子节点的前缀和连接符号。
func (d *DirComposite) Render(w io.Writer, prefix string, isLast bool, opts Options) {
	conn := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		conn = "└── "
		childPrefix = prefix + "    "
	}

	// 根节点以绝对路径显示
	if prefix == "" {
		fmt.Fprintln(w, d.Path())
	} else {
		// 目录名字 使用统一的 ColorizeName 处理
		name := ColorizeName(d.Name(), true, opts)
		fmt.Fprintf(w, "%s%s%s\n", prefix, conn, name)
	}

	for i, c := range d.children {
		c.Render(w, childPrefix, i == len(d.children)-1, opts)
	}
}
