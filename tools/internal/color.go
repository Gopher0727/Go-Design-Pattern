package tree

import (
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// ColorizeName 根据扩展名或目录类型为 name 添加 ANSI 颜色（若 opts.Color 为 true）
func ColorizeName(name string, isDir bool, opts Options) string {
	if !opts.Color {
		return name
	}

	// 只有在平台支持 ANSI 时才添加转义序列
	if !supportsANSIONWindows() {
		return name
	}

	if isDir {
		return "\033[34m" + name + "\033[0m"
	}

	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go":
		return "\033[36m" + name + "\033[0m"
	case ".md", ".txt":
		return "\033[35m" + name + "\033[0m"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "\033[33m" + name + "\033[0m"
	case ".pdf":
		return "\033[32m" + name + "\033[0m"
	default:
		return name
	}
}

const enableVirtualTerminalProcessing = 0x0004

// 尝试为标准输出启用虚拟终端处理（VT），成功返回 true。
// 实现使用 kernel32 GetStdHandle/GetConsoleMode/SetConsoleMode。
func supportsANSIONWindows() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle := kernel32.NewProc("GetStdHandle")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	// STD_OUTPUT_HANDLE == -11
	h, _, _ := procGetStdHandle.Call(uintptr(^uintptr(10) + 1)) // -11 as uintptr
	if h == 0 {
		return false
	}

	var mode uint32
	r, _, _ := procGetConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return false
	}

	// 如果已经包含标志，则视为支持
	if mode&enableVirtualTerminalProcessing != 0 {
		return true
	}

	// 尝试设置标志
	mode |= enableVirtualTerminalProcessing
	r, _, _ = procSetConsoleMode.Call(h, uintptr(mode))
	return r != 0
}
