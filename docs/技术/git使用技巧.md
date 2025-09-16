- 合并提交
``` shell
git rebase -i <start_commit_id>
# 使用 `s` 开头，表示将该提交合入前面的提交
# 使用 `f` 开头，表示将该提交合入前面的提交，且不保留该提交的message
git push --force-with-lease
```
