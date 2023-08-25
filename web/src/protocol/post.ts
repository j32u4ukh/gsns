export interface GetPostRequest {
}

export interface GetPostResponse {
    posts: PostResponseData[],
}

export interface PostResponseData{
    userId: number,
    id: number,
    title: String,
    body: String,
}