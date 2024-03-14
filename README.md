# go_vedio
带h5前端界面下载器， 目前只支持m3u8, mp4视频

# 说明

编译

`build.sh`

## 两种使用方式

### 服务 

`serve.exe` 可以通过前端页面进行手动添加需要执行的任务url，也可以通过代理方式自动填充任务

启动时候会读取相对路径下的etc/task_serve.yaml文件作为配置文件

后端服务端口为 http://localhost:8888/

前端端口为 http://localhost:9999/


#### 开启代理方式
首先现将 `mitm.crt` 自签证书加入到根证书，担心安全性的话可自己生成一个，只需将证书名字修改为`mitm.crt`并放在同级目录下

默认监听地址为 http://localhost:10888



### 工具

`dv.exe`提供工具方式下载

`-help` 可查看指令

将需要下载的url以下列方式填入url.txt

```text
文件名称
https://www.test.com/1.m3u8
文件名称
https://www.test.com/1.mp4
...
```

## feature
1. ffmpeg 用cgo调用
2. 前端自动刷新下载进度