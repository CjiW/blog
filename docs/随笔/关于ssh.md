## SSH 密钥
- 生成密钥对
```shell
ssh-keygen -t ed25519 -C "your_email@example.com"
```

- 添加私钥
```shell
# 后台启动 ssh 代理
eval "$(ssh-agent -s)"
# 将私钥添加到 ssh-agent
ssh-add ~/.ssh/id_ed25519
```
