import axios from 'axios'
import router from '../router'
import { useUserStore } from '../stores/user'

const baseURL = 'http://localhost:7777/api'

const instance = axios.create({
    baseURL,
    timeout: 10000
})

// 请求拦截器
instance.interceptors.request.use(
    (config) => {
        // 携带token
        const userStore = useUserStore()
        if (userStore.token) {
            config.headers.Authorization = `Bearer ${userStore.token}`
        }
        return config
    },
    (err) => Promise.reject(err)
)

// 响应拦截器
instance.interceptors.response.use(
    (res) => {
        // 只返回响应数据
        return res.data.data
    },
    (err) => {
        // 处理401错误：权限不足或者token过期
        if (err.response?.status === 401) {
            const userStore = useUserStore()
            userStore.logout()
        }
        return Promise.reject(err.response?.data || err)
    }
)

export default instance
