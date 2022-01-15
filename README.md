# README.md

<h1 align="center">
  <img src="logo.png" alt="murphysec" width="150">
  <br><a href="https://murphysec.com" target="_blank">murphysec</a><br>
  <h4 align="center">一款专注于软件供应链安全的开源工具，包含开源组件依赖分析、漏洞检测及漏洞修复等功能。</h4>
</h1>
<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/murphysec/murphysec.svg?style=flat-square">
  <a href="https://github.com/murphysec/murphysec/releases/latest">
    <img src="https://img.shields.io/github/release/murphysec/murphysec.svg?style=flat-square">
  </a>
  <a href="https://github.com/murphysec/murphysec/blob/master/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/murphysec/murphysec?style=flat-square">
  </a>
  <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/murphysec/murphysec?style=flat-square">
  <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/murphysec/murphysec?style=social">

</p>

## 安装

### macOS

```shell
wget https://github.com/murphysec/murphysec/releases/latest/download/murphysec-darwin-amd64
chmod 0755 murphysec-darwin-amd64 && mv murphysec-darwin-amd64 /usr/bin/murphysec
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
wget https://github.com/murphysec/murphysec/releases/latest/download/murphysec-linux-amd64
chmod 0755 murphysec-linux-amd64 && mv murphysec-linux-amd64 /usr/bin/murphysec
```

## 配置

执行`murphysec auth login`完成身份验证

小范围公测中，[点此申请](https://murphysec.com/register)访问令牌

## 用法

```
murphysec: A software supply chain security inspection tool.

Usage:
  murphysec [flags]
  murphysec [command]

Available Commands:
  auth        manage the API token
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  scan        Scan open source vulnerabilities in project

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
