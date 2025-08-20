"use client";

import { createContext, useContext, ReactNode, useState, useEffect } from 'react';
import useTlsFingerprint from '@/hooks/useTlsFingerprint';

interface FingerprintContextType {
  fingerprint: string;
  isLoading: boolean;
}

const FingerprintContext = createContext<FingerprintContextType>({
  fingerprint: '',
  isLoading: true
});

export function useFingerprintContext() {
  return useContext(FingerprintContext);
}

interface FingerprintProviderProps {
  children: ReactNode;
}

export function TlsFingerprintProvider({ children }: FingerprintProviderProps) {
  const { fingerprint, isLoading } = useTlsFingerprint();
  const [isClient, setIsClient] = useState(false);
  
  useEffect(() => {
    setIsClient(true);
    
    // 保存指紋到 localStorage，以便在 API 請求中使用
    if (fingerprint && !isLoading) {
      localStorage.setItem('tls-fingerprint', fingerprint);
    }
  }, [fingerprint, isLoading]);
  
  if (!isClient) {
    // 首次渲染在服務器端時提供空值
    return <FingerprintContext.Provider value={{ fingerprint: '', isLoading: true }}>
      {children}
    </FingerprintContext.Provider>;
  }
  
  return (
    <FingerprintContext.Provider value={{ fingerprint, isLoading }}>
      {children}
    </FingerprintContext.Provider>
  );
}
