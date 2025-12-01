<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { EditorView, keymap, placeholder } from '@codemirror/view'
import { EditorState, EditorSelection } from '@codemirror/state'
import { markdown } from '@codemirror/lang-markdown'
import { standardKeymap, history, historyKeymap } from '@codemirror/commands'
import { syntaxHighlighting, defaultHighlightStyle } from '@codemirror/language'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { UploadAttachment } from '../../wailsjs/go/main/App'

const emit = defineEmits(['save', 'cancel', 'command'])
const editorRef = ref(null)
const previewContent = ref('')
const isUploading = ref(false)
let view = null

// Command Suggestions
const showCommandSuggestions = ref(false)
const selectedSuggestionIndex = ref(0)
const suggestions = [
  { label: '/today', desc: '查看今日日志 (List)' },
  { label: '/week', desc: '导出本周日志 (Export)' },
  { label: '/month', desc: '导出本月日志 (Export)' },
  { label: '/list', desc: '查看今日日志 (同 /today)' }
]

const props = defineProps({
  isOpeningFile: Boolean
})

// Configure marked
marked.setOptions({
  gfm: true,
  breaks: true
})

// Helper: Wrap selection with open/close chars
const toggleWrapper = (open, close = open) => {
  return (view) => {
    const { state, dispatch } = view
    const changes = state.changeByRange(range => {
      const changes = []
      if (range.empty) {
        // Insert open+close and put cursor in middle
        changes.push({ from: range.from, insert: open + close })
      } else {
        // Wrap around
        changes.push({ from: range.from, insert: open })
        changes.push({ from: range.to, insert: close })
      }
      
      return {
        changes,
        range: EditorSelection.range(
            range.anchor + open.length,
            range.head + open.length
        )
      }
    })
    dispatch(changes)
    return true
  }
}

// Helper: Set or Toggle Heading
const setHeading = (level) => {
  return (view) => {
    const { state, dispatch } = view
    const changes = state.changeByRange(range => {
      const line = state.doc.lineAt(range.from)
      const text = line.text
      const headingRegex = /^(#{1,6})\s/
      const match = text.match(headingRegex)
      const targetPrefix = '#'.repeat(level) + ' '
      
      let change
      if (match) {
        // Replace existing heading
        if (match[1].length === level) {
             change = { from: line.from, to: line.from + match[0].length, insert: '' }
        } else {
             change = { from: line.from, to: line.from + match[0].length, insert: targetPrefix }
        }
      } else {
        // Insert new heading
        change = { from: line.from, insert: targetPrefix }
      }

      return {
        changes: [change],
        range: range 
      }
    })
    dispatch(changes)
    return true
  }
}

// Helper: Toggle Line Prefix (Lists, Task)
const toggleLinePrefix = (prefix) => {
  return (view) => {
    const { state, dispatch } = view
    const changes = state.changeByRange(range => {
      const line = state.doc.lineAt(range.from)
      const text = line.text
      
      let change
      
      if (text.startsWith(prefix)) {
        // Remove
        change = { from: line.from, to: line.from + prefix.length, insert: '' }
      } else {
        // Insert
        change = { from: line.from, insert: prefix }
      }
      
      return {
        changes: [change],
        range // Cursor stays
      }
    })
    dispatch(changes)
    return true
  }
}

// Helper: Insert Link
const insertLink = (view) => {
    const { state, dispatch } = view
    const changes = state.changeByRange(range => {
        const text = state.sliceDoc(range.from, range.to)
        const insert = `[${text}]()`
        return {
            changes: [{ from: range.from, to: range.to, insert }],
            range: EditorSelection.range(range.from + text.length + 3, range.from + text.length + 3) // Cursor inside ()
        }
    })
    dispatch(changes)
    return true
}

const selectCommand = (index) => {
  if (!view) return
  const cmd = suggestions[index].label
  view.dispatch({
      changes: { from: 0, to: view.state.doc.length, insert: cmd + ' ' }
  })
  // Move cursor to end
  view.dispatch({ selection: { anchor: cmd.length + 1 } })
  showCommandSuggestions.value = false
}

// Custom keymap
const customKeymap = [
  { key: 'Tab', run: (view) => {
      view.dispatch(view.state.replaceSelection('    '))
      return true
  }},
  { key: 'ArrowUp', run: (view) => {
      if (showCommandSuggestions.value) {
          selectedSuggestionIndex.value = (selectedSuggestionIndex.value - 1 + suggestions.length) % suggestions.length
          return true
      }
      return false
  }},
  { key: 'ArrowDown', run: (view) => {
      if (showCommandSuggestions.value) {
          selectedSuggestionIndex.value = (selectedSuggestionIndex.value + 1) % suggestions.length
          return true
      }
      return false
  }},
  { key: 'Enter', run: (view) => {
      if (showCommandSuggestions.value) {
          // If enter pressed while suggestion visible, select it?
          // Or just let user execute if they typed it fully?
          // Let's auto-complete if they are selecting
          selectCommand(selectedSuggestionIndex.value)
          return true
      }
      
      view.dispatch(view.state.replaceSelection('\n'))
      return true
    }
  },
  { key: 'Ctrl-Enter', run: (view) => {
      const content = view.state.doc.toString()
      
      // Check for commands first
      const trimmed = content.trim()
      if (trimmed.startsWith('/')) {
          emit('command', trimmed)
          view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: '' } })
          showCommandSuggestions.value = false
          return true
      }
      
      emit('save', content)
      view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: '' } })
      return true
    }
  },
  { key: 'Shift-Enter', run: () => false }, 
  { key: 'Escape', run: () => { 
      if (showCommandSuggestions.value) {
          showCommandSuggestions.value = false
          return true
      }
      emit('cancel'); return true 
  } },
  
  // Formatting
  { key: 'Mod-b', run: toggleWrapper('**') }, // Bold
  { key: 'Mod-i', run: toggleWrapper('*') },  // Italic
  { key: 'Mod-u', run: toggleWrapper('<u>', '</u>') }, // Underline
  { key: 'Ctrl-Shift-x', run: toggleWrapper('~~') }, // Strikethrough (Win/Linux)
  { key: 'Cmd-Shift-x', run: toggleWrapper('~~') },  // Strikethrough (Mac)
  { key: 'Mod-e', run: toggleWrapper('`') }, // Inline Code
  { key: 'Mod-k', run: insertLink }, // Link
  
  // Block Formatting
  { key: 'Ctrl-Alt-s', run: (view) => { 
      view.dispatch(view.state.replaceSelection('\n---\n')); return true 
  }}, // HR
  
  // Headings
  { key: 'Ctrl-Alt-1', run: setHeading(1) },
  { key: 'Ctrl-Alt-2', run: setHeading(2) },
  { key: 'Ctrl-Alt-3', run: setHeading(3) },
  { key: 'Ctrl-Alt-4', run: setHeading(4) },
  { key: 'Ctrl-Alt-5', run: setHeading(5) },
  { key: 'Ctrl-Alt-6', run: setHeading(6) },
  
  // Lists
  { key: 'Ctrl-Shift-8', run: toggleLinePrefix('- ') }, // Unordered List
  { key: 'Ctrl-Shift-7', run: toggleLinePrefix('1. ') }, // Ordered List
  { key: 'Ctrl-Alt-t', run: toggleLinePrefix('- [ ] ') }, // Task List
]

// Minimal theme
const minimalTheme = EditorView.theme({
  "&": {
    backgroundColor: "transparent",
    fontSize: "1.2rem",
    fontFamily: "inherit",
    height: "100%" 
  },
  ".cm-content": {
    padding: "0",
    fontFamily: "inherit",
    caretColor: "auto"
  },
  ".cm-line": {
    padding: "0"
  },
  "&.cm-focused": {
    outline: "none"
  },
  ".cm-activeLine": {
    backgroundColor: "transparent"
  },
  ".cm-gutters": {
    display: "none" 
  },
  ".cm-scroller": {
    overflow: "auto"
  }
})

const focus = () => {
  if (view) {
    view.focus()
  }
}

defineExpose({ focus })

// Handle Paste Event
const handlePaste = async (event, view) => {
  const items = event.clipboardData?.items
  if (!items) return false

  let hasHandledContent = false

  for (const item of items) {
    if (item.type.indexOf('image') !== -1) {
      event.preventDefault()
      hasHandledContent = true
      const file = item.getAsFile()
      if (!file) continue
      
      if (file.size > 50 * 1024 * 1024) { // 50MB limit
             alert(`Image ${file.name} exceeds 50MB limit.`)
             continue
      }

      try {
        isUploading.value = true
        const arrayBuffer = await file.arrayBuffer()
        const uint8Array = new Uint8Array(arrayBuffer)
        const array = Array.from(uint8Array)
        
        const webPath = await UploadAttachment(array, file.name || 'image.png')
        
        const md = `![image](${webPath})`
        view.dispatch(view.state.replaceSelection(md))
      } catch (err) {
        console.error('Failed to upload image:', err)
      } finally {
        isUploading.value = false
      }
      return true
    }
  }
  
  if (event.clipboardData?.files?.length > 0 && !hasHandledContent) {
    event.preventDefault()
    isUploading.value = true
    
    const processFiles = async () => {
        try {
            for (const file of event.clipboardData.files) {
                const isImage = file.type.indexOf('image') !== -1
                
                if (file.size > 50 * 1024 * 1024) { 
                     alert(`File ${file.name} exceeds 50MB limit.`)
                     continue
                }
        
                try {
                  const arrayBuffer = await file.arrayBuffer()
                  const uint8Array = new Uint8Array(arrayBuffer)
                  const array = Array.from(uint8Array)
                  
                  const webPath = await UploadAttachment(array, file.name)
                  
                  let md = ''
                  if (isImage) {
                     md = `![${file.name}](${webPath})`
                  } else {
                     md = `[${file.name}](${webPath})`
                  }
                  
                  view.dispatch(view.state.replaceSelection(md + ' '))
                } catch (err) {
                  console.error('Failed to upload file:', err)
                }
            }
        } finally {
            isUploading.value = false
        }
    }
    
    processFiles()
    return true
  }
  
  return false 
}

onMounted(() => {
  if (!editorRef.value) return

  const updateListener = EditorView.updateListener.of((update) => {
    if (update.docChanged) {
      const content = update.state.doc.toString()
      // Update Preview
      if (content.trim()) {
        previewContent.value = DOMPurify.sanitize(marked(content))
      } else {
        previewContent.value = ''
      }

      // Check for Command
      if (content.startsWith('/')) {
          showCommandSuggestions.value = true
          // Optionally filter suggestions based on typing?
          // For now show all, simple MVP.
      } else {
          showCommandSuggestions.value = false
      }
    }
  })
  
  const domEventHandler = EditorView.domEventHandlers({
    paste: (event, view) => {
      const result = handlePaste(event, view)
      // Sync check logic similar to before, simplified here for brevity
      const items = event.clipboardData?.items
      if (!items) return false
      
      let isCustom = false
      for (const item of items) {
          if (item.type.indexOf('image') !== -1) {
              isCustom = true
              break
          }
      }
      if (!isCustom && event.clipboardData?.files?.length > 0) {
          isCustom = true
      }
      
      if (isCustom) {
          handlePaste(event, view) 
          return true 
      }
      
      return false 
    }
  })

  const state = EditorState.create({
    doc: '',
    extensions: [
      history(),
      keymap.of([...customKeymap, ...standardKeymap, ...historyKeymap]),
      markdown(),
      syntaxHighlighting(defaultHighlightStyle), 
      minimalTheme,
      placeholder('Type your thought...'),
      EditorView.lineWrapping,
      updateListener,
      domEventHandler
    ]
  })

  view = new EditorView({
    state,
    parent: editorRef.value
  })
})

onBeforeUnmount(() => {
  if (view) {
    view.destroy()
  }
})
</script>

<template>
  <div class="input-area-wrapper">
    <div class="editors-row">
      <div ref="editorRef" class="editor-container" :class="{ 'has-preview': !!previewContent }"></div>
      
      <!-- Live Preview Area (Right column) -->
      <div v-if="previewContent" class="live-preview markdown-body" v-html="previewContent"></div>
    </div>

    <!-- Command Suggestions Popover -->
    <div v-if="showCommandSuggestions" class="command-suggestions">
        <div 
            v-for="(cmd, index) in suggestions" 
            :key="index" 
            class="suggestion-item"
            :class="{ active: index === selectedSuggestionIndex }"
            @click="selectCommand(index)"
        >
            <span class="cmd-label">{{ cmd.label }}</span>
            <span class="cmd-desc">{{ cmd.desc }}</span>
        </div>
    </div>

    <div class="hint">
      <span>Ctrl+↵ save</span>
      <span>/ command</span>
      <span>↵ newline</span>
      <span>Esc cancel</span>
      <span>Ctrl+B bold</span>
      <span>Ctrl+V paste img</span>
      <span v-if="isOpeningFile" class="status-opening">Opening file...</span>
      <span v-if="isUploading" class="status-opening">Uploading...</span>
    </div>
  </div>
</template>

<style scoped>
.input-area-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
  flex: 1; /* Allow it to grow */
  min-height: 0; /* Important for nested flex scrolling */
  position: relative;
}

.editors-row {
  display: flex;
  flex-direction: row;
  gap: 10px;
  /* Fill available height */
  flex: 1; 
  min-height: 0;
}

.editor-container {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  /* If no preview, it takes full width */
  transition: width 0.2s ease;
}

/* Optional: Visual separator or border when split? */
.editor-container.has-preview {
  border-right: 1px solid rgba(0,0,0,0.1);
  padding-right: 5px;
}

.live-preview {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  background: rgba(0,0,0,0.02);
  border-radius: 4px;
  padding: 5px 10px;
  font-size: 0.95rem; /* Match editor font size somewhat */
}

/* Ensure placeholder color matches design */
:deep(.cm-placeholder) {
  color: rgba(128, 128, 128, 0.6);
}

.hint {
  display: flex;
  gap: 12px;
  font-size: 0.75rem;
  color: rgba(128, 128, 128, 0.8);
  flex-wrap: wrap;
  margin-top: auto; /* Push to bottom */
}

.status-opening {
  color: #007acc;
  font-weight: bold;
  animation: pulse 1.5s infinite;
  margin-left: auto; /* Push to the right */
}

.command-suggestions {
    position: absolute;
    bottom: 30px; /* Above hints */
    left: 0;
    width: 300px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 6px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    z-index: 100;
    overflow: hidden;
}

.suggestion-item {
    padding: 8px 12px;
    display: flex;
    justify-content: space-between;
    cursor: pointer;
    border-bottom: 1px solid #eee;
}

.suggestion-item:last-child {
    border-bottom: none;
}

.suggestion-item.active {
    background-color: #f0f7ff;
}

.cmd-label {
    font-weight: bold;
    color: #007acc;
}

.cmd-desc {
    color: #888;
    font-size: 0.85rem;
}

@media (prefers-color-scheme: dark) {
    .command-suggestions {
        background: #2d2d2d;
        border-color: #444;
    }
    .suggestion-item {
        border-bottom-color: #444;
    }
    .suggestion-item.active {
        background-color: #3d3d3d;
    }
    .cmd-label {
        color: #4da6ff;
    }
    .cmd-desc {
        color: #aaa;
    }
}

@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

/* Markdown Styles (Copied from Timeline for consistency) */
:deep(.markdown-body p) {
  margin: 0; 
}
:deep(.markdown-body ul), :deep(.markdown-body ol) {
  margin: 0;
  padding-left: 20px;
}
:deep(.markdown-body code) {
  background-color: rgba(0,0,0,0.05);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
  font-size: 0.9em;
}
</style>
