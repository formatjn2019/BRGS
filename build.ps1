# # 父文件夹
# $parentFolder = Split-Path  $pwd  -Parent

# 当前文件夹全路径
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
# 当前文件夹名
$currentFolderName = Split-Path  $scriptPath -Leaf
# 参数数组
$goosArr = "linux", "linux", "windows"
$archArr = "arm64", "amd64", "amd64"
$targetArr = "linux_arm64", "linux_amd64", "windows_amd64"
$suffixArr = "", "", ".exe"


for ($i = 0; $i -lt $goosArr.GetLength(0); $i++) {
    # 设置环境变量
    go env -w ("GOOS=" + ($goosArr[$i]))
    go env -w ("GOARCH=" + ($archArr[$i]))
    # 打包程序
    go build .

    # 拼接目标路径
    $targetPath = "$scriptPath\dist\" + $targetArr[$i]
    if (!(Test-Path $targetPath))
    {
        New-Item -ItemType Directory -Force -Path $targetPath
    }
    # 移动程序
    Move-Item ("$scriptPath\$currentFolderName" + ($suffixArr[$i])) ("$targetPath\$currentFolderName" + ($suffixArr[$i])) -force
}

