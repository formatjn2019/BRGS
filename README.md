# BRGS

backup and restore of game saves

### 简介

设计之初是为了noita等类型的游戏做的存档备份工具，防止游戏机制或游戏bug导致的存档丢失造成玩家心态炸裂,同时能进行方便的SL操作。

设游戏存档目录为A,中转目录为B,程序基于fsnotify项目侦听系统底层的文件事件,可计算AB文件夹的差异，手动/定时将A中的文件同步到B,也可进行反向操作,仅同步文件差异，可用与快速的游戏SL操作

备份操作为将B中的全部文件存储为ZIP文件或者使用硬链接备份B中的文件(文件在C文件夹下首次存储为复制,二次及其以后为硬链接),C中存储为zip文件或文件夹,以时间戳命名，可由C中的备份还原至A,B->C可手动或自动操作，C->A由用户手动操作或不使用程序,使用操作系统的文件管理器进行操作。

包含一个简易的web界面，建议使用手机平板等设备操作，在windows下，同时启用web模式可自动打开浏览器，桌面版会比手机版多一个二维码界面，方便移动设备扫码接入。

### 使用方式

暂时还没弄完，不打算提供release

使用时ini文件必须存在

#### 自己编译

需要go环境和node环境

1. web界面
/web目录下

```shell
vite build
```

2. 主程序
   
windows

```cmd
./build.ps1
```
   
   

程序使用参数运行,不使用参数打开进入命令行界面，用于生成一个默认的csv文件。编辑并重命名该文件，可进行对应的配置
在同一界面读取该文件可生成或更新csv文件中对应配置的启动脚本~~也可自己手动改参数~~


### 版本

版本

v0.1.2 在桌面版网页增加了二维码显示，可由手机等移动设备直接扫码转到控制界面

v0.1.1 采用vue作为前端展示，命令行功能用于生成启动命令及手动模式
增加硬链接备份功能，减少重复文件空间

v0.0.2 命令行功能完成

v0.0.1 命令行基本功能完成



### 已知问题

1. 删除有子文件夹监控的目录系统会报错，在命令行模式可删除
2. 命令行模式，如果用通配符进行删除，会导致无法获得文件删除的通知
3. web端的二维码并未针对多网卡设备进行优化，所有若多网卡同时联网且二维码显示的地址与手机等设备不在同一网段下，会导致扫描二维码无法连接到web控制端，此时需要手动输入地址

### 其它

设计时是做了各个系统的兼容的，但目前mac因为设备问题没做，linux做了没测，建议只用windows的
使用哈希做文件的校验作为fsnotify的补充功能暂未实装


### Introduction

Originally designed as a save file backup tool for games like Noita, aiming to prevent player frustration caused by save file loss due to game mechanics or bugs, while also enabling convenient save/load operations.

The game save directory is denoted as A, the intermediary directory as B. The program, based on the fsnotify project, monitors low-level file events of the system, calculating the differences between folders A and B. It facilitates manual or scheduled synchronization of files from A to B, or vice versa, synchronizing only the file differences, enabling quick save/load operations in games.

Backup operations involve storing all files in B as a ZIP file or using hard links to backup files in B (files initially copied to directory C, subsequent backups are hard links). Files in C are named with timestamps and stored as ZIP files or directories. The synchronization between B and C can be done manually or automatically, while synchronization from C to A is a manual process by the user or not managed by the program, allowing the use of the operating system's file manager for operations.

Includes a simple web interface, recommended for use on mobile devices such as smartphones or tablets. In Windows, enabling web mode opens the browser automatically. The desktop version features a QR code interface in addition to the mobile version, facilitating access via scanning from mobile devices.

### Usage

Not fully implemented yet, no releases planned.

An ini file must be present for usage.

Compilation
Requires Go and Node environments.

Web Interface:
Located in the /web directory.
vite build

```shell
vite build
```

Main Program:
For Windows

```cmd
./build.ps1
```

Run the program with parameters. Without parameters, it opens in command-line mode, generating a default CSV file. Edit and rename this file for configuration. On the same interface, reading this file generates or updates startup scripts corresponding to the configurations in the CSV file. Alternatively, parameters can be manually adjusted.

### Version

v0.1.2: Added QR code display on the desktop version web page, allowing direct access to the control interface from mobile devices by scanning.

v0.1.1: Implemented Vue for frontend display, command-line functionality for generating startup commands and manual mode. Added hard link backup functionality to reduce redundant file space.

v0.0.2: Command-line functionality completed.

v0.0.1: Basic command-line functionality implemented.

Known Issues
1. Deleting a directory being monitored with subdirectories causes system errors; can be resolved in command-line mode.
2. In command-line mode, using wildcards for deletion prevents notification of file deletion. 
3. The QR code on the web interface is not optimized for devices with multiple network cards. If multiple network cards are active, and the address displayed on the QR code is not in the same network segment as mobile devices, scanning the QR code may fail to connect to the web control interface, requiring manual input of the address.

### Others

Compatibility with various systems was considered during design, but Mac compatibility has not been implemented due to device constraints. Linux compatibility has been implemented but not thoroughly tested; Windows usage is recommended. 
Hash-based file verification as a supplementary feature for fsnotify is not yet implemented.