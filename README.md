# Quick Capture (t-log)

极简的桌面闪念笔记应用。无压记录，用完即走。

## 功能特性

- **极速唤起**: 全局快捷键 (`Ctrl + Alt + Space`) 毫秒级响应。
- **无干扰界面**: 类似 Windows 便签的无边框半透明窗口。
- **本地存储**: 笔记自动按 `YYYY/MM/YYYY-MM-DD.md` 归档到本地目录。
- **历史回溯**: 自动展示最近 3 天的笔记记录。
- **外部编辑**: 输入 `open` 命令一键调用系统编辑器打开当日笔记。

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
3. 输入内容，按 **`Enter`** 保存并隐藏。
4. 按 **`Shift + Enter`** 换行。
5. 输入 **`open`** 并回车，使用默认编辑器打开今日笔记文件。
6. 按 **`Esc`** 或点击窗口外部可取消并隐藏。

## 配置

首次运行后，程序会在运行目录下生成 `config.json` 文件：

```json
{
  "root_path": "C:\\Users\\YourName\\QuickNotes",
  "hotkey": "Alt+Space", 
  "history_days": 3
}
```

*注: 当前版本快捷键暂时固定为 Ctrl+Alt+Space，配置文件中的 hotkey 字段将在未来版本生效。*

## 构建

构建生产版本安装包:

```bash
wails build
```
