#钉钉告警脚本
  执行 go build dingtalk.go 生成可执行文件

#  运行方法：
  -c 指定配置文件
  -port  指定监听端口
  ./dingtalk -c configpath  -port=3306


#alertmanager配置：
  webhook_configs:
  - url: 'http://192.168.x.x:3306/sendmsg'
