import apiClient from './client';

export interface Event {
  id: string;
  title: string;
  description: string;
  location: string;
  start_time: string;
  end_time: string;
  created_at: string;
  updated_at: string;
}

export interface TicketType {
  id: string;
  event_id: string;
  name: string;
  price: number;
  total_quantity: number;
  available_quantity: number;
  sale_start: string;
  sale_end: string;
}

// 獲取所有活動
export const getAllEvents = async () => {
  const response = await apiClient.get<Event[]>('/events');
  return response.data;
};

// 獲取熱門活動
export const getFeaturedEvents = async (limit: number = 6) => {
  const response = await apiClient.get<Event[]>(`/events/featured?limit=${limit}`);
  return response.data;
};

// 獲取單個活動詳情
export const getEvent = async (eventId: string) => {
  const response = await apiClient.get<Event>(`/events/${eventId}`);
  return response.data;
};

// 獲取活動的票種
export const getEventTicketTypes = async (eventId: string) => {
  const response = await apiClient.get<TicketType[]>(`/events/${eventId}/ticket-types`);
  return response.data;
};

// 搜索活動
export const searchEvents = async (query: string) => {
  const response = await apiClient.get<Event[]>(`/events/search?q=${encodeURIComponent(query)}`);
  return response.data;
};
