<script setup>
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps({
  notes: {
    type: Array,
    default: () => []
  }
})

// Configure marked
marked.setOptions({
  gfm: true,
  breaks: true
})

const renderMarkdown = (content) => {
  const rawHtml = marked(content)
  return DOMPurify.sanitize(rawHtml)
}

// Group notes by date
const groupedNotes = computed(() => {
  const groups = {}
  props.notes.forEach(note => {
    // If date format is YYYY-MM-DD
    if (!groups[note.date]) {
      groups[note.date] = []
    }
    groups[note.date].push(note)
  })
  
  // Sort groups by date desc
  return Object.keys(groups)
    .sort((a, b) => new Date(b) - new Date(a))
    .map(date => {
        // Friendly date label
        const today = new Date().toISOString().split('T')[0]
        const yesterday = new Date(Date.now() - 86400000).toISOString().split('T')[0]
        let label = date
        if (date === today) label = 'Today'
        else if (date === yesterday) label = 'Yesterday'
        
        return {
            date,
            label,
            items: groups[date]
        }
    })
})
</script>

<template>
  <div class="timeline">
    <div v-for="group in groupedNotes" :key="group.date" class="day-group">
      <div class="day-label">{{ group.label }}</div>
      <div class="entries">
        <div v-for="(note, idx) in group.items" :key="idx" class="entry">
          <span class="time">{{ note.timestamp }}</span>
          <div class="content markdown-body" v-html="renderMarkdown(note.content)"></div>
        </div>
      </div>
    </div>
    <div v-if="notes.length === 0" class="empty-state">
      No recent notes
    </div>
  </div>
</template>

<style scoped>
.timeline {
  flex: 1;
  overflow-y: auto;
  padding-right: 5px; /* Space for scrollbar */
}

/* Custom Scrollbar */
.timeline::-webkit-scrollbar {
  width: 4px;
}
.timeline::-webkit-scrollbar-track {
  background: transparent;
}
.timeline::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
}
.timeline::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

.day-group {
  margin-bottom: 15px;
}

.day-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: rgba(128, 128, 128, 0.8);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.entry {
  display: flex;
  gap: 10px;
  padding: 4px 0;
  font-size: 0.95rem;
  line-height: 1.4;
  color: #333;
}

.time {
  font-family: monospace;
  color: #888;
  font-size: 0.85rem;
  min-width: 45px;
  margin-top: 2px;
  flex-shrink: 0;
}

.content {
  word-break: break-word;
  flex: 1;
}

.empty-state {
  text-align: center;
  color: #999;
  margin-top: 40px;
  font-size: 0.9rem;
}

@media (prefers-color-scheme: dark) {
  .entry {
    color: #eee;
  }
}

/* Markdown Styles - Minimalist */
:deep(.markdown-body) {
  font-size: 0.95rem;
}
/* Image Styles */
:deep(.markdown-body img) {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
  cursor: pointer;
  margin: 4px 0;
  display: block;
}
:deep(.markdown-body p) {
  margin: 0; /* Remove default p margins */
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
:deep(.markdown-body pre) {
  background-color: rgba(0,0,0,0.05);
  padding: 8px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 4px 0;
}
:deep(.markdown-body pre code) {
  background-color: transparent;
  padding: 0;
}
:deep(.markdown-body blockquote) {
  margin: 4px 0;
  padding-left: 10px;
  border-left: 3px solid #ddd;
  color: #777;
}
:deep(.markdown-body a) {
  color: #0066cc;
  text-decoration: none;
}
:deep(.markdown-body a:hover) {
  text-decoration: underline;
}
/* Dark mode adjustments for markdown */
@media (prefers-color-scheme: dark) {
  :deep(.markdown-body code) {
    background-color: rgba(255,255,255,0.1);
  }
  :deep(.markdown-body pre) {
    background-color: rgba(255,255,255,0.1);
  }
  :deep(.markdown-body blockquote) {
    border-left-color: #555;
    color: #aaa;
  }
  :deep(.markdown-body a) {
    color: #4da6ff;
  }
}
</style>
