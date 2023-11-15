import themes from 'daisyui/src/theming/themes';

/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/daisyui/**/*.{js,jsx,ts,tsx}'],
	theme: {
		extend: {}
	},
	plugins: [require('daisyui')],
	daisyui: {
		themes: [
			{
				light: {
					...themes['light'],
					primary: '#55BCD9',
					secondary: '#D4A35C',
					'primary-content': '#ffffff'
				}
			},
			{
				dark: {
					...themes['dark'],
					primary: '#55BCD9',
					secondary: '#D4A35C',
					'base-300': '#12161b',
					'--rounded-btn': '1.9rem',
					'--tab-border': '2px',
					'--tab-radius': '.5rem'
				}
			},
			{
				cupcake: {
					...themes['cupcake'],
					primary: '#55BCD9',
					secondary: '#D4A35C',
					'primary-content': '#ffffff'
				}
			},
			'bumblebee',
			'emerald',
			'corporate',
			'synthwave',
			'retro',
			'cyberpunk',
			'valentine',
			'halloween',
			'garden',
			'forest',
			'aqua',
			'lofi',
			'pastel',
			'fantasy',
			'wireframe',
			'black',
			'luxury',
			'dracula',
			'cmyk',
			'autumn',
			'business',
			'acid',
			'lemonade',
			'night',
			'coffee',
			'winter',
			'dim',
			'nord',
			'sunset'
		]
	}
};
