import apiClient from './client';
import { Event, TicketType } from './events';

// 用戶管理 API
export interface User {
  id: string;
  email: string;
  name: string;
  phone: string;
  role: string;
  created_at: string;
}

export const getUsers = async (page: number = 1, limit: number = 10) => {
  const response = await apiClient.get<{ users: User[], total: number }>(`/admin/users?page=${page}&limit=${limit}`);
  return response.data;
};

export const getUserById = async (id: string) => {
  const response = await apiClient.get<User>(`/admin/users/${id}`);
  return response.data;
};

export const updateUserRole = async (id: string, role: string) => {
  const response = await apiClient.patch(`/admin/users/${id}/role`, { role });
  return response.data;
};

// 活動管理 API
export const getAdminEvents = async (page: number = 1, limit: number = 10) => {
  const response = await apiClient.get<{ events: Event[], total: number }>(`/admin/events?page=${page}&limit=${limit}`);
  return response.data;
};

export const createEvent = async (eventData: Omit<Event, 'id' | 'created_at' | 'updated_at'>) => {
  const response = await apiClient.post<Event>('/admin/events', eventData);
  return response.data;
};

export const updateEvent = async (id: string, eventData: Partial<Event>) => {
  const response = await apiClient.put<Event>(`/admin/events/${id}`, eventData);
  return response.data;
};

export const deleteEvent = async (id: string) => {
  const response = await apiClient.delete(`/admin/events/${id}`);
  return response.data;
};

// 票券管理 API
export const getTicketTypes = async (eventId: string) => {
  const response = await apiClient.get<TicketType[]>(`/admin/events/${eventId}/ticket-types`);
  return response.data;
};

export const createTicketType = async (eventId: string, ticketTypeData: Omit<TicketType, 'id' | 'event_id' | 'created_at' | 'updated_at'>) => {
  const response = await apiClient.post<TicketType>(`/admin/events/${eventId}/ticket-types`, ticketTypeData);
  return response.data;
};

export const updateTicketType = async (id: string, ticketTypeData: Partial<TicketType>) => {
  const response = await apiClient.put<TicketType>(`/admin/ticket-types/${id}`, ticketTypeData);
  return response.data;
};

export const deleteTicketType = async (id: string) => {
  const response = await apiClient.delete(`/admin/ticket-types/${id}`);
  return response.data;
};

// 訂單管理 API
export interface Order {
  id: string;
  user_id: string;
  user_name: string;
  total_amount: number;
  status: string;
  payment_method: string;
  payment_status: string;
  created_at: string;
}

export const getOrders = async (page: number = 1, limit: number = 10, status?: string) => {
  let url = `/admin/orders?page=${page}&limit=${limit}`;
  if (status) {
    url += `&status=${status}`;
  }
  const response = await apiClient.get<{ orders: Order[], total: number }>(url);
  return response.data;
};

export const getOrderById = async (id: string) => {
  const response = await apiClient.get<Order>(`/admin/orders/${id}`);
  return response.data;
};

export const updateOrderStatus = async (id: string, status: string) => {
  const response = await apiClient.patch(`/admin/orders/${id}/status`, { status });
  return response.data;
};

// 統計數據 API
export interface DashboardStats {
  total_sales: number;
  total_users: number;
  total_tickets: number;
  total_events: number;
  recent_activities: {
    id: string;
    type: string;
    status: string;
    message: string;
    timestamp: string;
    amount?: number;
    user?: string;
  }[];
  sales_trend: {
    date: string;
    amount: number;
  }[];
}

export const getDashboardStats = async () => {
  const response = await apiClient.get<DashboardStats>('/admin/stats/dashboard');
  return response.data;
};
