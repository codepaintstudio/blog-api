<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Link from '@tiptap/extension-link'
import css from 'highlight.js/lib/languages/css'
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import html from 'highlight.js/lib/languages/xml'
import 'highlight.js/styles/github.css'

const lowlight = createLowlight(common)
lowlight.register('html', html)
lowlight.register('css', css)
lowlight.register('js', javascript)
lowlight.register('javascript', javascript)
lowlight.register('typescript', typescript)
lowlight.register('ts', typescript)

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
      defaultLanguage: 'javascript',
      HTMLAttributes: {
        class: 'hljs',
        spellcheck: 'true'
      }
    }),
    TaskList,
    TaskItem.configure({
      nested: true
    }),
    Link.configure({
      openOnClick: true,
      HTMLAttributes: {
        class: 'text-[#1f80ff] hover:text-[#1f80ff]/90 transition-colors'
      }
    })
  ]
})
</script>

<template>
  <div class="prose max-w-none text-gray-900">
    <editor-content :editor="editor" />
  </div>
</template>

<style scoped>
:root {
  --tw-prose-body: #333;
  --tw-prose-headings: var(--color-accent-emphasis);
  --tw-prose-links: #1f80ff;
  --tw-prose-bold: #333;
  --tw-prose-counters: var(--color-accent-emphasis);
  --tw-prose-bullets: var(--color-accent-emphasis);
  --tw-prose-hr: var(--color-accent-emphasis);
  --tw-prose-quotes: #333;
  --tw-prose-quote-borders: var(--color-accent-emphasis);
  --tw-prose-captions: var(--color-accent-emphasis);
  --tw-prose-code: #333;
  --tw-prose-pre-code: #333;
  --tw-prose-pre-bg: #f6f8fa;
  --tw-prose-th-borders: var(--color-accent-emphasis);
  --tw-prose-td-borders: var(--color-accent-emphasis);
}

.prose :where(code):not(:where([class~='not-prose'] *)) {
  background-color: #f6f8fa;
  padding: 0.125rem 0.25rem;
  font-size: 0.875em;
  border: 1px solid #ddd;
}

.prose :where(pre):not(:where([class~='not-prose'] *)) {
  background-color: #f6f8fa;
  border: 1px solid #ddd;
  overflow-x: auto;
  box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  padding: 1rem;
  margin: 1rem 0;
}

.prose div[data-type='codeBlock'] {
  position: relative;
  margin: 1rem 0;
}

/* 语言标签 */
.prose div[data-type='codeBlock']::before {
  content: attr(data-language);
  position: absolute;
  top: 0;
  right: 1rem;
  z-index: 1;
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-family: ui-monospace, monospace;
  line-height: 1;
  color: #666;
  text-transform: uppercase;
  background-color: #f6f8fa;
  border-bottom-left-radius: 0.25rem;
  border-bottom-right-radius: 0.25rem;
  border: 1px solid #ddd;
  border-top: none;
}

/* 代码块容器 */
.prose pre {
  position: relative;
  font-family:
    ui-monospace,
    SFMono-Regular,
    Menlo,
    Monaco,
    Consolas,
    Liberation Mono,
    Courier New,
    monospace;
  font-size: 0.875rem;
  line-height: 1.5;
}

.prose :where(pre code):not(:where([class~='not-prose'] *)) {
  background-color: transparent;
  border: none;
  padding: 0;
}

.prose :where(ul[data-type='taskList']):not(:where([class~='not-prose'] *)) {
  list-style: none;
  padding-left: 0;
}

.prose :where(ul[data-type='taskList'] li):not(:where([class~='not-prose'] *)) {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.prose :where(ul[data-type='taskList'] li > div):not(:where([class~='not-prose'] *)) {
  flex: 1;
}
</style>
