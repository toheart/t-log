<script setup>
import { useApp } from './composables/useApp'
import InputArea from './components/InputArea.vue'
import ContextPanel from './components/ContextPanel.vue'
import CommandPalette from './components/CommandPalette.vue'
import SettingsModal from './components/SettingsModal.vue'
import { ref } from 'vue'

const {
  inputRef,
  recentNotes,
  contextPanelMode,
  isContextPanelVisible,
  isCommandPaletteVisible,
  isOpeningFile,
  handleSave,
  handleEsc,
  handleCommand,
  openDailyNote,
  closeCommandPalette
} = useApp()

const settingsRef = ref(null)
</script>

<template>
  <main class="app-container" :class="{ expanded: isContextPanelVisible }">
    <div class="main-column">
        <InputArea 
          ref="inputRef"
          :is-opening-file="isOpeningFile"
          @save="handleSave"
          @cancel="handleEsc"
          @command="handleCommand"
        />
        <div class="toggle-hint" @click="openDailyNote" title="Open Today's Note (Ctrl+H)">
            <span>Open MD</span>
        </div>
    </div>
    
    <div class="side-panel" v-if="isContextPanelVisible">
        <div class="divider-vertical"></div>
        <ContextPanel :notes="recentNotes" :mode="contextPanelMode" />
    </div>

    <CommandPalette 
      :visible="isCommandPaletteVisible"
      @close="closeCommandPalette"
    />

    <SettingsModal ref="settingsRef" />
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
  flex-direction: row; /* Layout changed to Row */
  height: 100%;
  background-color: rgba(255, 255, 255, 0.9); 
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  padding: 20px;
  box-sizing: border-box;
  transition: width 0.3s ease; /* Transition width instead of height */
  position: relative;
  overflow: hidden;
}

.main-column {
    flex: 1; /* Takes full width when side-panel is hidden */
    display: flex;
    flex-direction: column;
    height: 100%;
    min-width: 0; /* Prevents flex item from overflowing */
    position: relative;
}

.side-panel {
    width: 50%; /* Takes half width when visible */
    display: flex;
    flex-direction: row;
    height: 100%;
    min-width: 0;
    animation: slideIn 0.3s ease;
}

@keyframes slideIn {
    from { opacity: 0; transform: translateX(20px); }
    to { opacity: 1; transform: translateX(0); }
}

.divider-vertical {
    width: 1px;
    background-color: rgba(0,0,0,0.1);
    margin: 0 15px;
    flex-shrink: 0;
}

.toggle-hint {
    position: absolute;
    bottom: 0;
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
  .divider-vertical {
      background-color: rgba(255,255,255,0.1);
  }
}
</style>
