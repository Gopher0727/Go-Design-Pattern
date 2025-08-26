# tree

- cmd/main.go
  
  解析命令行、构造 options、调用 internal/tree 的 Build + Render，负责进程层面的退出码/日志。

- internal/options.go
  
  定义 Options（showFiles、showHidden、color、maxDepth、sort、parallel、followLinks 等）和 ParseFlags()（便于测试）。

- internal/node.go
  
  组合模式的核心接口 FSNode 及实现 FileLeaf、DirComposite；只包含数据与 Render 接口（Render 接收 io.Writer 及 Options）。

- internal/builder.go
  
  负责文件系统遍历与构建节点树（并发控制、符号链接处理、环路检测、排序），提供 Build(root string, opts Options) (*DirComposite, error)。

- internal/render.go
  
  负责把节点树格式化输出（颜色、时间、大小、人类可读），Render(root *DirComposite, w io.Writer, opts Options)。

- internal/tree.go
  
  对外封装（例如 NewTree(opts) -> t.Build(root) -> t.Render(w)），便于单元测试与重用。

- internal/color
  
  封装颜色输出（支持 Windows via go-colorable），cmd 只传 opts.Color 即可。

- internal/fsutil
  
  实现跨平台 inode/文件标识提取与环路检测（可选）。