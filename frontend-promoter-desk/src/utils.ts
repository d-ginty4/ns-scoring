export async function logHttpError(resp: Response, context: string): Promise<void> {
	let body = '';
	try {
		body = await resp.text();
	} catch {
		body = '(unable to read body)';
	}
	console.error(`${resp.statusText}(${resp.status}) - ${context}: ${body}`);
}
