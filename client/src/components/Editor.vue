<script setup>
import { onBeforeUnmount, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
const lowlight = createLowlight(common)
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Link from '@tiptap/extension-link'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  editable: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  autofocus: false,
  content: props.modelValue,
  editable: props.editable,
  extensions: [
    StarterKit,
    CodeBlockLowlight.configure({
      lowlight
    }),
    TaskList,
    TaskItem.configure({
      nested: true
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        class:
          'text-[#1f80ff] hover:text-[#1f80ff]/90 transition-colors'
      }
    })
  ],
  editorProps: {
    attributes: {
      class: 'prose max-w-none min-h-[200px] p-4 focus:outline-none'
    }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

watch(
  () => props.modelValue,
  (newValue) => {
    const isSame = newValue === editor.value.getHTML()
    if (!isSame) {
      editor.value.commands.setContent(newValue, false)
    }
  }
)

onBeforeUnmount(() => {
  editor.value.destroy()
})
</script>

<template>
  <div
    class="border border-[var(--color-accent-emphasis)] overflow-hidden bg-[var(--color-canvas-default)] shadow-sm"
  >
    <editor-content :editor="editor" class="prose max-w-none min-h-[200px] p-4" />
    <div
      class="bg-[var(--color-canvas-subtle)] border-t border-[var(--color-accent-emphasis)] p-2 flex items-center gap-2"
    >
      <button
        v-for="item in toolbarItems"
        :key="item.icon"
        @click="item.action"
        class="p-2 rounded hover:bg-[var(--color-canvas-default)] text-[var(--color-fg-default)]"
        :class="{ 'bg-[var(--color-canvas-default)]': item.isActive?.() }"
      >
        <span class="block w-5 h-5" v-html="item.icon"></span>
      </button>
    </div>
  </div>
</template>

<style scoped>
/* Base Editor Styles */
.ProseMirror:focus {
  outline: none;
}

.ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  color: var(--color-fg-subtle);
  float: left;
  height: 0;
  pointer-events: none;
}

/* Markdown Styles */
.prose pre {
  background-color: var(--color-canvas-subtle);
  border: 1px solid var(--color-border-default);
  border-radius: 0.25rem;
}

.prose code {
  background-color: var(--color-canvas-subtle);
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  border: 1px solid var(--color-border-default);
}

.prose pre code {
  background-color: transparent;
  border: none;
  padding: 0;
}

.prose blockquote {
  border-left: 4px solid var(--color-border-default);
  background-color: var(--color-canvas-subtle);
  padding-left: 1rem;
  padding-top: 0.25rem;
  padding-bottom: 0.25rem;
}

.prose hr {
  border-color: var(--color-border-default);
}

.prose table {
  border-collapse: collapse;
  border: 1px solid var(--color-border-default);
}

.prose th,
.prose td {
  border: 1px solid var(--color-border-default);
  padding: 0.5rem 1rem;
}

.prose thead {
  background-color: var(--color-canvas-subtle);
}

.prose ul[data-type='taskList'] {
  list-style: none;
  padding: 0;
}

.prose ul[data-type='taskList'] li {
  display: flex;
  align-items: flex-start;
}

.prose ul[data-type='taskList'] li > label {
  margin-right: 0.5rem;
  margin-top: 0.25rem;
}

.prose ul[data-type='taskList'] li > div {
  flex: 1;
}
</style>
