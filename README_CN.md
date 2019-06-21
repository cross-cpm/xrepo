[EN](README.md) | 中文

xrepo 用于管理大型项目外部仓库引用

## 快速开始

### 安装

```shell
go get -u -v github.com/cross-cpm/xrepo
```

### 使用

创建 externals.json 文件，内容格式如下所示：

```json
{
    "https://github.com/felixqin/miniboost.git" : {
        "branch": "master",
        "ref": "HEAD",
        "targets": {
            "./": ["./externals/miniboost"]
        }
    },
    "https://github.com/felixqin/zbuild.git" : {
        "branch": "master",
        "ref": "HEAD",
        "targets": {
            "./": ["./externals/zbuild"]
        }
    }
}
```

在 externals.json 同级目录下，执行以下命令，检出所有文件中定义的外部依赖仓库：

```shell
xrepo checkout
```

此例中，目录树结果如下：

```text
<workdir>
    |- externals
    |   |- miniboost
    |   |- zbuild
    |- externals.json
```

## 更多命令

显示工具使用说明

```shell
xrepo help
```

全部外部仓库更新到最新版本

```shell
xrepo pull
```

推送全部外部仓库

```shell
xrepo push
```

查看全部外部仓库状态

```shell
xrepo status
```

列出外部仓库当前版本号和文件定义的版本号

```shell
xrepo rev list
```

将外部仓库当前版本号写入定义文件中

```shell
xrepo rev save
```
