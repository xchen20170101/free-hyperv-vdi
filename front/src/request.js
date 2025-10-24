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

// response 拦截器
// 处理响应数据，包括token失效的情况
let isRedirecting = false; // 防止重复跳转的标志

request.interceptors.response.use(
    response => {
        // 检查响应数据中的code字段
        const res = response.data;
        
        // 获取当前路由路径
        const currentPath = router.currentRoute.path;
        const isOnLoginPage = currentPath === '/' || currentPath === '/login';
        
        // 如果code为2（表示失败），检查是否需要跳转登录页
        if (res.code === 2 && res.msg && !isOnLoginPage) {
            // 定义会触发跳转登录页的错误消息（token相关错误）
            const tokenErrors = [
                'Token.Invalid',      // Token无效
                'Token.Expired',      // Token过期
                'Token.NotExist',     // Token不存在
                'Auth.Failed',        // 认证失败
                'Session.Expired',    // 会话过期
                'Unauthorized'        // 未授权
            ];
            
            // 检查是否为token失效错误
            const isTokenError = tokenErrors.some(keyword => res.msg.includes(keyword));
            
            // 如果是token失效错误，跳转到登录页
            if (isTokenError && !isRedirecting) {
                isRedirecting = true; // 设置跳转标志
                
                // 清除本地存储和会话存储
                localStorage.clear();
                sessionStorage.clear();
                
                // 显示提示消息
                ElementUI.Message.warning('登录已过期，请重新登录');
                
                // 跳转到登录页
                router.push('/login').then(() => {
                    // 延迟重置标志，避免并发请求多次触发
                    setTimeout(() => {
                        isRedirecting = false;
                    }, 1000);
                }).catch(() => {
                    isRedirecting = false;
                });
            }
        }
        
        return response;
    },
    error => {
        // 处理HTTP错误状态码（如401, 403等）
        if (error.response) {
            const status = error.response.status;
            const currentPath = router.currentRoute.path;
            const isOnLoginPage = currentPath === '/' || currentPath === '/login';
            
            // 401未授权或403禁止访问
            if ((status === 401 || status === 403) && !isOnLoginPage && !isRedirecting) {
                isRedirecting = true;
                
                // 清除本地存储
                localStorage.clear();
                sessionStorage.clear();
                
                // 显示提示消息
                ElementUI.Message.warning('登录已过期，请重新登录');
                
                // 跳转到登录页
                router.push('/login').then(() => {
                    setTimeout(() => {
                        isRedirecting = false;
                    }, 1000);
                }).catch(() => {
                    isRedirecting = false;
                });
            } else if (status === 500) {
                ElementUI.Message.error('服务器错误，请稍后重试');
            } else if (status === 404) {
                ElementUI.Message.error('请求的资源不存在');
            }
        } else if (error.message) {
            // 网络错误或超时
            if (error.message.includes('timeout')) {
                ElementUI.Message.error('请求超时，请检查网络连接');
            } else if (error.message.includes('Network Error')) {
                ElementUI.Message.error('网络错误，请检查网络连接');
            }
        }
        
        return Promise.reject(error);
    }
);

export default request