import { logHttpError } from "../utils"
import type { Response } from "../types"

export const host = "http://localhost:3000/api/v1"

type RequestOptions = {
	method?: "GET" | "POST" | "PUT" | "DELETE"
	body?: unknown
}

async function send(path: string, options: RequestOptions = {}) {
	const method = options.method ?? "GET"
	const resp = await fetch(`${host}${path}`, {
		method,
		headers: options.body === undefined ? undefined : { "Content-Type": "application/json" },
		body: options.body === undefined ? undefined : JSON.stringify(options.body),
	})

	if (!resp.ok) {
		await logHttpError(resp, `Failed to ${method} ${path}`)
		return null
	}

	return resp
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T | null> {
	const resp = await send(path, options)
	if (!resp) return null

	const payload: Response<T> = await resp.json()
	return payload.data
}

export async function requestOk(path: string, options: RequestOptions = {}): Promise<boolean> {
	const resp = await send(path, options)
	return resp !== null
}

export function resourcePath(base: string, ...ids: Array<string | number>): string {
	const encoded = ids.map((id) => encodeURIComponent(String(id)))
	return [base, ...encoded].join("/")
}
