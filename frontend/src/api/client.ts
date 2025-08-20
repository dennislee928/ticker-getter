import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// 創建 Axios 實例
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 請求攔截器，添加授權標頭
apiClient.interceptors.request.use(
  (config) => {
    // 從本地存儲獲取令牌
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 響應攔截器
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    
    // 如果響應狀態為 401（未授權）並且還沒有重試過
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // 嘗試刷新令牌
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const { data } = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });
          
          // 存儲新的令牌
          localStorage.setItem('auth_token', data.token);
          localStorage.setItem('refresh_token', data.refresh_token);
          
          // 更新授權標頭
          apiClient.defaults.headers.common['Authorization'] = `Bearer ${data.token}`;
          originalRequest.headers['Authorization'] = `Bearer ${data.token}`;
          
          // 重試原始請求
          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        // 刷新令牌失敗，清除本地存儲
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        
        // 重定向到登入頁面
        window.location.href = '/auth/login';
      }
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;
