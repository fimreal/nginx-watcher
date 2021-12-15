# nginx-watcher

golang 编写，利用 github.com/fsnotify/fsnotify 来监控配置文件变化，向 nginx master 进程发送 SIGHUP 信号。

作用是方便 k8s 多副本实例自动重载配置，附带 dockerfile 和 yaml 例子，部署方便。

目前已改进的点包括：

- 使用 filepath.Walk 来递归添加监控目录
- 同时支持 nginx、openresty、tengine 进程监测，但只会触发同一 process namespace 中第一个启动的 nginx 进行 reload
- 可以通过环境变量配置需要监控的配置文件目录，默认使用 "/etc/nginx/"
- 即使是单次改动，fsnotify 监控触发事件也比较多，使用 chan 保证第一次监测到改动 10s 后执行一次 nginx reload 操作
- 添加 Dockerfile 和 k8s yaml 配置

待改进：

- 添加 github action 自动打包上传镜像到 dockerhub
- 过滤不敏感修改，减少 reload 操作
