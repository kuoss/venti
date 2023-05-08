import type { App } from 'vue';
import type { AxiosInstance } from 'axios';
import type { HeadersDefaults } from 'axios';
import { useErrorStore } from '../stores/error';
import router from '../router';
import axios from 'axios';

interface CommonHeaderProperties extends HeadersDefaults {
  Authorization: string;
  UserID: string;
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance,
  }
}

export default {
  install: (app: App) => {
    const instance = axios.create();

    instance.defaults.baseURL = import.meta.env.AXIOS_BASE_URL;
    instance.defaults.withCredentials = true;
    instance.defaults.headers = {
      Authorization: localStorage.getItem('token'),
      UserID: localStorage.getItem('userID'),
    } as CommonHeaderProperties;

    instance.interceptors.request.use(
      function (config) {
        useErrorStore().$reset();
        return config;
      },
      function (error) {
        return Promise.reject(error);
      }
    );

    instance.interceptors.response.use(
      function (response) {
        return response;
      },
      function (error) {
        console.log('error=', error);
        switch (error.response.status) {
          case 401:
            console.log(401);
            localStorage.removeItem('token');
            localStorage.removeItem('userID');
            window.location.reload();
            break;
          case 403:
          case 404:
            router.push({
              name: 'error',
              params: {
                error: {
                  message: error.response.data.message,
                  status: error.status,
                }.toString()
              },
            });
            break;
          case 422:
            useErrorStore().$state = error.response.data;
            break;
          default:
            console.log('axiosPlugin report', error.response.data);
        }
        return Promise.reject(error);
      }
    );
    app.config.globalProperties.$axios = instance;
  }
}