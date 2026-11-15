import {
	CognitoIdentityProviderClient,
	AdminGetUserCommand,
	AdminDeleteUserCommand,
	UserNotFoundException,
	InitiateAuthCommand,
	AdminCreateUserCommand,
	AdminSetUserPasswordCommand,
} from "@aws-sdk/client-cognito-identity-provider";
import { userPoolId, clientId } from "./constants";
import { faker } from "@faker-js/faker";

const cognito = new CognitoIdentityProviderClient({
	region: "eu-west-1",
});

export const createCognitoUser = async () => {
	const u = {
		id: "",
		email: faker.internet.email(),
		name: faker.person.fullName(),
		password: `${faker.internet.password()}Aa1.`,
	}

	const user = await cognito.send(new AdminCreateUserCommand({
		UserPoolId: userPoolId,
		Username: u.email,
		UserAttributes: [
			{ Name: "email", Value: u.email },
			{ Name: "email_verified", Value: "true" },
			{ Name: "name", Value: u.name },
		],
		MessageAction: "SUPPRESS", // no welcome email
	}));
	const sub = user.User?.Attributes?.find((a) => a.Name === "sub")?.Value!
	u.id = sub

	await cognito.send(new AdminSetUserPasswordCommand({
		UserPoolId: userPoolId,
		Username: u.email,
		Password: u.password,
		Permanent: true,
	}));

	return u
}

export const getCognitoUser = async (id: string) => {
	console.log(`Getting user ${id} from Cognito ${userPoolId}`)
	const user = await cognito.send(new AdminGetUserCommand({
		UserPoolId: userPoolId,
		Username: id,
	})).catch((err) => {
		if (err instanceof UserNotFoundException) {
			return null
		}

		console.log(err)
		return null
	})

	return user
}

export const generateUserToken = async (email: string, password: string) => {
	const result = await cognito.send(
		new InitiateAuthCommand({
			ClientId: clientId,
			AuthFlow: "USER_PASSWORD_AUTH",
			AuthParameters: {
				USERNAME: email,
				PASSWORD: password,
			},
		}),
	);

	const token = result.AuthenticationResult?.AccessToken;
	if (!token) {
		throw new Error("No token returned");
	}

	return token;
}

export const deleteCognitoUser = async (id: string) => {
	console.log(`Deleting user ${id} from Cognito ${userPoolId}`)
	await cognito.send(new AdminDeleteUserCommand({
		UserPoolId: userPoolId,
		Username: id,
	})).catch((err) => {
		console.log(`Error deleting user ${id} from Cognito ${userPoolId}: ${err}`)
	})
}
