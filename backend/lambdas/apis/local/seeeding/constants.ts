const host = 'http://localhost:3000'
export const postUrl = `${host}/api/v1/post`
export const commentUrl = (postId: string) => `${host}/api/v1/post/${postId}/comment`
export const clearPosts = true
export const isLocal = true
export const userPoolId = "eu-west-1_CrIf3kzpW";
export const clientId = "1vrqh099jug02hbm95jipq4el";
export const userTableName = "ns-forum-backend-dev-users"
