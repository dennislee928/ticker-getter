"use client";

import { useState, useEffect } from 'react';
import StatsCard from '@/components/admin/StatsCard';
import RecentActivities from '@/components/admin/RecentActivities';
import { 
  TicketIcon, 
  UsersIcon, 
  CurrencyDollarIcon, 
  CalendarIcon 
} from '@heroicons/react/24/outline';
import { formatPrice } from '@/lib/utils';

export default function AdminDashboard() {
  const [isLoading, setIsLoading] = useState(true);

  // 模擬從 API 獲取數據
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 1000);

    return () => clearTimeout(timer);
  }, []);

  // 模擬統計數據
  const stats = {
    totalSales: 1235000,
    totalUsers: 2456,
    totalTickets: 8765,
    totalEvents: 45,
  };

  // 模擬最近活動
  const recentActivities = [
    {
      id: '1',
      type: 'order',
      status: 'success',
      message: '新訂單 #12345',
      timestamp: new Date().toISOString(),
      amount: 2500,
      user: '張三',
    },
    {
      id: '2',
      type: 'ticket',
      status: 'pending',
      message: '票券兌換請求',
      timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(), // 30 分鐘前
      user: '李四',
    },
    {
      id: '3',
      type: 'refund',
      status: 'failed',
      message: '退款請求失敗',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(), // 2 小時前
      amount: 1200,
      user: '王五',
    },
    {
      id: '4',
      type: 'user',
      status: 'success',
      message: '新用戶註冊',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 3).toISOString(), // 3 小時前
      user: '趙六',
    },
  ] as const;

  if (isLoading) {
    return (
      <div className="animate-pulse">
        <h1 className="text-2xl font-semibold text-gray-900 mb-6">儀表板</h1>
        
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-6">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="bg-white shadow rounded-lg h-28">
              <div className="p-5">
                <div className="flex items-center">
                  <div className="flex-shrink-0 rounded-md bg-gray-200 p-3 h-12 w-12"></div>
                  <div className="ml-5 w-0 flex-1">
                    <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                    <div className="h-6 bg-gray-300 rounded w-1/2"></div>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
        
        <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
          <div className="bg-white shadow rounded-lg h-80">
            <div className="p-5">
              <div className="h-6 bg-gray-200 rounded w-1/4 mb-6"></div>
              <div className="space-y-4">
                {[1, 2, 3, 4].map((i) => (
                  <div key={i} className="h-12 bg-gray-100 rounded"></div>
                ))}
              </div>
            </div>
          </div>
          
          <div className="bg-white shadow rounded-lg h-80">
            <div className="p-5">
              <div className="h-6 bg-gray-200 rounded w-1/4 mb-6"></div>
              <div className="h-64 bg-gray-100 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-semibold text-gray-900 mb-6">儀表板</h1>
      
      {/* 統計卡片 */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-6">
        <StatsCard 
          title="總銷售額" 
          value={formatPrice(stats.totalSales)} 
          icon={<CurrencyDollarIcon className="h-6 w-6 text-primary-600" />} 
          change={{ value: 12.5, positive: true }}
        />
        
        <StatsCard 
          title="用戶數量" 
          value={stats.totalUsers.toLocaleString()} 
          icon={<UsersIcon className="h-6 w-6 text-blue-600" />}
          textColor="text-blue-600"
          bgColor="bg-blue-50"
          change={{ value: 8.3, positive: true }}
        />
        
        <StatsCard 
          title="票券銷量" 
          value={stats.totalTickets.toLocaleString()} 
          icon={<TicketIcon className="h-6 w-6 text-green-600" />}
          textColor="text-green-600"
          bgColor="bg-green-50"
          change={{ value: 5.2, positive: true }}
        />
        
        <StatsCard 
          title="活動數量" 
          value={stats.totalEvents} 
          icon={<CalendarIcon className="h-6 w-6 text-purple-600" />}
          textColor="text-purple-600"
          bgColor="bg-purple-50"
          change={{ value: 2.1, positive: true }}
        />
      </div>
      
      {/* 最近活動和圖表 */}
      <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
        <RecentActivities activities={recentActivities} />
        
        <div className="bg-white shadow rounded-lg overflow-hidden">
          <div className="px-4 py-5 sm:px-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              銷售趨勢
            </h3>
          </div>
          <div className="px-4 py-5 sm:p-6 h-64 flex items-center justify-center">
            <p className="text-gray-500">
              圖表將在此顯示 (需要圖表庫)
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
