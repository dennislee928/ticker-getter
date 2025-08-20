"use client";

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import Navbar from '@/components/Navbar';
import EventList from '@/components/events/EventList';
import SearchInput from '@/components/ui/SearchInput';
import { searchEvents, Event } from '@/api/events';

export default function EventSearchPage() {
  const searchParams = useSearchParams();
  const query = searchParams.get('q') || '';
  
  const [searchResults, setSearchResults] = useState<Event[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    const performSearch = async () => {
      if (!query) {
        setSearchResults([]);
        setIsLoading(false);
        return;
      }
      
      try {
        setIsLoading(true);
        const results = await searchEvents(query);
        setSearchResults(results);
      } catch (err) {
        console.error('搜尋錯誤:', err);
        setError('搜尋過程中發生錯誤，請稍後再試。');
      } finally {
        setIsLoading(false);
      }
    };
    
    performSearch();
  }, [query]);
  
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">搜尋結果</h1>
            {query && (
              <p className="mt-1 text-sm text-gray-600">
                關鍵字: <span className="font-medium">"{query}"</span>
              </p>
            )}
          </div>
          
          <div className="max-w-md">
            <SearchInput defaultValue={query} />
          </div>
          
          {isLoading ? (
            <div className="text-center py-12">
              <p className="text-lg text-gray-600">載入中...</p>
            </div>
          ) : error ? (
            <div className="rounded-md bg-red-50 p-4">
              <div className="flex">
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-red-800">錯誤</h3>
                  <div className="mt-2 text-sm text-red-700">
                    <p>{error}</p>
                  </div>
                </div>
              </div>
            </div>
          ) : searchResults.length > 0 ? (
            <div>
              <p className="mb-4 text-gray-600">
                找到 {searchResults.length} 個符合的結果
              </p>
              <EventList initialEvents={searchResults} />
            </div>
          ) : (
            <div className="text-center py-12 bg-white rounded-lg shadow">
              <p className="text-lg text-gray-600">
                找不到符合 "{query}" 的活動。
              </p>
              <p className="mt-2 text-sm text-gray-500">
                請嘗試使用不同的關鍵字或瀏覽所有活動。
              </p>
            </div>
          )}
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
