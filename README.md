# Quick Capture (t-log)

极简的桌面闪念笔记应用。无压记录，用完即走。

## 功能特性

- **极速唤起**: 全局快捷键 (`Ctrl + Alt + Space`) 毫秒级响应。
- **无干扰界面**: 类似 Windows 便签的无边框半透明窗口。
- **Markdown 编辑**: 支持粗体、斜体、列表、标题等 Markdown 语法，实时预览。
- **图片附件**: 支持剪贴板直接粘贴图片 (`Ctrl + V`)，自动保存到本地。
- **快捷指令**: 输入 `/` 唤起指令菜单，快速查看今日、本周、本月日志。
- **命令面板**: `Ctrl + P` 唤起命令面板，支持全文搜索、打开特定日期笔记、设置等。
- **本地存储**: 笔记自动按 `YYYY/MM/YYYY-MM-DD.md` 归档到本地目录。
- **外部编辑**: 输入 `open` 或按 `Ctrl + H` 一键调用系统编辑器打开当日笔记。

## 快速开始

### 开发环境

1. 确保已安装 [Go 1.21+](https://go.dev/dl/) 和 [Node.js 18+](https://nodejs.org/)。
2. 安装 Wails: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`。
3. 启动开发模式:
   ```bash
   wails dev
   ```

### 使用说明

1. 启动后，应用会自动隐藏到后台。
2. 按下 **`Ctrl + Alt + Space`** 唤起窗口。
3. 输入内容，按 **`Ctrl + Enter`** 保存并隐藏。
4. 按 **`Enter`** 换行。
5. 输入 **`/`** 查看可用指令 (如 `/today` 查看今日日志)。
6. 按 **`Ctrl + P`** 打开命令面板，可搜索历史笔记。
7. 按 **`Ctrl + V`** 粘贴剪贴板中的图片。
8. 输入 **`open`** 或按 **`Ctrl + H`** 使用默认编辑器打开今日笔记文件。
9. 按 **`Esc`** 或点击窗口外部可取消并隐藏。

## 配置

首次运行后，程序会在运行目录下生成 `config.json` 文件。也可以通过命令面板 (`Ctrl + P`) -> `Settings` 进行可视化配置。

```json
{
  "root_path": "C:\\Users\\YourName\\QuickNotes",
  "hotkey": "Ctrl+Alt+Space", 
  "history_days": 3
}
```

## 构建

构建生产版本安装包:

```bash
wails build
```
