# CodeQL扫描工具

## 简介

这是一个基于CodeQL的代码安全扫描工具，提供了图形化界面，支持对Java、Python、JavaScript和Go语言项目进行安全漏洞扫描。该工具可以使用默认的安全规则库，也支持配置自定义的规则库，方便进行定制化的安全检查。

## 环境要求

- Windows操作系统
- CodeQL CLI工具
- Java项目扫描需要Maven环境（如果要扫描Java项目）

## 安装步骤

1. 下载并安装CodeQL CLI工具：
   - 访问 [GitHub CodeQL CLI](https://github.com/github/codeql-cli-binaries/releases) 下载最新版本
   - 解压到本地目录，例如：`D:\Sec\CodeQL\codeql`

2. 如需扫描Java项目，请安装Maven：
   - 访问 [Apache Maven](https://maven.apache.org/download.cgi) 下载最新版本
   - 解压到本地目录，例如：`D:\Maven\apache-maven-3.6.0`
   - 配置环境变量（工具会自动处理）

## 配置说明

首次运行工具时，需要配置以下内容：

1. **CodeQL路径**：
   - 选择CodeQL可执行文件路径（codeql.exe）
   - 例如：`D:\Sec\CodeQL\codeql\codeql.exe`

2. **规则库路径**（可选）：
   - 可以启用自定义规则库
   - 选择CodeQL规则库所在目录
   - 例如：`D:\Sec\CodeQL\codeql-codeql-cli-v2.20.5\java\ql\src\Security`

3. **Maven路径**（扫描Java项目必需）：
   - 选择Maven可执行文件路径（mvn.cmd）
   - 例如：`D:\Maven\apache-maven-3.6.0\bin\mvn.cmd`

4. **临时工作目录**：
   - 选择用于存储临时文件的目录
   - 建议选择空间充足的目录

## 使用方法

1. 启动程序后，首先完成上述配置并点击"保存配置"按钮

2. 选择要扫描的项目：
   - 点击"选择目录"按钮
   - 选择要扫描的代码项目根目录

3. 选择项目语言：
   - 从下拉菜单中选择项目所用的编程语言
   - 支持Java、Python、JavaScript和Go

4. 开始扫描：
   - 点击"开始扫描"按钮
   - 等待扫描完成
   - 结果将显示在界面下方的文本框中
   - 同时会在项目目录生成CSV格式的详细报告

## 注意事项

1. 首次扫描前请确保所有配置正确，特别是CodeQL和Maven的路径

2. 扫描Java项目时：
   - 确保配置了正确的Maven路径
   - 项目应该是标准的Maven项目结构

3. 扫描过程可能需要一定时间，取决于：
   - 项目代码量
   - 选择的规则数量
   - 系统配置

4. 建议定期更新CodeQL规则库，以获取最新的安全检查规则

## 常见问题

1. **创建数据库失败**
   - 检查CodeQL路径是否正确
   - 确认选择的项目目录包含源代码
   - Java项目需要检查Maven配置
   - 请删除目录中的数据库文件试试

2. **分析失败**
   - 检查规则库路径是否正确
   - 确认选择的语言与项目实际使用的语言匹配

3. **扫描结果为空**
   - 检查是否选择了正确的规则库
   - 确认项目代码是否符合扫描要求

## 技术支持

如果在使用过程中遇到问题，请检查：

1. 配置文件（config.json）是否正确


## 许可说明

本工具仅用于安全测试和代码审计，请勿用于非法用途。使用本工具进行测试时，应确保已获得相应的授权。