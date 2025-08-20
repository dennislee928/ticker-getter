"use client";

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useRouter } from 'next/navigation';
import { logout } from '@/api/auth';
import { toast } from 'react-hot-toast';
import {
  HomeIcon,
  TicketIcon,
  UsersIcon,
  ChartBarIcon,
  CalendarIcon,
  CurrencyDollarIcon,
  Cog6ToothIcon,
  ArrowLeftOnRectangleIcon
} from '@heroicons/react/24/outline';

interface AdminLayoutProps {
  children: React.ReactNode;
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const pathname = usePathname();
  const router = useRouter();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isClient, setIsClient] = useState(false);
  const [userName, setUserName] = useState('管理員');

  // 客戶端檢測
  useEffect(() => {
    setIsClient(true);
    
    // 從本地存儲獲取用戶信息
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        if (user.name) {
          setUserName(user.name);
        }
      } catch (e) {
        console.error('解析用戶信息失敗:', e);
      }
    }
  }, []);
  
  // 檢查是否為管理員
  useEffect(() => {
    if (isClient) {
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        router.push('/auth/login');
        return;
      }
      
      try {
        const user = JSON.parse(userStr);
        if (user.role !== 'admin') {
          toast.error('您沒有管理員權限');
          router.push('/');
        }
      } catch (e) {
        router.push('/auth/login');
      }
    }
  }, [isClient, router]);

  const navigation = [
    { name: '儀表板', href: '/admin', icon: HomeIcon, current: pathname === '/admin' },
    { name: '活動管理', href: '/admin/events', icon: CalendarIcon, current: pathname.startsWith('/admin/events') },
    { name: '票券管理', href: '/admin/tickets', icon: TicketIcon, current: pathname.startsWith('/admin/tickets') },
    { name: '用戶管理', href: '/admin/users', icon: UsersIcon, current: pathname.startsWith('/admin/users') },
    { name: '訂單管理', href: '/admin/orders', icon: CurrencyDollarIcon, current: pathname.startsWith('/admin/orders') },
    { name: '報表統計', href: '/admin/reports', icon: ChartBarIcon, current: pathname.startsWith('/admin/reports') },
    { name: '系統設置', href: '/admin/settings', icon: Cog6ToothIcon, current: pathname.startsWith('/admin/settings') },
  ];

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  const handleLogout = async () => {
    try {
      await logout();
      toast.success('已登出系統');
      router.push('/auth/login');
    } catch (error) {
      console.error('登出失敗:', error);
      toast.error('登出失敗，請稍後再試');
    }
  };

  if (!isClient) {
    // 渲染伺服器端骨架
    return (
      <div className="min-h-screen bg-gray-100">
        <div className="h-full">
          <div className="flex h-full">
            <div className="hidden md:flex md:w-64 md:flex-col">
              <div className="flex flex-col flex-grow bg-primary-700 pt-5 pb-4 overflow-y-auto">
                <div className="flex items-center flex-shrink-0 px-4">
                  <span className="text-xl font-bold text-white">票務管理後台</span>
                </div>
                <div className="mt-5 flex-1 flex flex-col">
                  <nav className="flex-1 px-2 space-y-1">
                    {/* 骨架加載 */}
                    {Array(7).fill(0).map((_, i) => (
                      <div key={i} className="flex items-center px-2 py-2 rounded-md">
                        <div className="h-5 w-5 bg-gray-300 rounded-md mr-3"></div>
                        <div className="h-4 bg-gray-300 rounded w-24"></div>
                      </div>
                    ))}
                  </nav>
                </div>
              </div>
            </div>
            <div className="flex flex-col w-0 flex-1 overflow-hidden">
              <div className="relative z-10 flex-shrink-0 flex h-16 bg-white shadow">
                <div className="flex-1 px-4 flex justify-between">
                  <div className="flex-1 flex items-center">
                    <div className="h-8 w-48 bg-gray-300 rounded"></div>
                  </div>
                </div>
              </div>
              <main className="flex-1 relative overflow-y-auto focus:outline-none">
                <div className="py-6">
                  <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
                    {/* 骨架加載 */}
                    <div className="h-96 rounded-lg border-4 border-dashed border-gray-200"></div>
                  </div>
                </div>
              </main>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="h-full">
        <div className="flex h-full">
          {/* 側邊欄 - 手機版 */}
          {isSidebarOpen && (
            <div className="md:hidden fixed inset-0 flex z-40">
              <div className="fixed inset-0 bg-gray-600 bg-opacity-75" onClick={toggleSidebar}></div>
              <div className="relative flex-1 flex flex-col max-w-xs w-full bg-primary-700">
                <div className="absolute top-0 right-0 -mr-12 pt-2">
                  <button
                    type="button"
                    className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                    onClick={toggleSidebar}
                  >
                    <span className="sr-only">關閉側邊欄</span>
                    <svg className="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
                
                <div className="flex-1 h-0 pt-5 pb-4 overflow-y-auto">
                  <div className="flex-shrink-0 flex items-center px-4">
                    <span className="text-xl font-bold text-white">票務管理後台</span>
                  </div>
                  <nav className="mt-5 px-2 space-y-1">
                    {navigation.map((item) => (
                      <Link
                        key={item.name}
                        href={item.href}
                        className={`group flex items-center px-2 py-2 text-base font-medium rounded-md ${
                          item.current
                            ? 'bg-primary-800 text-white'
                            : 'text-white hover:bg-primary-600'
                        }`}
                      >
                        <item.icon className="mr-4 h-6 w-6 text-primary-300" aria-hidden="true" />
                        {item.name}
                      </Link>
                    ))}
                    
                    <button
                      onClick={handleLogout}
                      className="w-full text-left group flex items-center px-2 py-2 text-base font-medium rounded-md text-white hover:bg-primary-600"
                    >
                      <ArrowLeftOnRectangleIcon className="mr-4 h-6 w-6 text-primary-300" aria-hidden="true" />
                      登出系統
                    </button>
                  </nav>
                </div>
              </div>
            </div>
          )}
          
          {/* 側邊欄 - 桌面版 */}
          <div className="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0">
            <div className="flex-1 flex flex-col min-h-0 bg-primary-700">
              <div className="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
                <div className="flex items-center flex-shrink-0 px-4">
                  <span className="text-xl font-bold text-white">票務管理後台</span>
                </div>
                <nav className="mt-5 flex-1 px-2 space-y-1">
                  {navigation.map((item) => (
                    <Link
                      key={item.name}
                      href={item.href}
                      className={`group flex items-center px-2 py-2 text-sm font-medium rounded-md ${
                        item.current
                          ? 'bg-primary-800 text-white'
                          : 'text-white hover:bg-primary-600'
                      }`}
                    >
                      <item.icon className="mr-3 h-5 w-5 text-primary-300" aria-hidden="true" />
                      {item.name}
                    </Link>
                  ))}
                </nav>
              </div>
              <div className="flex-shrink-0 flex border-t border-primary-800 p-4">
                <button
                  onClick={handleLogout}
                  className="w-full text-left flex items-center px-2 py-2 text-sm font-medium rounded-md text-white hover:bg-primary-600"
                >
                  <ArrowLeftOnRectangleIcon className="mr-3 h-5 w-5 text-primary-300" aria-hidden="true" />
                  登出系統
                </button>
              </div>
            </div>
          </div>
          
          <div className="md:pl-64 flex flex-col flex-1">
            <div className="sticky top-0 z-10 md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3 bg-gray-100">
              <button
                type="button"
                className="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
                onClick={toggleSidebar}
              >
                <span className="sr-only">開啟側邊欄</span>
                <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
            
            <main className="flex-1">
              <div className="py-6">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
                  <div className="py-4">
                    {children}
                  </div>
                </div>
              </div>
            </main>
          </div>
        </div>
      </div>
    </div>
  );
}
