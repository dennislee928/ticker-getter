"use client";

import { useEffect, useState } from 'react';
import { getCookies, setCookie } from 'cookies-next';

// 指紋特徵
interface FingerprintFeatures {
  userAgent: string;
  language: string;
  platform: string;
  screenResolution: string;
  timezone: string;
  cookiesEnabled: boolean;
  canvas: string;
  webgl: string;
  fonts: string[];
  plugins: string[];
  audioContext: string;
}

/**
 * 生成瀏覽器 TLS 指紋，用於識別使用者裝置
 * 此指紋將用於防止重複購買票券
 */
export function useTlsFingerprint() {
  const [fingerprint, setFingerprint] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(true);
  
  // 從 Canvas 獲取特徵
  const getCanvasFingerprint = (): string => {
    try {
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      if (!ctx) return '';
      
      // 繪製一些文字和形狀
      canvas.width = 200;
      canvas.height = 50;
      ctx.textBaseline = 'top';
      ctx.font = '14px Arial';
      ctx.fillStyle = '#f60';
      ctx.fillRect(10, 10, 100, 30);
      ctx.fillStyle = '#069';
      ctx.fillText('Fingerprint', 15, 15);
      
      // 獲取 Canvas 數據
      return canvas.toDataURL().substr(-50);
    } catch (e) {
      return 'canvas-unsupported';
    }
  };
  
  // 從 WebGL 獲取特徵
  const getWebGLFingerprint = (): string => {
    try {
      const canvas = document.createElement('canvas');
      const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
      if (!gl) return '';
      
      const info = gl.getExtension('WEBGL_debug_renderer_info');
      if (!info) return '';
      
      const vendor = gl.getParameter(info.UNMASKED_VENDOR_WEBGL);
      const renderer = gl.getParameter(info.UNMASKED_RENDERER_WEBGL);
      
      return `${vendor}-${renderer}`;
    } catch (e) {
      return 'webgl-unsupported';
    }
  };
  
  // 從 AudioContext 獲取特徵
  const getAudioFingerprint = (): string => {
    try {
      const audioContext = window.AudioContext || (window as any).webkitAudioContext;
      if (!audioContext) return '';
      
      const context = new audioContext();
      const oscillator = context.createOscillator();
      const analyser = context.createAnalyser();
      const gain = context.createGain();
      
      oscillator.connect(analyser);
      analyser.connect(gain);
      
      // 獲取音頻參數
      const hashValue = analyser.fftSize + 
                      analyser.frequencyBinCount + 
                      analyser.minDecibels + 
                      analyser.maxDecibels;
                      
      return hashValue.toString();
    } catch (e) {
      return 'audio-unsupported';
    }
  };
  
  // 生成完整特徵
  const generateFeatures = async (): Promise<FingerprintFeatures> => {
    // 基本特徵
    const features: FingerprintFeatures = {
      userAgent: navigator.userAgent,
      language: navigator.language,
      platform: navigator.platform,
      screenResolution: `${window.screen.width}x${window.screen.height}x${window.screen.colorDepth}`,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      cookiesEnabled: navigator.cookieEnabled,
      canvas: getCanvasFingerprint(),
      webgl: getWebGLFingerprint(),
      audioContext: getAudioFingerprint(),
      fonts: [],
      plugins: []
    };
    
    // 嘗試獲取已安裝字體
    if ((document as any).fonts && (document as any).fonts.ready) {
      try {
        await (document as any).fonts.ready;
        const fontFamilies = ['Arial', 'Courier New', 'Georgia', 'Times New Roman', 'Verdana'];
        const availableFonts = [];
        
        for (const font of fontFamilies) {
          if ((document as any).fonts.check(`12px "${font}"`)) {
            availableFonts.push(font);
          }
        }
        
        features.fonts = availableFonts;
      } catch (e) {
        console.error('Error getting fonts:', e);
      }
    }
    
    // 獲取瀏覽器外掛
    if (navigator.plugins) {
      features.plugins = Array.from(navigator.plugins).map(p => p.name);
    }
    
    return features;
  };
  
  // 哈希特徵為指紋
  const hashFeatures = async (features: FingerprintFeatures): Promise<string> => {
    try {
      const featureStr = JSON.stringify(features);
      const encoder = new TextEncoder();
      const data = encoder.encode(featureStr);
      
      // 使用 SHA-256 哈希
      const hashBuffer = await crypto.subtle.digest('SHA-256', data);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
      
      return hashHex;
    } catch (e) {
      console.error('Error hashing features:', e);
      
      // 降級到簡單串接
      return Object.values(features)
        .flat()
        .join('|')
        .replace(/[^a-zA-Z0-9]/g, '');
    }
  };
  
  // 生成或取回儲存的指紋
  useEffect(() => {
    const initFingerprint = async () => {
      try {
        setIsLoading(true);
        
        // 首先檢查 Cookie 中是否已有指紋
        const cookies = getCookies();
        const storedFingerprint = cookies['tls-fingerprint'];
        
        if (storedFingerprint) {
          setFingerprint(storedFingerprint);
        } else {
          // 生成新指紋
          const features = await generateFeatures();
          const newFingerprint = await hashFeatures(features);
          
          // 儲存到 Cookie，有效期 1 年
          setCookie('tls-fingerprint', newFingerprint, { maxAge: 60 * 60 * 24 * 365 });
          setFingerprint(newFingerprint);
        }
      } catch (error) {
        console.error('Error generating fingerprint:', error);
      } finally {
        setIsLoading(false);
      }
    };
    
    initFingerprint();
  }, []);
  
  return { fingerprint, isLoading };
}

export default useTlsFingerprint;
