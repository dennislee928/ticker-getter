"use client";

import { useState, useEffect } from 'react';
import { formatPrice } from '@/lib/utils';
import { TicketIcon, ShoppingCartIcon } from '@heroicons/react/24/outline';
import { TicketType } from '@/api/events';
import { useFingerprintContext } from '@/components/TlsFingerprintProvider';
import { checkFingerprintStatus, purchaseTicket } from '@/api/tickets';
import { toast } from 'react-hot-toast';

interface TicketTypeCardProps {
  ticketType: TicketType;
  eventId: string;
  onPurchase?: () => void;
}

export default function TicketTypeCard({ ticketType, eventId, onPurchase }: TicketTypeCardProps) {
  const [quantity, setQuantity] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [alreadyPurchased, setAlreadyPurchased] = useState(false);
  const { fingerprint } = useFingerprintContext();
  
  const now = new Date();
  const saleStart = new Date(ticketType.sale_start);
  const saleEnd = new Date(ticketType.sale_end);
  
  const isSaleActive = now >= saleStart && now <= saleEnd;
  const isSoldOut = ticketType.available_quantity <= 0;

  // 檢查是否已購買
  const checkIfAlreadyPurchased = async () => {
    if (!fingerprint) return;
    
    try {
      const result = await checkFingerprintStatus(ticketType.id);
      setAlreadyPurchased(result.already_purchased);
    } catch (err) {
      console.error('檢查購買狀態失敗:', err);
    }
  };

  // 當指紋加載完成後檢查
  useEffect(() => {
    if (fingerprint) {
      checkIfAlreadyPurchased();
    }
  }, [fingerprint]);

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    if (value >= 1 && value <= Math.min(10, ticketType.available_quantity)) {
      setQuantity(value);
    }
  };

  const handlePurchase = async () => {
    try {
      setIsLoading(true);
      
      // 檢查是否已購買
      if (!fingerprint) {
        toast.error('無法識別您的設備，請嘗試重新整理頁面');
        return;
      }
      
      const statusCheck = await checkFingerprintStatus(ticketType.id);
      if (statusCheck.already_purchased) {
        toast.error('您已購買過此票券，無法重複購買');
        setAlreadyPurchased(true);
        return;
      }
      
      // 嘗試購買
      const result = await purchaseTicket(ticketType.id, quantity);
      
      // 購買成功
      toast.success('購票成功！');
      
      // 通知父組件
      if (onPurchase) {
        onPurchase();
      }
      
      // 標記為已購買
      setAlreadyPurchased(true);
    } catch (err: any) {
      console.error('購票失敗:', err);
      toast.error(err.response?.data?.error || '購票失敗，請稍後再試');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="border border-gray-200 rounded-lg p-4 bg-white shadow-sm hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start">
        <div>
          <div className="flex items-center">
            <TicketIcon className="h-5 w-5 text-primary-500 mr-2" />
            <h3 className="text-lg font-medium text-gray-900">{ticketType.name}</h3>
          </div>
          <p className="mt-1 text-lg font-semibold text-primary-600">
            {formatPrice(ticketType.price)}
          </p>
        </div>
        
        <div className="text-right">
          <p className="text-sm text-gray-500">
            剩餘 <span className="font-medium">{ticketType.available_quantity}</span> 張
          </p>
          <p className="text-sm text-gray-500 mt-1">
            限購 <span className="font-medium">10</span> 張
          </p>
        </div>
      </div>
      
      {!isSaleActive && (
        <div className="mt-4">
          {now < saleStart ? (
            <p className="text-sm text-amber-600">
              尚未開始販售，{new Date(ticketType.sale_start).toLocaleDateString()} 開始
            </p>
          ) : (
            <p className="text-sm text-gray-500">
              已結束販售
            </p>
          )}
        </div>
      )}
      
      {isSaleActive && !isSoldOut && !alreadyPurchased && (
        <div className="mt-4 flex items-center">
          <div className="w-20">
            <label htmlFor={`quantity-${ticketType.id}`} className="sr-only">
              數量
            </label>
            <input
              id={`quantity-${ticketType.id}`}
              type="number"
              min="1"
              max={Math.min(10, ticketType.available_quantity)}
              value={quantity}
              onChange={handleQuantityChange}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            />
          </div>
          <button
            type="button"
            disabled={isLoading}
            onClick={handlePurchase}
            className="ml-4 flex-1 bg-primary-600 border border-transparent rounded-md py-2 px-4 flex items-center justify-center text-sm font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? (
              '處理中...'
            ) : (
              <>
                <ShoppingCartIcon className="h-5 w-5 mr-2" />
                購票
              </>
            )}
          </button>
        </div>
      )}
      
      {isSaleActive && isSoldOut && (
        <div className="mt-4">
          <p className="text-sm text-red-600 font-medium">
            已售完
          </p>
        </div>
      )}
      
      {alreadyPurchased && (
        <div className="mt-4">
          <p className="text-sm text-green-600 font-medium">
            您已購買此票券
          </p>
        </div>
      )}
    </div>
  );
}