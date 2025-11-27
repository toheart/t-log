<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { SaveNote, HideWindow, GetRecentNotes, OpenDailyNote } from '../wailsjs/go/main/App'
import { WindowSetSize } from '../wailsjs/runtime/runtime'
import InputArea from './components/InputArea.vue'
import Timeline from './components/Timeline.vue'

const inputRef = ref(null)
const recentNotes = ref([])
const isLocked = ref(false) // Lock auto-hide temporarily on focus
const isSaving = ref(false) // Lock auto-hide during save operation
const showTimeline = ref(false) // Toggle timeline visibility

const COLLAPSED_HEIGHT = 160
const EXPANDED_HEIGHT = 500
const WIDTH = 400

const loadNotes = async () => {
  try {
    recentNotes.value = await GetRecentNotes() || []
  } catch (error) {
    console.error('Failed to load notes:', error)
  }
}

const toggleTimeline = async () => {
  showTimeline.value = !showTimeline.value
  const height = showTimeline.value ? EXPANDED_HEIGHT : COLLAPSED_HEIGHT
  WindowSetSize(WIDTH, height)
}

const handleSave = async (content) => {
  const trimmed = content.trim()
  if (!trimmed) return

  isSaving.value = true // Prevent blur from triggering hide while we are saving

  // Slash command: open
  if (trimmed === 'open') {
    try {
      await OpenDailyNote()
      await HideWindow()
    } catch (error) {
      console.error('Failed to open editor:', error)
    } finally {
      setTimeout(() => { isSaving.value = false }, 500)
    }
    return
  }

  try {
    await SaveNote(content)
    await loadNotes() // Refresh timeline immediately
    await HideWindow()
  } catch (error) {
    console.error('Failed to save note:', error)
    isSaving.value = false // Reset if error, so blur works again
  } finally {
    if (isSaving.value) {
       setTimeout(() => { isSaving.value = false }, 500)
    }
  }
}

const handleEsc = () => {
  HideWindow()
}

// Handle global hotkeys
const handleKeydown = (e) => {
  // Toggle timeline with Ctrl+H (or F1 for simplicity)
  // Let's use Ctrl+H as requested context implies a shortcut
  if (e.ctrlKey && e.key.toLowerCase() === 'h') {
    e.preventDefault()
    toggleTimeline()
  }
}

window.addEventListener('keydown', handleKeydown)

window.addEventListener('focus', () => {
  isLocked.value = true
  setTimeout(() => {
    isLocked.value = false
  }, 500)

  isSaving.value = false

  loadNotes()
  // Force focus back to input
  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.focus()
    }
    // Reset to collapsed state on fresh open? 
    // Spec implies "minimal", so collapsing on re-open is good practice.
    // If user wants to keep it open, we can remove this.
    // For now, let's keep previous state or default to collapsed?
    // Defaulting to collapsed ensures "Quick Capture" focus.
    if (showTimeline.value) {
        showTimeline.value = false
        WindowSetSize(WIDTH, COLLAPSED_HEIGHT)
    }
  })
})

window.addEventListener('blur', () => {
  if (isLocked.value || isSaving.value) return
  HideWindow()
})

onMounted(() => {
  loadNotes()
  // Initial size
  WindowSetSize(WIDTH, COLLAPSED_HEIGHT)
})

</script>

<template>
  <main class="app-container" :class="{ expanded: showTimeline }">
    <InputArea 
      ref="inputRef"
      @save="handleSave"
      @cancel="handleEsc"
    />
    <div class="timeline-wrapper" v-show="showTimeline">
        <div class="divider"></div>
        <Timeline :notes="recentNotes" />
    </div>
    <div class="toggle-hint" @click="toggleTimeline" title="Toggle History (Ctrl+H)">
        <span v-if="showTimeline">▲</span>
        <span v-else>▼</span>
    </div>
  </main>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  background-color: transparent;
  height: 100vh;
  overflow: hidden;
}

#app {
  height: 100vh;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.9); /* Slightly more opaque */
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  padding: 20px;
  box-sizing: border-box;
  transition: height 0.3s ease;
  position: relative;
}

.timeline-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    margin-top: 10px;
}

.divider {
    height: 1px;
    background-color: rgba(0,0,0,0.1);
    margin-bottom: 10px;
}

.toggle-hint {
    position: absolute;
    bottom: 5px;
    left: 50%;
    transform: translateX(-50%);
    font-size: 10px;
    color: #999;
    cursor: pointer;
    padding: 2px 10px;
    opacity: 0.5;
}
.toggle-hint:hover {
    opacity: 1;
}

/* Dark mode support if needed */
@media (prefers-color-scheme: dark) {
  .app-container {
    background-color: rgba(30, 30, 30, 0.9);
    color: white;
  }
  .divider {
      background-color: rgba(255,255,255,0.1);
  }
}
</style>
