import { generateRandomQuestion, generateRandomComment } from './generators'
import { postUrl, commentUrl } from './constants'

export const createPost = async (token: string) => {
	console.log("Creating post")
	const response = await fetch(postUrl, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`,
		},
		body: JSON.stringify(generateRandomQuestion()),
	})

	const body = await response.json()
	if (!response.ok) {
		throw new Error(body.error?.message ?? `HTTP ${response.status}`)
	}
	return body.data
}

export const createComment = async (token: string, postId: string) => {
	console.log("Creating comment")
	const response = await fetch(commentUrl(postId), {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`,
		},
		body: JSON.stringify(generateRandomComment()),
	})

	const body = await response.json()
	if (!response.ok) {
		throw new Error(body.error?.message ?? `HTTP ${response.status}`)
	}
	return body.data
}
