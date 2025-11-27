import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { SaveNote, HideWindow, GetRecentNotes, OpenDailyNote } from '../../wailsjs/go/main/App'
import { WindowSetSize, EventsOn } from '../../wailsjs/runtime/runtime'

// State Constants
export const ViewState = {
  DEFAULT: 'default',
  TIMELINE: 'timeline'
}

export const ModalState = {
  NONE: 'none',
  COMMAND_PALETTE: 'command_palette'
}

export const ActivityState = {
  IDLE: 'idle',
  LOCKED: 'locked', // Focus debounce
  SAVING: 'saving',
  OPENING: 'opening'
}

const COLLAPSED_HEIGHT = 600
const EXPANDED_HEIGHT = 1000
const WIDTH = 800

export function useApp() {
  // State Machine
  const appState = reactive({
    view: ViewState.DEFAULT,
    modal: ModalState.NONE,
    activity: ActivityState.IDLE
  })

  const inputRef = ref(null)
  const recentNotes = ref([])
  let resetEventCancel = null

  // Computed Helpers
  const isTimelineVisible = computed(() => appState.view === ViewState.TIMELINE)
  const isCommandPaletteVisible = computed(() => appState.modal === ModalState.COMMAND_PALETTE)
  const isOpeningFile = computed(() => appState.activity === ActivityState.OPENING)
  const isSaving = computed(() => appState.activity === ActivityState.SAVING)
  const isLocked = computed(() => appState.activity === ActivityState.LOCKED)

  // Actions
  const loadNotes = async () => {
    try {
      recentNotes.value = await GetRecentNotes() || []
    } catch (error) {
      console.error('Failed to load notes:', error)
    }
  }

  const hideAndReset = async () => {
    await HideWindow()
  }

  const toggleTimeline = async () => {
    appState.view = appState.view === ViewState.DEFAULT ? ViewState.TIMELINE : ViewState.DEFAULT
    const height = isTimelineVisible.value ? EXPANDED_HEIGHT : COLLAPSED_HEIGHT
    WindowSetSize(WIDTH, height)
  }

  const toggleCommandPalette = () => {
    appState.modal = appState.modal === ModalState.NONE ? ModalState.COMMAND_PALETTE : ModalState.NONE
  }

  const closeCommandPalette = () => {
    appState.modal = ModalState.NONE
  }

  const handleSave = async (content) => {
    const trimmed = content.trim()
    if (!trimmed) return

    appState.activity = ActivityState.SAVING

    // Slash command: open
    if (trimmed === 'open') {
      appState.activity = ActivityState.OPENING
      try {
        await OpenDailyNote()
        await hideAndReset()
      } catch (error) {
        console.error('Failed to open editor:', error)
      } finally {
        appState.activity = ActivityState.IDLE
      }
      return
    }

    try {
      await SaveNote(content)
      await loadNotes()
      await hideAndReset()
    } catch (error) {
      console.error('Failed to save note:', error)
      appState.activity = ActivityState.IDLE
    } finally {
      if (appState.activity === ActivityState.SAVING) {
        setTimeout(() => { appState.activity = ActivityState.IDLE }, 500)
      }
    }
  }

  const handleEsc = () => {
    hideAndReset()
  }

  const openDailyNoteWrapper = async () => {
    appState.activity = ActivityState.OPENING
    try {
      await OpenDailyNote()
    } finally {
      setTimeout(() => { appState.activity = ActivityState.IDLE }, 1000)
    }
  }

  // Event Handlers
  const handleKeydown = async (e) => {
    // Open Today's Note with Ctrl+H
    if (e.ctrlKey && e.key.toLowerCase() === 'h') {
      e.preventDefault()
      openDailyNoteWrapper()
    }
    // Toggle command palette with Ctrl+P or Cmd+P
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'p') {
      e.preventDefault()
      toggleCommandPalette()
    }
  }

  const handleFocus = () => {
    appState.activity = ActivityState.LOCKED
    setTimeout(() => {
      if (appState.activity === ActivityState.LOCKED) {
         appState.activity = ActivityState.IDLE
      }
    }, 500)
  
    loadNotes()
    // Force focus back to input
    nextTick(() => {
      if (inputRef.value) {
        inputRef.value.focus()
      }
    })
  }

  const handleBlur = () => {
    if (appState.activity !== ActivityState.IDLE || isCommandPaletteVisible.value) return
    hideAndReset()
  }

  // Lifecycle
  onMounted(() => {
    loadNotes()
    WindowSetSize(WIDTH, COLLAPSED_HEIGHT)

    window.addEventListener('keydown', handleKeydown)
    window.addEventListener('focus', handleFocus)
    window.addEventListener('blur', handleBlur)

    resetEventCancel = EventsOn("app:reset", () => {
      nextTick(() => {
        if (inputRef.value) {
          inputRef.value.focus()
        }
      })
      loadNotes()
    })
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
    window.removeEventListener('focus', handleFocus)
    window.removeEventListener('blur', handleBlur)
    if (resetEventCancel) {
      resetEventCancel()
    }
  })

  return {
    appState,
    inputRef,
    recentNotes,
    // Computed
    isTimelineVisible,
    isCommandPaletteVisible,
    isOpeningFile,
    isSaving,
    isLocked,
    // Methods
    toggleTimeline,
    toggleCommandPalette,
    closeCommandPalette,
    handleSave,
    handleEsc,
    openDailyNote: openDailyNoteWrapper
  }
}

