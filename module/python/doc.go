/*
Package python 用于 python 扫描

触发文件:
  - *.py (exclude conanfile.py)
  - pyproject.toml
  - requirements*

被触发路径下的 pyproject.toml 文件将被识别为一个模块；
被触发的路径下所有 requirements* 文件将被识别为一个模块。
如果以上两者未被发现，目录树遍历，所有名称为 venv 的文件夹将被跳过，解析包名称自所有 *.py 文件的 import 语句，将黑名单过滤后的包名与 pip list 输出合并作为一个模块。
*/
package python
