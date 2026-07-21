export function setCookie(name: string, value: string, minutes: number, secure: boolean) {
	const date = new Date();
	date.setTime(date.getTime() + minutes * 60 * 1000);
	let stringCookie = `${name}=${value}; expires=${date.toUTCString()}; path=/; SameSite=Lax`;

	if (secure) {
		stringCookie += '; Secure';
	}

	document.cookie = stringCookie;
}
