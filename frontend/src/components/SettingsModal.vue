<template>
  <div v-if="isOpen" class="settings-modal-overlay" @click.self="close">
    <div class="settings-modal">
      <h2>Settings</h2>
      <div class="form-group">
        <label>Notes Save Path:</label>
        <div class="input-group">
          <input type="text" v-model="config.root_path" readonly />
          <button @click="browsePath">Browse</button>
        </div>
      </div>
      <div class="form-group">
        <label>Hotkey (Requires Restart):</label>
        <input type="text" v-model="config.hotkey" />
      </div>
      <div class="form-group">
        <label>History Days:</label>
        <input type="number" v-model.number="config.history_days" />
      </div>
      <div class="actions">
        <button @click="save">Save</button>
        <button @click="close" class="secondary">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
// Wails runtime bindings would be here, but we use window.go directly or imports
// Assuming window.go.main.App available globally

const isOpen = ref(false)
const config = ref({
  root_path: '',
  hotkey: '',
  history_days: 3
})

const open = async () => {
  try {
    const cfg = await window.go.main.App.GetConfig()
    config.value = { ...cfg } // clone
    isOpen.value = true
  } catch (err) {
    console.error("Failed to load config:", err)
  }
}

const close = () => {
  isOpen.value = false
}

const browsePath = async () => {
  try {
    const path = await window.go.main.App.SelectRootPath()
    if (path) {
      config.value.root_path = path
    }
  } catch (err) {
    console.error("Failed to browse path:", err)
  }
}

const save = async () => {
  try {
    await window.go.main.App.UpdateConfig(config.value)
    isOpen.value = false
  } catch (err) {
    console.error("Failed to save config:", err)
  }
}

// Event listener
const handleOpenSettings = () => {
  open()
}

onMounted(() => {
  if (window.runtime) {
    window.runtime.EventsOn("app:open-settings", handleOpenSettings)
  }
})

onUnmounted(() => {
  // Cleanup if needed
})

defineExpose({ open })
</script>

<style scoped>
.settings-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.settings-modal {
  background: #1e1e1e;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.5);
  color: #fff;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-size: 0.9em;
  color: #aaa;
}

.input-group {
  display: flex;
  gap: 5px;
}

input {
  width: 100%;
  padding: 8px;
  background: #2d2d2d;
  border: 1px solid #3d3d3d;
  color: #fff;
  border-radius: 4px;
}

button {
  padding: 8px 15px;
  background: #007acc;
  border: none;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background: #005999;
}

button.secondary {
  background: #3d3d3d;
}

button.secondary:hover {
  background: #4d4d4d;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}
</style>

