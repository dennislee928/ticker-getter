"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import EventList from '@/components/events/EventList';
import SearchInput from '@/components/ui/SearchInput';
import { FunnelIcon } from '@heroicons/react/24/outline';

export default function EventsPage() {
  const router = useRouter();
  const [showFilters, setShowFilters] = useState(false);
  
  // 過濾條件 (可以根據需要擴展)
  const [filters, setFilters] = useState({
    priceRange: [0, 5000],
    dateRange: '',
    location: ''
  });
  
  const handleSearch = (query: string) => {
    router.push(`/events/search?q=${encodeURIComponent(query)}`);
  };

  const toggleFilters = () => {
    setShowFilters(!showFilters);
  };

  const handleFilterChange = (key: string, value: any) => {
    setFilters({
      ...filters,
      [key]: value
    });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          <h1 className="text-3xl font-bold text-gray-900">活動列表</h1>

          <div className="sm:flex sm:justify-between sm:items-center space-y-3 sm:space-y-0">
            <div className="relative sm:max-w-md w-full">
              <SearchInput 
                placeholder="搜尋活動名稱、地點、日期..." 
                onSearch={handleSearch}
              />
            </div>

            <button
              type="button"
              onClick={toggleFilters}
              className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <FunnelIcon className="h-5 w-5 mr-2 text-gray-500" />
              篩選條件
            </button>
          </div>

          {showFilters && (
            <div className="bg-white p-4 rounded-md shadow">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label htmlFor="price-range" className="block text-sm font-medium text-gray-700">
                    價格範圍
                  </label>
                  <div className="mt-1 flex gap-2">
                    <input
                      type="number"
                      min="0"
                      max={filters.priceRange[1]}
                      value={filters.priceRange[0]}
                      onChange={(e) => handleFilterChange('priceRange', [parseInt(e.target.value), filters.priceRange[1]])}
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                      placeholder="最低"
                    />
                    <span className="text-gray-500 self-center">至</span>
                    <input
                      type="number"
                      min={filters.priceRange[0]}
                      value={filters.priceRange[1]}
                      onChange={(e) => handleFilterChange('priceRange', [filters.priceRange[0], parseInt(e.target.value)])}
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                      placeholder="最高"
                    />
                  </div>
                </div>

                <div>
                  <label htmlFor="date-range" className="block text-sm font-medium text-gray-700">
                    日期
                  </label>
                  <select
                    id="date-range"
                    value={filters.dateRange}
                    onChange={(e) => handleFilterChange('dateRange', e.target.value)}
                    className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
                  >
                    <option value="">所有日期</option>
                    <option value="today">今天</option>
                    <option value="this-week">本週</option>
                    <option value="this-month">本月</option>
                    <option value="next-month">下個月</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="location" className="block text-sm font-medium text-gray-700">
                    地點
                  </label>
                  <select
                    id="location"
                    value={filters.location}
                    onChange={(e) => handleFilterChange('location', e.target.value)}
                    className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
                  >
                    <option value="">所有地點</option>
                    <option value="taipei">台北</option>
                    <option value="taichung">台中</option>
                    <option value="kaohsiung">高雄</option>
                  </select>
                </div>
              </div>
            </div>
          )}

          <div className="pt-4">
            <EventList />
          </div>
        </div>
      </main>

      <footer className="bg-gray-800 mt-12">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <p className="text-center text-gray-300 text-sm">
            &copy; 2024 票務購買平台. 保留所有權利.
          </p>
        </div>
      </footer>
    </div>
  );
}
