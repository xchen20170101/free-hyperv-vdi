import axios from 'axios'
import ElementUI from "element-ui";
import router from './router';

const request = axios.create({
    //baseURL: 'http://127.0.0.1:8090',
    baseURL: '',
    timeout: 30000,
    withCredentials: true
})

// request 拦截器
// 可以自请求发送前对请求做一些处理
// 比如统一加token，对请求参数统一加密
request.interceptors.request.use(config => {

    config.headers['Content-Type'] = 'application/x-www-form-urlencoded;charset=utf-8';
    return config
}, error => {
    return Promise.reject(error)
});

// 移除了响应拦截器中的401状态码处理

export default request