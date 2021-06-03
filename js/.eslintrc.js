module.exports = {
  env: {
    browser: true,
    es6: true,
    node: true,
    jest: true,
  },
  extends: ['plugin:@typescript-eslint/recommended', 'plugin:jest/recommended'],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly',
  },
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint', 'jest', 'prettier'],
  rules: {
    'prettier/prettier': [
      'error',
      {
        printWidth: 120,
        singleQuote: true,
        useTabs: false,
        tabWidth: 2,
        trailingComma: 'all',
      },
    ],
    '@typescript-eslint/no-namespace': 'off',
    '@typescript-eslint/camelcase': 0,
    'no-process-env': 'error',
  },
};
