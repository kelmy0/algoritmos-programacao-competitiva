export function sanitizeTitle(input: string): string {
	return input.replace(/[<>]/g, "").replace(/\s+/g, " ").trim();
}

export function sanitizeHumanName(name: string): string {
	const clean = name.replace(/[^\p{L}\s.'-]/gu, "");

	const words = clean.replace(/\s+/g, " ").trim().split(" ");

	if (words.length === 0 || words[0] === "") return "";

	return words.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()).join(" ");
}

export function sanitizeUsername(username: string): string {
	const clean = username.replace(/[^\p{L}\p{N}_-]/gu, "");
	return clean.replace(/\s+/g, "").toLowerCase();
}

export function isValidEmail(email: string): boolean {
	const clean = email.trim().toLowerCase();
	const atIndex = clean.indexOf("@");
	const lastDotIndex = clean.lastIndexOf(".");

	return atIndex > 0 && lastDotIndex > atIndex + 1 && lastDotIndex < clean.length - 1;
}
