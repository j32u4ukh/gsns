import request from '../utils/request'
import * as ptc from "../protocol/post"

export const getPosts = () =>{
	return request<ptc.GetPostResponse>({
		url: '/posts',
		method: 'get',
		// data,
	});
}