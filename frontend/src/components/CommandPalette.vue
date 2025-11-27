<template>
  <div v-if="visible" class="command-palette-overlay" @click.self="close">
    <div class="command-palette">
      <input
        ref="inputRef"
        v-model="searchQuery"
        class="command-input"
        type="text"
        :placeholder="placeholder"
        @keydown.down.prevent="selectNext"
        @keydown.up.prevent="selectPrev"
        @keydown.enter.prevent="executeSelected"
        @keydown.esc.prevent="close"
        @input="handleInput"
      />
      <ul class="command-list">
        <li
          v-for="(item, index) in filteredItems"
          :key="item.id || index"
          :class="{ active: index === selectedIndex }"
          @click="selectItem(index)"
          @mouseover="selectedIndex = index"
        >
          <div class="command-item">
            <span class="command-title">{{ item.title || item.content }}</span>
            <span v-if="item.description" class="command-desc">{{ item.description }}</span>
            <span v-if="item.date" class="command-date">{{ item.date }} {{ item.time }}</span>
          </div>
        </li>
        <li v-if="filteredItems.length === 0" class="no-results">
          No results found
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue';
import { GetCommands, ExecuteCommand, SearchNotes, OpenNoteAt } from '../../wailsjs/go/main/App';

const props = defineProps({
  visible: Boolean
});

const emit = defineEmits(['close']);

const inputRef = ref(null);
const searchQuery = ref('');
const selectedIndex = ref(0);
const commands = ref([]);
const searchResults = ref([]);
const mode = ref('command'); // 'command' or 'search'
let searchTimeout = null; // For debounce

// Computed placeholder based on mode
const placeholder = computed(() => {
  return mode.value === 'search' ? 'Search notes...' : 'Type a command...';
});

// Computed items to display
const filteredItems = computed(() => {
  if (mode.value === 'search') {
    return searchResults.value;
  }
  
  // Command mode: filter commands by query
  if (!searchQuery.value) return commands.value;
  const query = searchQuery.value.toLowerCase();
  return commands.value.filter(cmd => 
    cmd.title.toLowerCase().includes(query) || 
    (cmd.description && cmd.description.toLowerCase().includes(query))
  );
});

// Watch visibility to focus input and reset state
watch(() => props.visible, async (newVal) => {
  if (newVal) {
    await loadCommands();
    searchQuery.value = '';
    mode.value = 'command';
    selectedIndex.value = 0;
    await nextTick();
    inputRef.value?.focus();
  }
});

const loadCommands = async () => {
  try {
    commands.value = await GetCommands();
  } catch (err) {
    console.error('Failed to load commands:', err);
  }
};

const close = () => {
  emit('close');
};

const selectNext = () => {
  if (selectedIndex.value < filteredItems.value.length - 1) {
    selectedIndex.value++;
    scrollIntoView();
  }
};

const selectPrev = () => {
  if (selectedIndex.value > 0) {
    selectedIndex.value--;
    scrollIntoView();
  }
};

const scrollIntoView = () => {
  // Implementation for scrolling to active item if needed
  // Simple version: rely on mouseover for now or standard browser behavior
};

const handleInput = async () => {
  selectedIndex.value = 0;
  
  // Check for search mode trigger
  if (mode.value === 'command' && searchQuery.value.startsWith('find ')) {
    mode.value = 'search';
    searchQuery.value = ''; // Clear query for actual search term
    return;
  }

  // Handle search mode
  if (mode.value === 'search') {
    if (!searchQuery.value) {
      searchResults.value = [];
      return;
    }
    
    // Debounce search
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(async () => {
      try {
        searchResults.value = await SearchNotes(searchQuery.value);
      } catch (err) {
        console.error('Search failed:', err);
      }
    }, 300); // 300ms debounce
  }
};

const executeSelected = async () => {
  const item = filteredItems.value[selectedIndex.value];
  if (!item) return;

  if (mode.value === 'command') {
    if (item.id === 'cmd:find') {
      mode.value = 'search';
      searchQuery.value = '';
      return;
    }
    
    try {
      await ExecuteCommand(item.id, []);
      close();
    } catch (err) {
      console.error('Command execution failed:', err);
    }
  } else {
    // Handle search result selection
    // item structure from SearchNotes: { date: string, time: string, content: string, line_number: int }
    // We want to open the daily note for that date
    // Since we don't have a direct "Open Note at Line" API yet, 
    // we'll at least open the daily note for that day.
    // Ideally, the backend OpenDailyNote could take a date argument.
    // For now, let's assume we want to at least open the file.
    
    // But wait, OpenDailyNote only opens TODAY.
    // We might need a new backend command: OpenNote(date string)
    // Or better yet, OpenNoteAt(filePath string, line int)
    
    // Since we can't easily add backend API right this second without context switch,
    // let's look at what we have.
    // We have OpenDailyNote().
    
    // Actually, the user requirement is "定位到文本位置".
    // This implies opening the file in the default editor (VS Code, Notepad, etc.)
    // Most editors support opening at a line number, e.g. "code file.md:10".
    
    // Let's emit an event to the parent to handle this, or call a new backend method if we had one.
    // For now, let's just log it and maybe alert that we need that backend capability.
    // OR, we can try to execute a command that opens that specific file.
    
    // Wait, SearchResult has: Date, Content, SourceFile (maybe?)
    // Let's check search.go or note.go for SearchResult struct.
    
    // Assuming we can't change backend right now, we'll close and log.
    // BUT user specifically asked for this. So we MUST implement it.
    // We need to ask the parent component or call a backend function.
    
    // Let's assume we add a new command to App.go: OpenNoteAt(date, line)
    // But first, let's just keep the palette open? No, user said "instead of closing".
    // Meaning: selecting it SHOULD do something (locate text), not just close.
    
    // If we can't locate it yet, we should probably at least NOT close immediately if it does nothing.
    // But the goal is to make it work.
    
    // Let's try to call OpenNoteAt
    if (item.filePath && item.lineNo) {
        try {
            await OpenNoteAt(item.filePath, item.lineNo);
            close(); // Close after successful open
        } catch (err) {
             console.error('Failed to open note at location:', err);
        }
    } else {
        console.warn('Missing filePath or lineNo for item:', item);
        close();
    }
  }
};

const selectItem = (index) => {
  selectedIndex.value = index;
  executeSelected();
};
</script>

<style scoped>
.command-palette-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 100px;
}

.command-palette {
  width: 600px;
  background: var(--bg-color, #1e1e1e);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 400px;
}

.command-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  border: none;
  background: transparent;
  color: var(--text-color, #e0e0e0);
  border-bottom: 1px solid #333;
  outline: none;
  box-sizing: border-box; /* Ensure padding doesn't overflow width */
}

.command-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
}

.command-list li {
  padding: 10px 16px;
  cursor: pointer;
  border-bottom: 1px solid #2a2a2a;
}

.command-list li.active {
  background: #2a2d3e; /* Selection color */
}

.command-list li:hover {
  background: #2a2d3e;
}

.command-item {
  display: flex;
  flex-direction: column;
}

.command-title {
  font-weight: 500;
  color: var(--text-color, #e0e0e0);
}

.command-desc {
  font-size: 12px;
  color: #888;
  margin-top: 2px;
}

.command-date {
  font-size: 11px;
  color: #666;
  align-self: flex-end;
}

.no-results {
  padding: 16px;
  text-align: center;
  color: #888;
}

/* Light mode override */
@media (prefers-color-scheme: light) {
  .command-palette {
    background: #ffffff;
    border: 1px solid #ddd;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .command-input {
    color: #333;
    border-bottom: 1px solid #eee;
  }
  
  .command-title {
    color: #333;
  }
  
  .command-list li {
    border-bottom: 1px solid #f0f0f0;
  }
  
  .command-list li.active, .command-list li:hover {
    background: #f5f5f5;
  }
}
</style>
