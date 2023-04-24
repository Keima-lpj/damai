package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

var commands = map[string]string{
	"windows": "cmd",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var Version = "0.1.0"

// 打开系统默认的浏览器的对应链接
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	cmd := exec.Command(run, "/C", "start", uri)
	return cmd.Start()
}

func main() {
	uri := "https://www.damai.cn"
	// 打开大麦网
	err := Open(uri)
	if err != nil {
		fmt.Println("打开网页失败", uri, err)
	}
	// 点击演唱会按钮

	// 跳转到订单详情页面

	// 选中需要选中的观演人，这里的人数要和前面的数量保持一致

	// 模拟点击提交订单按钮
}
