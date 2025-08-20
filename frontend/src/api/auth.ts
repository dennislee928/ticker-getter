import apiClient from './client';

interface LoginData {
  email: string;
  password: string;
}

interface RegisterData {
  name: string;
  email: string;
  phone: string;
  password: string;
}

interface AuthResponse {
  token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
}

// 用戶登入
export const login = async (data: LoginData): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>('/auth/login', data);
  return response.data;
};

// 用戶註冊
export const register = async (data: RegisterData): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>('/auth/register', data);
  return response.data;
};

// 登出
export const logout = async (): Promise<void> => {
  await apiClient.post('/auth/logout');
  // 清除本地存儲
  localStorage.removeItem('auth_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('user');
};

// 獲取當前用戶信息
export const getCurrentUser = async () => {
  const response = await apiClient.get('/auth/me');
  return response.data;
};

// 重設密碼請求
export const requestPasswordReset = async (email: string) => {
  const response = await apiClient.post('/auth/forgot-password', { email });
  return response.data;
};

// 重設密碼
export const resetPassword = async (token: string, password: string) => {
  const response = await apiClient.post('/auth/reset-password', { 
    token, 
    password 
  });
  return response.data;
};
