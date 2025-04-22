<!-- ArticleContent.vue -->
<template>
  <div class="prose prose-invert max-w-none">
    <editor-content :editor="editor" />
  </div>
</template>

<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Link from '@tiptap/extension-link'

const lowlight = createLowlight(common)

const props = defineProps({
  content: {
    type: String,
    required: true
  }
})

const editor = useEditor({
  content: props.content,
  editable: false,
  extensions: [
    StarterKit,
    CodeBlockLowlight.configure({
      lowlight,
    }),
    TaskList,
    TaskItem.configure({
      nested: true,
    }),
    Link.configure({
      openOnClick: true,
      HTMLAttributes: {
        class: 'text-[var(--color-accent-fg)] hover:text-[var(--color-accent-emphasis)] transition-colors'
      }
    }),
  ]
})
</script>

<style>
/* Base styles for dark mode prose */
:root {
  --prose-body: var(--color-fg-default);
  --prose-headings: var(--color-fg-default);
  --prose-links: var(--color-accent-fg);
  --prose-bold: var(--color-fg-default);
  --prose-counters: var(--color-fg-subtle);
  --prose-bullets: var(--color-fg-subtle);
  --prose-hr: var(--color-border-default);
  --prose-quotes: var(--color-fg-default);
  --prose-quote-borders: var(--color-border-default);
  --prose-captions: var(--color-fg-subtle);
  --prose-code: var(--color-fg-default);
  --prose-pre-code: var(--color-fg-default);
  --prose-pre-bg: var(--color-canvas-subtle);
  --prose-th-borders: var(--color-border-default);
  --prose-td-borders: var(--color-border-default);
}

.prose {
  --tw-prose-body: var(--prose-body);
  --tw-prose-headings: var(--prose-headings);
  --tw-prose-links: var(--prose-links);
  --tw-prose-bold: var(--prose-bold);
  --tw-prose-counters: var(--prose-counters);
  --tw-prose-bullets: var(--prose-bullets);
  --tw-prose-hr: var(--prose-hr);
  --tw-prose-quotes: var(--prose-quotes);
  --tw-prose-quote-borders: var(--prose-quote-borders);
  --tw-prose-captions: var(--prose-captions);
  --tw-prose-code: var(--prose-code);
  --tw-prose-pre-code: var(--prose-pre-code);
  --tw-prose-pre-bg: var(--prose-pre-bg);
  --tw-prose-th-borders: var(--prose-th-borders);
  --tw-prose-td-borders: var(--prose-td-borders);
}

.prose :where(code):not(:where([class~="not-prose"] *)) {
  background-color: var(--color-canvas-subtle);
  border-radius: 0.25rem;
  padding: 0.125rem 0.25rem;
  font-size: 0.875em;
  border: 1px solid var(--color-border-default);
}

.prose :where(pre):not(:where([class~="not-prose"] *)) {
  background-color: var(--color-canvas-subtle);
  border: 1px solid var(--color-border-default);
  border-radius: 0.25rem;
}

.prose :where(pre code):not(:where([class~="not-prose"] *)) {
  background-color: transparent;
  border: none;
  padding: 0;
}

.prose :where(ul[data-type="taskList"]):not(:where([class~="not-prose"] *)) {
  list-style: none;
  padding-left: 0;
}

.prose :where(ul[data-type="taskList"] li):not(:where([class~="not-prose"] *)) {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.prose :where(ul[data-type="taskList"] li > div):not(:where([class~="not-prose"] *)) {
  flex: 1;
}
</style>
