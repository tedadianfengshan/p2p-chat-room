# p2p-chat-room

开防火墙

```shell
firewall-cmd --zone=public --add-port=8002/tcp --permanent   # 开放8002端口
firewall-cmd --zone=public --remove-port=8002/tcp --permanent  #关闭8002端口
firewall-cmd --reload   # 配置立即生效
```