module.exports = {
    settings: {
        react: {
            version: 'detect',
        },
    },
    env: {
        browser: true,
    },
    extends: [
        'eslint:recommended',
        'plugin:react/recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:react-hooks/recommended',
    ],
    parser: '@typescript-eslint/parser',
    parserOptions: {
        ecmaFeatures: {
            jsx: true
        },
        ecmaVersion: 'latest',
        sourceType: 'module',
    },
    plugins: [
        'react',
        '@typescript-eslint',
        'react-hooks'
    ],
    rules: {
        'no-console': ['error', { allow: ['log', 'info']}],
        'react/prop-types': 'off',
        'no-multiple-empty-lines': [
            'error',
            {
                'max': 1
            }
        ],
        'indent': [
            'error',
            4,
            {
                SwitchCase: 1,
                ImportDeclaration: 'first',
            }
        ],
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/no-var-requires': 'off',
        'react-hooks/rules-of-hooks': 'error',
        'react-hooks/exhaustive-deps': ['warn', {
            'additionalHooks': '(useMyCustomHook|useMyOtherCustomHook)'
        }],
        'react/display-name': 'off',
        'no-undef': 'off',
        'no-misleading-character-class': 'off',
        'no-empty': 'off',
        'no-useless-escape': 'off',
        '@typescript-eslint/no-unused-vars': 'off',
        'no-prototype-builtins': 'off',
        'react/no-string-refs': 'off',
        'react/no-unknown-property': 'off',
        '@typescript-eslint/no-this-alias': 'off',
        'no-constant-condition': 'off',
        '@typescript-eslint/ban-types': 'off',
        'react/no-find-dom-node': 'off',
        'prefer-rest-params': 'off',
        '@typescript-eslint/triple-slash-reference': 'off',
        'react/no-deprecated': 'off',
        'prefer-spread': 'off',
        'prefer-const': 'off',
        'no-fallthrough': 'off',
    }
}
