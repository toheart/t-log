<script setup>
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { ClipboardSetText } from '../../wailsjs/runtime'

const props = defineProps({
  notes: {
    type: Array,
    default: () => []
  },
  mode: {
    type: String,
    default: 'list', // 'list' | 'export'
    validator: (value) => ['list', 'export'].includes(value)
  }
})

// Configure marked
marked.setOptions({
  gfm: true,
  breaks: true
})

const renderMarkdown = (content) => {
  // Pre-process: transform "- [HH:MM]" into a styled badge, keeping the list structure
  // regex matches "- [10:30]" at start of line
  const processed = (content || '').replace(/^-\s+\[(\d{2}:\d{2})\]/gm, '- <span class="log-time">$1</span>')
  const rawHtml = marked(processed)
  return DOMPurify.sanitize(rawHtml)
}

// Helper to format date label
const getLabel = (date) => {
    const today = new Date().toISOString().split('T')[0]
    const yesterday = new Date(Date.now() - 86400000).toISOString().split('T')[0]
    if (date === today) return 'Today'
    if (date === yesterday) return 'Yesterday'
    return date
}

// Format notes for Export Mode
const exportText = computed(() => {
  return props.notes.map(note => {
    const label = getLabel(note.date)
    return `## ${label} (${note.date})\n\n${note.content}`
  }).join('\n\n')
})

const copyToClipboard = async () => {
  try {
    await ClipboardSetText(exportText.value)
    alert('已复制到剪贴板')
  } catch (err) {
    console.error('Failed to copy:', err)
    alert('复制失败')
  }
}
</script>

<template>
  <div class="context-panel">
    <!-- List Mode -->
    <div v-if="mode === 'list'" class="list-view">
      <div v-for="note in notes" :key="note.date" class="day-group">
        <div class="day-label">{{ getLabel(note.date) }}</div>
        <div class="entries">
            <div class="content markdown-body" v-html="renderMarkdown(note.content)"></div>
        </div>
      </div>
      <div v-if="notes.length === 0" class="empty-state">
        暂无记录
      </div>
    </div>

    <!-- Export Mode -->
    <div v-else-if="mode === 'export'" class="export-view">
      <div class="export-header">
        <span>导出预览</span>
        <button @click="copyToClipboard" class="copy-btn">复制全部</button>
      </div>
      <textarea readonly class="export-content" :value="exportText"></textarea>
    </div>
  </div>
</template>

<style scoped>
.context-panel {
  flex: 1;
  overflow-y: auto;
  padding-right: 5px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 6px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  min-height: 0; /* Important for flexbox scrolling */
}

/* List View Styles */
.day-group {
  margin-bottom: 20px;
  background: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.day-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: rgba(128, 128, 128, 0.8);
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
  padding-bottom: 8px;
}

.entries {
  font-size: 0.95rem;
  line-height: 1.5;
  color: #333;
}

.empty-state {
  text-align: center;
  color: #999;
  margin-top: 40px;
  font-size: 0.9rem;
}

/* Export View Styles */
.export-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 10px;
}

.export-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: #666;
}

.copy-btn {
  padding: 4px 12px;
  background-color: #007acc;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background-color 0.2s;
}

.copy-btn:hover {
  background-color: #005999;
}

.export-content {
  flex: 1;
  width: 100%;
  resize: none;
  padding: 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9rem;
  background-color: white;
  color: #333;
  outline: none;
}

/* Dark Mode */
@media (prefers-color-scheme: dark) {
  .day-group {
    background: #2d2d2d;
    box-shadow: 0 1px 3px rgba(0,0,0,0.2);
  }
  .entries {
    color: #eee;
  }
  .export-content {
    background-color: #2d2d2d;
    color: #eee;
    border-color: #444;
  }
  .context-panel {
    background: rgba(255, 255, 255, 0.02);
  }
}

/* Markdown Styles */
:deep(.markdown-body) {
  font-size: 0.95rem;
}
:deep(.markdown-body p) {
  margin-bottom: 0.5em;
}

/* Remove default bullets for the top-level list in markdown body, 
   assuming it's the log list */
:deep(.markdown-body > ul) {
  list-style-type: none;
  padding-left: 0;
}

:deep(.markdown-body > ul > li) {
  position: relative;
  padding-left: 0;
  margin-bottom: 12px;
}

/* Keep bullets/numbers for nested lists */
:deep(.markdown-body ul ul),
:deep(.markdown-body ul ol),
:deep(.markdown-body ol ul), 
:deep(.markdown-body ol ol) {
  margin-top: 4px;
  margin-bottom: 4px;
  padding-left: 24px;
}

/* Style the time badge */
:deep(.log-time) {
  display: inline-block;
  background-color: rgba(0, 0, 0, 0.06);
  color: #555;
  padding: 1px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.85em;
  margin-right: 8px;
  font-weight: 600;
  vertical-align: middle;
}

@media (prefers-color-scheme: dark) {
  :deep(.log-time) {
    background-color: rgba(255, 255, 255, 0.1);
    color: #ccc;
  }
}

:deep(.markdown-body code) {
  background-color: rgba(0,0,0,0.05);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
  font-size: 0.9em;
}
@media (prefers-color-scheme: dark) {
  :deep(.markdown-body code) {
    background-color: rgba(255,255,255,0.1);
  }
}
</style>
