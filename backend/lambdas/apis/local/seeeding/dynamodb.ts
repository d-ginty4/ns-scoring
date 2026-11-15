import { ScanCommand, BatchWriteItemCommand, DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { DynamoDBDocumentClient, GetCommand, PutCommand, DeleteCommand } from "@aws-sdk/lib-dynamodb";
import { isLocal, userTableName } from "./constants";

const client = isLocal ? new DynamoDBClient({
	region: "local",
	endpoint: "http://localhost:8000",
	credentials: {
		accessKeyId: "local",
		secretAccessKey: "local",
	}
}) : new DynamoDBClient({});

// DynamoDB Document Client used dynamodb in dev
const docClient = DynamoDBDocumentClient.from(new DynamoDBClient({}));

export const clearPostsAndComments = async () => {
	const tableName = `ns-forum-backend-${isLocal ? "local" : "dev"}-main`

	await deleteAllItems(tableName, client)
}

async function deleteAllItems(tableName: string, client: DynamoDBClient) {
	let lastEvaluatedKey;

	do {
		const result = await client.send(
			new ScanCommand({
				TableName: tableName,
				ExclusiveStartKey: lastEvaluatedKey,

				ProjectionExpression: "#postId, #meta",
				ExpressionAttributeNames: {
					"#postId": "postId",
					"#meta": "meta",
				},
			})
		);

		const items = result.Items ?? [];

		for (let i = 0; i < items.length; i += 25) {
			const batch = items.slice(i, i + 25);

			await client.send(
				new BatchWriteItemCommand({
					RequestItems: {
						[tableName]: batch.map((item) => ({
							DeleteRequest: {
								Key: {
									postId: item.postId,
									meta: item.meta,
								},
							},
						})),
					},
				})
			);
		}

		lastEvaluatedKey = result.LastEvaluatedKey;
	} while (lastEvaluatedKey);
}

export const checkUserExistsDynamo = async (id: string) => {
	const { Item } = await docClient.send(new GetCommand({
		TableName: userTableName,
		Key: {
			id
		},
	}));

	if (!Item) {
		return false
	}
	return true
}

export const createDynamoUser = async (id: string, email: string, name: string) => {
	const isAdmin = Math.random() < 0.2

	if (isAdmin) {
		console.log(`Admin user ${name}`)
	}

	await docClient.send(new PutCommand({
		TableName: userTableName,
		Item: {
			id,
			name,
			email,
			isAdmin: isAdmin,
			title: isAdmin ? "Admin" : "",
			preferences: { receiveEmails: true },
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		},
	}));
}

export const deleteDynamoUser = async (id: string) => {
	console.log(`Deleting user ${id} from DynamoDB ${userTableName}`)
	await docClient.send(new DeleteCommand({
		TableName: userTableName,
		Key: {
			id
		},
	}));
}
