import { faker } from '@faker-js/faker'

export const generateRandomQuestion = () => {
	const category = [
		'General',
		'Competitions',
		'Drug Testing',
		'Live Stream',
		'Merchandise',
	]

	const post = {
		postType: "question",
		title: faker.lorem.words({ min: 5, max: 15 }),
		body: faker.lorem.paragraphs({ min: 3, max: 15 }),
		category: category[Math.floor(Math.random() * category.length)],
	}

	return post
}

export const generateRandomComment = () => {
	const comment = {
		body: faker.lorem.sentences({ min: 3, max: 15 }),
	}

	return comment
}
