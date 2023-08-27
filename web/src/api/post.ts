import request from './index'
import * as ptc from "../protocol/post"

// /post/:postId
export const getThePost = (postId: number) =>{
	return request<ptc.GetThePostResponse>({
		url: `/post/${postId}`,
		method: 'get',
	});
}