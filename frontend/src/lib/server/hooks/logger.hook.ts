import type { Handle } from "@sveltejs/kit";

function formatDate(d: Date): string {
	const pad = (n: number) => n.toString().padStart(2, "0");
	return `${d.getFullYear()}/${pad(d.getMonth() + 1)}/${pad(d.getDate())} - ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

export const handleLogger: Handle = async ({ event, resolve }) => {
	const start = performance.now();
	const startTime = new Date();

	const response = await resolve(event);

	const duration = `${(performance.now() - start).toFixed(2)}ms`.padStart(8, " ");
	const method = event.request.method.padEnd(8, " ");
	const url = event.url.pathname + event.url.search;
	const status = response.status;
	const clientIp = event.getClientAddress().padStart(15, " ");
	const timestamp = formatDate(startTime);

	console.log(`[BFF] ${timestamp} | ${status} | ${duration} | ${clientIp} | ${method} "${url}"`);

	return response;
};
