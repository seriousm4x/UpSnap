export function toggleVisibility(el: HTMLInputElement) {
	if (el.type === 'password') {
		el.type = 'text';
	} else {
		el.type = 'password';
	}
}
