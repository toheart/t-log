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

const emit = defineEmits(['save', 'cancel'])
const editorRef = ref(null)
const previewContent = ref('')
const isUploading = ref(false)
let view = null

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
             // Toggle off if same level? Usually typora just keeps it or changes it.
             // Let's make it strictly "Set" per shortcut, maybe toggle if same?
             // Typora behavior: Ctrl+1 on H1 -> Body text.
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
        range: range // Keep relative position? range.map(change) handled by CM?
                     // Actually simpler to just return range. 
                     // For line changes, keeping cursor in place is default behavior of state.update
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
      
      // Check if line already starts with prefix (ignoring whitespace)
      // Using simple startsWith for MVP, or regex if needed
      // Task list: "- [ ] "
      // Unordered: "- " or "* "
      // Ordered: "1. " (regex needed for d.)
      
      let change
      // Handle dynamic regex for ordered list "1. " vs "2. " etc? 
      // Typora shortcut creates "1. ".
      
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

// Custom keymap
const customKeymap = [
  { key: 'Tab', run: (view) => {
      view.dispatch(view.state.replaceSelection('    '))
      return true
  }},
  { key: 'Ctrl-Enter', run: (view) => {
      const content = view.state.doc.toString()
      emit('save', content)
      view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: '' } })
      return true
    }
  },
  { key: 'Enter', run: (view) => {
      view.dispatch(view.state.replaceSelection('\n'))
      return true
    }
  },
  { key: 'Shift-Enter', run: () => false }, // Default behavior is fine, or map to newline too
  { key: 'Escape', run: () => { emit('cancel'); return true } },
  
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
    height: "100%" // Take full height of container
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
    display: "none" // Hide line numbers
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

  // Check if we have files to handle
  let hasHandledContent = false

  for (const item of items) {
    // Handle Images
    if (item.type.indexOf('image') !== -1) {
      event.preventDefault()
      hasHandledContent = true
      const file = item.getAsFile()
      if (!file) continue
      
      if (file.size > 50 * 1024 * 1024) { // 50MB limit
             alert(`Image ${file.name} exceeds 50MB limit.`)
             continue
      }

      // Upload logic
      try {
        isUploading.value = true
        const arrayBuffer = await file.arrayBuffer()
        const uint8Array = new Uint8Array(arrayBuffer)
        const array = Array.from(uint8Array)
        
        const webPath = await UploadAttachment(array, file.name || 'image.png')
        
        // Insert Markdown
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
  
  // Handle Files (Non-Image)
  // Check if there are files AND it's not just plain text content being pasted as file
  // Sometimes text is also available as file.
  // But usually event.clipboardData.files is empty for pure text.
  if (event.clipboardData?.files?.length > 0 && !hasHandledContent) {
    // Double check if it is text?
    // If I copy text "abc", files is usually empty.
    // If I copy a file from explorer, files has length 1.
    
    event.preventDefault()
    isUploading.value = true
    
    const processFiles = async () => {
        try {
            for (const file of event.clipboardData.files) {
                const isImage = file.type.indexOf('image') !== -1
                
                if (file.size > 50 * 1024 * 1024) { // 50MB limit
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
  
  // Default behavior for text
  return false
}

onMounted(() => {
  if (!editorRef.value) return

  const updateListener = EditorView.updateListener.of((update) => {
    if (update.docChanged) {
      const content = update.state.doc.toString()
      // Only render if not empty
      if (content.trim()) {
        previewContent.value = DOMPurify.sanitize(marked(content))
      } else {
        previewContent.value = ''
      }
    }
  })
  
  const domEventHandler = EditorView.domEventHandlers({
    paste: (event, view) => {
      // Wait, handlePaste is async. CodeMirror expects return bool immediately.
      // If we return false, default happens.
      // But handlePaste might decide later it was a file?
      // Actually handlePaste preventsDefault synchronously if it finds a file.
      // The async part is the upload.
      // So we can just call it.
      
      const result = handlePaste(event, view)
      // If handlePaste returns a Promise (because async), we can't use its return value directly for sync prevention.
      // But handlePaste calls event.preventDefault() internally if matches.
      // So we just return false generally? Or rely on event default prevented?
      
      // CodeMirror logic: "If a handler returns true, the event is assumed to be handled..."
      // We need to detect if we handled it.
      
      // Fix: handlePaste is async, so it returns a Promise.
      // We need to check items synchronously.
      // Let's make handlePaste synchronous in decision making, async in execution.
      // But `item.getAsFile()` is sync.
      // The issue is `await handlePaste` isn't possible here.
      
      // We can refactor handlePaste to NOT be async wrapper, but fire async operations.
      // But for now, let's just rely on the fact that we call preventDefault() inside.
      // However, if we return false (because it's a promise object), CM might try to handle it too?
      // If preventDefault is called, CM usually respects it.
      
      // Let's see. Promise object is truthy? Yes.
      // So if we return handlePaste(), we are returning true (Promise).
      // This prevents CM default paste! Even for text!
      // THIS IS THE BUG.
      
      // We must ONLY return true if we actually found a file/image.
      // But we can't know that from the Promise result immediately unless we peek.
      
      // Refactor: Check synchronously.
      const items = event.clipboardData?.items
      if (!items) return false
      
      let isCustom = false
      // Check for images
      for (const item of items) {
          if (item.type.indexOf('image') !== -1) {
              isCustom = true
              break
          }
      }
      // Check for files (if not image item found)
      if (!isCustom && event.clipboardData?.files?.length > 0) {
          isCustom = true
      }
      
      if (isCustom) {
          handlePaste(event, view) // Fire and forget async
          return true // Tell CM we handled it
      }
      
      return false // Tell CM to handle it (Text)
    }
  })

  const state = EditorState.create({
    doc: '',
    extensions: [
      history(),
      keymap.of([...customKeymap, ...standardKeymap, ...historyKeymap]),
      markdown(),
      syntaxHighlighting(defaultHighlightStyle), // Use default highligting
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

    <div class="hint">
      <span>Ctrl+↵ save</span>
      <span>↵ newline</span>
      <span>Esc cancel</span>
      <span>Ctrl+B bold</span>
      <span>Ctrl+I italic</span>
      <span>Ctrl+K link</span>
      <span>Ctrl+E code</span>
      <span>Ctrl+V paste img</span>
      <span>Ctrl+1..6 H1-6</span>
      <span>Ctrl+Shift+8 list</span>
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
