/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  extends: [
    "plugin:vue/vue3-essential",
    'eslint:recommended',
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting'
  ],
  env: {
    // 'vue/setup-compiler-macros': true,
  },
  parserOptions: {
    ecmaVersion: 'latest'
  },
  // plugins: ['prettier'],
  rules: {
    '@typescript-eslint/no-unused-vars': ['error', { 'destructuredArrayIgnorePattern': '^_' }],
    'vue/multi-word-component-names': 'off',
    'vue/require-v-for-key': 'off',
    'prettier/prettier': [
      'error',
      {
        singleQuote: true,
        semi: true,
        useTabs: false,
        tabWidth: 2,
        trailingComma: 'all',
        printWidth: 120,
        bracketSpacing: true,
        arrowParens: 'avoid',
      },
    ],
  }
};
