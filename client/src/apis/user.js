import instance from '../utils/request'

// 用户登录
export const login = (email, password) => {
    return instance.post('/user/login', { email, password })
}

// 用户注册
export const register = (username, email, password, re_password) => {
    return instance.post('/user/register', { username, email, password, re_password })
}
