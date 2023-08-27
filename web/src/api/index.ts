import { ref } from 'vue';
import axios, { InternalAxiosRequestConfig } from 'axios';
import { Local } from '../utils/storage';

const service = ref(axios.create({
	// baseURL: import.meta.env.API_URL || "http://localhost:5001/",
    baseURL: "http://localhost:1023",
	timeout: 50000,
	headers: { 'Content-Type': 'application/json' },
}));

service.value.interceptors.request.use(
	(config:InternalAxiosRequestConfig) => {
		if (Local.get('accessToken')) {
            // with JWT token
			config.headers.Authorization = `Bearer ${Local.get('accessToken')}`;
		} else if (Local.get('refreshToken')) {
            config.headers.Authorization = `Bearer ${Local.get('refreshToken')}`;
        }
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

service.value.interceptors.response.use(
	(response) => {
		const res = response.data;
		if (res.code && res.code !== 200) {
			return Promise.reject(service.value.interceptors.response);
		} else {
			return response.data;
		}
	},
	async (error) => {
        const originalRequest = error.config;
        
		if (error.message.indexOf('timeout') != -1) {
            // timeout
		} else if (error.message == 'Network Error') {
            // network error
		} else {
            const res = error.response.data;
            console.log(res);
		}
		return Promise.reject(error);
	}
);
export default service.value;