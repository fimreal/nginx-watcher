package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/fimreal/goutils/sys"
	"github.com/fsnotify/fsnotify"
)

const (
	ngx_conf_dir_env_name     = "nginx_Confdir"
	default_ngx_conf_dir_name = "/etc/nginx/"
)

func main() {
	// 获取环境变量作为监控路径
	path, ok := os.LookupEnv(ngx_conf_dir_env_name)
	if !ok {
		path = default_ngx_conf_dir_name
	}
	watchConfigFile(path)
}

func reloadNginx() {
	pid, err := getPid()
	if err != nil {
		ezap.Error(err)
	}
	ezap.Debug("获取到 nginx pid: ", pid)

	ezap.Info("执行 nginx reload")
	err = sys.ProcessReload(pid)
	if err != nil {
		ezap.Errorf("向 nginx(%d) 进程发送 SUP 信号时出现错误: %s", pid, err)
	}
}

func getPid() (int, error) {
	ngxName := []string{"nginx", "openresty", "tengine"}
	for _, name := range ngxName {
		pid, err := sys.GetMasterPidByName(name)
		if err != nil {
			continue
		}
		return pid, nil
	}
	return 0, fmt.Errorf("not found nginx pid, maybe nginx is not running")
}

func watchConfigFile(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				ezap.Info("监控发现文件改动, ", event.Name, " ", event.Op, ", 等待 10s 后自动 reload 配置")
				time.Sleep(10 * time.Second)
				if event.Op&fsnotify.Create != fsnotify.Create {
					reloadNginx()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				ezap.Error(err)
			}
		}
	}()

	// for _, file := range path {
	// 	err = watcher.Add(file)
	// 	if err != nil {
	// 		ezap.Fatal("Add failed:", err)
	// 	}
	// }
	err = watcher.Add(path)
	if err != nil {
		ezap.Fatal("Add failed:", err)
	}

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ezap.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		}
		if info.IsDir() {
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			err = watcher.Add(path)
			if err != nil {
				return err
			}
			ezap.Info("添加监控目录: ", path)
		}
		return nil
	})

	<-done
}
