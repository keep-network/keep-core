module.exports = {
	"parser": "babel-eslint",
	'env': {
		'browser': true,
		'es6': true,
	},
	'extends': [
		'eslint-config-keep',
		'plugin:react/recommended'
	],
	'globals': {
		'Atomics': 'readonly',
		'SharedArrayBuffer': 'readonly'
	},
	'parserOptions': {
		'ecmaFeatures': {
			'jsx': true,
		},
		'ecmaVersion': 2018,
		'sourceType': 'module'
	},
	'plugins': [
		'react'
	],
	'rules': {
		'react/prop-types': 0,
		'react/display-name': 0,
		'no-invalid-this': 0,
		'indent': ["error", 2, { "SwitchCase": 1 }],
		'no-unused-vars': ["error", { "ignoreRestSiblings": true }]
	}
}