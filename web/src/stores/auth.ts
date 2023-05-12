import { defineStore } from 'pinia';
import { useErrorStore } from '@/stores/error';
import axios from 'axios';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    loggedIn: localStorage.getItem('token') ? true : false,
    userID: localStorage.getItem('userID'),
    username: localStorage.getItem('username'),
  }),
  actions: {
    async login(username: string, password: string) {
      let formData = new FormData();
      formData.append('username', username);
      formData.append('password', password);

      try {
        const response = await axios.post('/auth/login', formData);
        const data = response.data;
        if (data) {
          useErrorStore().clear();
          const token = `Bearer ${data.token}`;
          localStorage.setItem('token', token);
          localStorage.setItem('userID', data.userID);
          localStorage.setItem('username', data.username);
          axios.defaults.headers.common['Authorization'] = token;
          axios.defaults.headers.common['UserID'] = data.userID;
          axios.defaults.headers.common['Username'] = data.username;
          this.userID = data.userID;
          this.username = data.username;
          this.loggedIn = true;
        }
      } catch (error) {
        console.log('login error response:', error.response)
        useErrorStore().set(error.response.data);
      }
    },
    async logout() {
      try {
        const response = await fetch('/auth/logout', { method: 'post' });
        if (response.status != 200) {
          console.warn('logout response status is not 200')
        }
      } catch (error) {
        console.error('logout error response:', error.response)
      }
      localStorage.removeItem('token');
      localStorage.removeItem('userID');
      localStorage.removeItem('username');
      this.$reset();
    },
  },
});

export type AuthStore = ReturnType<typeof useAuthStore>