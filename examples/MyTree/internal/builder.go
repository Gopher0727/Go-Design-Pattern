package tree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type builder struct {
	opts Options
	sem  chan struct{} // 控制并发
}

func Build(root string, opts Options) (*DirComposite, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}

	rootNode := NewDirComposite(root, info)

	var wg sync.WaitGroup
	b := &builder{
		opts: opts,
		sem:  make(chan struct{}, max(1, opts.Parallel)),
	}
	wg.Go(func() {
		b.buildDir(rootNode, 0)
	})
	wg.Wait()

	return rootNode, nil
}

func (b *builder) buildDir(dir *DirComposite, depth int) {
	// 并发限制
	b.sem <- struct{}{}
	defer func() { <-b.sem }()
	if b.opts.MaxDepth >= 0 && depth > b.opts.MaxDepth {
		return
	}

	entries, err := os.ReadDir(dir.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir %s: %v\n", dir.path, err)
		return
	}

	children := make([]FSNode, 0, len(entries))
	var wg sync.WaitGroup
	for _, e := range entries {
		if !b.opts.ShowHidden && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		full := filepath.Join(dir.path, e.Name())
		var info fs.FileInfo
		if b.opts.FollowLinks {
			info, err = os.Stat(full)
		} else {
			info, err = e.Info()
		}
		if err != nil {
			// 忽略无法读取的条目
			continue
		}
		if info.IsDir() {
			// 如果达到最大深度，添加目录但不下探
			if b.opts.MaxDepth >= 0 && depth >= b.opts.MaxDepth {
				children = append(children, NewDirComposite(full, info))
				continue
			}
			sub := NewDirComposite(full, info)
			children = append(children, sub)

			wg.Add(1)
			go func(sd *DirComposite, d int) {
				defer wg.Done()
				b.buildDir(sd, d)
			}(sub, depth+1)
		} else {
			if !b.opts.ShowFiles {
				continue
			}
			children = append(children, NewFileLeaf(full, info))
		}
	}
	wg.Wait()

	keys := make([]string, 0, 3)
	for _, k := range strings.Split(b.opts.SortBy, "|") {
		if s := strings.TrimSpace(strings.ToLower(k)); s != "" {
			keys = append(keys, s)
		}
	}

	// 排序：目录优先，然后按选项排序（支持反序）
	sort.Slice(children, func(i, j int) bool {
		aNode, bNode := children[i], children[j]
		// 目录优先
		if aNode.IsDir() != bNode.IsDir() {
			return aNode.IsDir()
		}

		// 按 keys 依次比较，遇到不同值立即返回比较结果
		for _, key := range keys {
			switch key {
			case "size":
				if aNode.Size() != bNode.Size() {
					if b.opts.Reverse {
						return aNode.Size() > bNode.Size()
					}
					return aNode.Size() < bNode.Size()
				}
			case "time", "mtime", "modtime":
				if !aNode.ModTime().Equal(bNode.ModTime()) {
					if b.opts.Reverse {
						return aNode.ModTime().After(bNode.ModTime())
					}
					return aNode.ModTime().Before(bNode.ModTime())
				}
			default: // name 或未知的 key 当作 name
				an := strings.ToLower(aNode.Name())
				bn := strings.ToLower(bNode.Name())
				if an != bn {
					if b.opts.Reverse {
						return an > bn
					}
					return an < bn
				}
			}
		}

		// 所有 key 比较后相等：回退到 name 做稳定比较（不考虑 Reverse）
		return strings.ToLower(aNode.Name()) < strings.ToLower(bNode.Name())
	})

	dir.children = children
}
