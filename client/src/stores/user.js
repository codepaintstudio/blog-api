import { defineStore } from 'pinia'
import request from '../utils/request'
import router from '../router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}')
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  
  actions: {
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    },
    
    setUserInfo(userInfo) {
      this.userInfo = userInfo
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    
    async login(email, password) {
      try {
        const data = await request.post('/user/login', { email, password })
        this.setToken(data.token)
        this.setUserInfo(data.user)
        router.push('/')
        return data
      } catch (error) {
        throw error
      }
    },
    
    async register(username, email, password) {
      try {
        const data = await request.post('/user/register', { username, email, password })
        this.setToken(data.token)
        this.setUserInfo(data.user)
        router.push('/')
        return data
      } catch (error) {
        throw error
      }
    },
    
    logout() {
      this.token = ''
      this.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      router.push('/login')
    }
  }
})
