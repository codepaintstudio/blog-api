import js from '@eslint/js'
import globals from 'globals'
import pluginVue from 'eslint-plugin-vue'
import prettier from 'eslint-plugin-prettier'
import { defineConfig } from 'eslint/config'

export default defineConfig([
  { ignores: ['**/node_modules/**', '**/dist/**'] },
  pluginVue.configs['flat/essential'],
  {
    files: ['**/*.{js,mjs,cjs,vue}'],
    plugins: { js, prettier, vue: pluginVue },
    extends: ['js/recommended'],
    rules: {
      'prettier/prettier': 'error',
      'no-useless-catch': 'off',
      'vue/multi-word-component-names': 'off',
      quotes: ['error', 'single'],
      semi: ['error', 'never']
    }
  },
  { files: ['**/*.{js,mjs,cjs,vue}'], languageOptions: { globals: globals.browser } }
])
