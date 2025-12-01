import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { SaveNote, HideWindow, GetRecentNotes, OpenDailyNote, OpenDateNote, GetDailyNotes } from '../../wailsjs/go/main/App'
import { WindowSetSize, EventsOn } from '../../wailsjs/runtime/runtime'

// State Constants
export const ViewState = {
  DEFAULT: 'default',
  TIMELINE: 'timeline', // Legacy name, mapped to ContextPanel
  CONTEXT_PANEL: 'context_panel'
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
// const EXPANDED_HEIGHT = 1000
const DEFAULT_WIDTH = 800
const EXPANDED_WIDTH = 1400 // Increase width for side-by-side

export function useApp() {
  // State Machine
  const appState = reactive({
    view: ViewState.DEFAULT,
    modal: ModalState.NONE,
    activity: ActivityState.IDLE,
    contextMode: 'list' // 'list' | 'export'
  })

  const inputRef = ref(null)
  const recentNotes = ref([])
  let resetEventCancel = null

  // Computed Helpers
  const isContextPanelVisible = computed(() => appState.view === ViewState.CONTEXT_PANEL)
  const isCommandPaletteVisible = computed(() => appState.modal === ModalState.COMMAND_PALETTE)
  const isOpeningFile = computed(() => appState.activity === ActivityState.OPENING)
  const isSaving = computed(() => appState.activity === ActivityState.SAVING)
  const isLocked = computed(() => appState.activity === ActivityState.LOCKED)
  const contextPanelMode = computed(() => appState.contextMode)

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
    appState.view = ViewState.DEFAULT 
    WindowSetSize(DEFAULT_WIDTH, COLLAPSED_HEIGHT)
  }

  const toggleContextPanel = (mode = 'list') => {
    if (appState.view === ViewState.CONTEXT_PANEL && appState.contextMode === mode) {
        appState.view = ViewState.DEFAULT
        WindowSetSize(DEFAULT_WIDTH, COLLAPSED_HEIGHT)
    } else {
        appState.view = ViewState.CONTEXT_PANEL
        appState.contextMode = mode
        // Expand width instead of height
        WindowSetSize(EXPANDED_WIDTH, COLLAPSED_HEIGHT)
    }
  }

  const toggleCommandPalette = () => {
    appState.modal = appState.modal === ModalState.NONE ? ModalState.COMMAND_PALETTE : ModalState.NONE
  }

  const closeCommandPalette = () => {
    appState.modal = ModalState.NONE
  }

  // Command Handler
  const handleCommand = async (cmd) => {
      const command = cmd.toLowerCase().trim()
      
      if (command === '/today' || command === '/list') {
          const today = new Date().toISOString().split('T')[0]
          try {
              recentNotes.value = await GetDailyNotes(today, today) || []
              toggleContextPanel('list')
          } catch (err) {
              console.error(err)
          }
          return
      }
      
      if (command === '/week') {
          const now = new Date()
          const day = now.getDay() || 7 
          if (day !== 1) now.setHours(-24 * (day - 1)) 
          const startOfWeek = now.toISOString().split('T')[0]
          const today = new Date().toISOString().split('T')[0]
          
          try {
              recentNotes.value = await GetDailyNotes(startOfWeek, today) || []
              toggleContextPanel('export')
          } catch (err) {
              console.error(err)
          }
          return
      }

      if (command === '/month') {
        const now = new Date()
        const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0]
        const today = new Date().toISOString().split('T')[0]
        
        try {
            recentNotes.value = await GetDailyNotes(startOfMonth, today) || []
            toggleContextPanel('export')
        } catch (err) {
            console.error(err)
        }
        return
    }
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
    
    // Slash command: open YYYY-MM-DD
    if (trimmed.startsWith('open ')) {
        const dateStr = trimmed.substring(5).trim()
        if (/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) {
            appState.activity = ActivityState.OPENING
            try {
                await OpenDateNote(dateStr)
                await hideAndReset()
            } catch (error) {
                console.error('Failed to open date note:', error)
            } finally {
                appState.activity = ActivityState.IDLE
            }
            return
        }
    }

    try {
      await SaveNote(content)
      if (isContextPanelVisible.value && appState.contextMode === 'list') {
          const today = new Date().toISOString().split('T')[0]
          recentNotes.value = await GetDailyNotes(today, today) || []
      } else {
          await loadNotes()
      }
      
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
    if (isContextPanelVisible.value) {
        appState.view = ViewState.DEFAULT
        WindowSetSize(DEFAULT_WIDTH, COLLAPSED_HEIGHT)
        return
    }
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
    if (e.ctrlKey && e.key.toLowerCase() === 'h') {
      e.preventDefault()
      openDailyNoteWrapper()
    }
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
  
    nextTick(() => {
      if (inputRef.value) {
        inputRef.value.focus()
      }
    })
  }

  // Lifecycle
  onMounted(() => {
    WindowSetSize(DEFAULT_WIDTH, COLLAPSED_HEIGHT)

    window.addEventListener('keydown', handleKeydown)
    window.addEventListener('focus', handleFocus)

    resetEventCancel = EventsOn("app:reset", () => {
      nextTick(() => {
        if (inputRef.value) {
          inputRef.value.focus()
        }
      })
      appState.view = ViewState.DEFAULT
      WindowSetSize(DEFAULT_WIDTH, COLLAPSED_HEIGHT)
    })
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
    window.removeEventListener('focus', handleFocus)
    if (resetEventCancel) {
      resetEventCancel()
    }
  })

  return {
    appState,
    inputRef,
    recentNotes,
    isContextPanelVisible,
    contextPanelMode,
    isCommandPaletteVisible,
    isOpeningFile,
    isSaving,
    isLocked,
    toggleContextPanel,
    toggleCommandPalette,
    closeCommandPalette,
    handleSave,
    handleEsc,
    handleCommand,
    openDailyNote: openDailyNoteWrapper
  }
}
