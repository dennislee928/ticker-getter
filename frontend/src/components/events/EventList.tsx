"use client";

import { useState, useEffect } from 'react';
import EventCard from './EventCard';
import { getAllEvents, Event } from '@/api/events';

interface EventListProps {
  initialEvents?: Event[];
}

export default function EventList({ initialEvents }: EventListProps) {
  const [events, setEvents] = useState<Event[]>(initialEvents || []);
  const [isLoading, setIsLoading] = useState<boolean>(!initialEvents);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!initialEvents) {
      const fetchEvents = async () => {
        try {
          setIsLoading(true);
          const eventsData = await getAllEvents();
          setEvents(eventsData);
        } catch (err) {
          console.error('Failed to fetch events:', err);
          setError('無法加載活動列表，請稍後再試。');
        } finally {
          setIsLoading(false);
        }
      };

      fetchEvents();
    }
  }, [initialEvents]);

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {[...Array(6)].map((_, index) => (
          <div key={index} className="animate-pulse rounded-lg overflow-hidden bg-gray-100 h-72">
            <div className="h-48 bg-gray-200"></div>
            <div className="p-4 space-y-2">
              <div className="h-4 bg-gray-200 rounded w-3/4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              <div className="h-4 bg-gray-200 rounded w-2/3"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (error) {
    return (
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
    );
  }

  if (events.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-lg text-gray-600">
          目前沒有可用的活動。
        </p>
      </div>
    );
  }

  // 設定每個活動的示例最低價格 (實際應用中應從 API 獲取)
  const eventPrices: Record<string, number> = events.reduce((acc, event) => {
    // 為每個活動設置一個隨機的最低價格，範圍在 300 到 3000 之間
    acc[event.id] = Math.floor(Math.random() * 2700) + 300;
    return acc;
  }, {} as Record<string, number>);

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      {events.map(event => (
        <EventCard 
          key={event.id} 
          event={event} 
          minPrice={eventPrices[event.id]}
        />
      ))}
    </div>
  );
}
