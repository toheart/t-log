<script setup>
import { useApp } from './composables/useApp'
import InputArea from './components/InputArea.vue'
import Timeline from './components/Timeline.vue'
import CommandPalette from './components/CommandPalette.vue'
import SettingsModal from './components/SettingsModal.vue'
import { ref } from 'vue'

const {
  inputRef,
  recentNotes,
  isTimelineVisible,
  isCommandPaletteVisible,
  isOpeningFile,
  handleSave,
  handleEsc,
  openDailyNote,
  closeCommandPalette
} = useApp()

const settingsRef = ref(null)
</script>

<template>
  <main class="app-container" :class="{ expanded: isTimelineVisible }">
    <InputArea 
      ref="inputRef"
      :is-opening-file="isOpeningFile"
      @save="handleSave"
      @cancel="handleEsc"
    />
    <div class="timeline-wrapper" v-show="isTimelineVisible">
        <div class="divider"></div>
        <Timeline :notes="recentNotes" />
    </div>
    <div class="toggle-hint" @click="openDailyNote" title="Open Today's Note (Ctrl+H)">
        <span>Open MD</span>
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
