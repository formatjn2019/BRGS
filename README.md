# BRGS

backup and restore of game saves

能对文件夹进行备份与还原及压缩操作

适用于对游戏存档文件的自动备份

可用于存档文件的保存或者频繁的sl操作

版本
v0.1.1 采用vue作为前端展示，命令行功能用于生成启动命令及手动模式
增加硬链接备份功能，减少重复文件空间

v0.0.2 命令行功能完成

v0.0.1 命令行基本功能完成

### 已知问题

1. 删除有子文件夹监控的目录系统会报错，在命令行模式可删除
2. 命令行模式，如果用通配符进行删除，会导致无法获得文件删除的通知
