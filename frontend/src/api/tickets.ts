import apiClient from './client';

// 檢查票券可用性
export const checkTicketAvailability = async (ticketTypeId: string, quantity: number = 1) => {
  const response = await apiClient.get(`/tickets/check-availability/${ticketTypeId}?quantity=${quantity}`);
  return response.data;
};

// 檢查指紋購買狀態
export const checkFingerprintStatus = async (ticketTypeId: string) => {
  const response = await apiClient.get(`/tickets/check-fingerprint/${ticketTypeId}`);
  return response.data;
};

// 驗證票券
export const validateTicket = async (ticketCode: string) => {
  const response = await apiClient.get(`/tickets/validate/${ticketCode}`);
  return response.data;
};

// 使用票券
export const useTicket = async (ticketCode: string) => {
  const response = await apiClient.post(`/tickets/use/${ticketCode}`);
  return response.data;
};

// 購買票券
export const purchaseTicket = async (ticketTypeId: string, quantity: number = 1) => {
  // 設置 TLS 指紋頭部
  const fingerprint = localStorage.getItem('tls-fingerprint');
  
  // 包含指紋信息
  const config = fingerprint ? {
    headers: {
      'X-TLS-Fingerprint': fingerprint
    }
  } : {};
  
  const response = await apiClient.post('/orders', {
    items: [
      {
        ticket_type_id: ticketTypeId,
        quantity: quantity
      }
    ]
  }, config);
  
  return response.data;
};
