package help

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//获取main的绝对路径
func GetMainCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

//获取文件的相对路径
//path 相对与项目的相对路径
func GetFileAbsPath(path string) string {
	main_path := GetMainCurrPath()
	return filepath.Join(main_path, path)
}
