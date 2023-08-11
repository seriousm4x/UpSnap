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
					...themes['[data-theme=light]'],
					primary: '#55BCD9',
					secondary: '#D4A35C'
				}
			},
			{
				dark: {
					...themes['[data-theme=dark]'],
					primary: '#55BCD9',
					secondary: '#D4A35C',
					'base-300': '#101316'
				}
			},
			{
				cupcake: {
					...themes['[data-theme=cupcake]'],
					primary: '#55BCD9',
					secondary: '#D4A35C'
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
			'winter'
		],
		darkTheme: 'dracula'
	}
};
