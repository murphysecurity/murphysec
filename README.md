# README.md

<h1 align="center">
  <img src="logo.png" alt="murphysec" width="150">
  <br><a href="https://murphysec.com" target="_blank">murphysec</a><br>
  <h4 align="center">一款专注于软件供应链安全的开源工具，包含开源组件依赖分析、漏洞检测及漏洞修复等功能。</h4>
</h1>
<p align="center">
  
</p>

![](media/16404924273208/16404949624351.jpg)

## 安装

### macOS

使用Homebrew安装

```shell
// TODO
```

### Windows

使用scoop安装

```
scoop bucket add murphysec https://github.com/murphysec/scoop-bucket
scoop update
scoop install murphysec
```

### Linux

```shell
// TODO
```

## 配置

执行`murphysec auth login`完成身份验证

`小范围内测中，点此获取访问令牌`[申请](https://murphysec.com/dasdsa)

## 用法

```
murphysec: An open source component security detection tool.

Usage:
  murphysec [flags]
  murphysec [command]

Available Commands:
  auth        manage the API token
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  scan

Flags:
      --color          colorize the output (default true)
  -h, --help           help for murphysec
      --token string   specify the API token
  -v, --verbose        show verbose log
      --version        output version information and exit

Use "murphysec [command] --help" for more information about a command.
```

## 开源协议

[Apache 2.0](LICENSE)