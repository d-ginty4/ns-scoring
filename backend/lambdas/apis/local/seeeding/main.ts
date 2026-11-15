import { clearPostsAndComments, checkUserExistsDynamo, createDynamoUser, deleteDynamoUser } from "./dynamodb";
import { clearPosts } from "./constants";
import { getCognitoUser, generateUserToken, createCognitoUser, deleteCognitoUser } from "./cognito";
import { createPost, createComment } from "./requests";
import { randomInt, sleep } from "./utils";

const main = async () => {
	if (clearPosts) {
		console.log("Clearing existing posts and comments")
		clearPostsAndComments()
	}

	// Create mock users
	const users = []
	console.log("Creating mock users")
	for (let i = 0; i < 5; i++) {
		const user = await createCognitoUser()
		await createDynamoUser(user.id, user.email, user.name)

		console.log(`Created user ${user.name} ${user.id} ${user.email}`)
		users.push(user)
	}

	console.log("Sleeping for 5 seconds to allow users to be created")
	await sleep(5000)

	// Create tokens for each of our mock users
	const tokens = []
	for (const user of users) {
		const u = await getCognitoUser(user.id)
		if (!u) {
			throw new Error("User does not exist in Cognito")
		}

		const userExists = await checkUserExistsDynamo(user.id)
		if (!userExists) {
			throw new Error("User does not exist in DynamoDB")
		}

		console.log(`Creating token for user ${user.name} ${user.id} ${user.email}`)
		tokens.push(await generateUserToken(user.email, user.password))
	}

	console.log("Creating 10 posts with random number of comments")
	for (let i = 0; i < 10; i++) {
		const post = await createPost(tokens[Math.floor(Math.random() * tokens.length)]).catch((err) => {
			console.log(err)
		})

		for (let j = 0; j < randomInt(0, 5); j++) {
			await createComment(tokens[Math.floor(Math.random() * tokens.length)], post.postId).catch((err) => {
				console.log(err)
			})
		}
	}


	console.log("Deleted the mock users")
	for (const user of users) {
		await deleteCognitoUser(user.id)
		await deleteDynamoUser(user.id)
	}
}

main()
