"use client";

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Image from 'next/image';
import Link from 'next/link';
import Navbar from '@/components/Navbar';
import { getEvent, getEventTicketTypes, Event, TicketType } from '@/api/events';
import { formatDate, formatTime } from '@/lib/utils';
import { CalendarIcon, MapPinIcon, ClockIcon, ArrowLeftIcon } from '@heroicons/react/24/outline';
import TicketTypeCard from '@/components/events/TicketTypeCard';
import { toast } from 'react-hot-toast';

export default function EventDetailPage() {
  const { id } = useParams();
  const eventId = Array.isArray(id) ? id[0] : id;
  
  const [event, setEvent] = useState<Event | null>(null);
  const [ticketTypes, setTicketTypes] = useState<TicketType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    const fetchEventDetails = async () => {
      try {
        setIsLoading(true);
        
        // 同時獲取事件和票種信息
        const [eventData, ticketTypesData] = await Promise.all([
          getEvent(eventId),
          getEventTicketTypes(eventId)
        ]);
        
        setEvent(eventData);
        setTicketTypes(ticketTypesData);
      } catch (err) {
        console.error('獲取活動詳情失敗:', err);
        setError('無法載入活動詳情，請稍後再試。');
      } finally {
        setIsLoading(false);
      }
    };
    
    if (eventId) {
      fetchEventDetails();
    }
  }, [eventId]);
  
  const refreshTicketTypes = async () => {
    try {
      const ticketTypesData = await getEventTicketTypes(eventId);
      setTicketTypes(ticketTypesData);
    } catch (err) {
      console.error('刷新票種信息失敗:', err);
      toast.error('無法更新票券信息，請刷新頁面');
    }
  };
  
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="animate-pulse space-y-8">
            <div className="h-8 bg-gray-200 rounded w-1/3"></div>
            
            <div className="h-64 bg-gray-200 rounded"></div>
            
            <div className="space-y-4">
              <div className="h-4 bg-gray-200 rounded w-3/4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              <div className="h-4 bg-gray-200 rounded w-2/3"></div>
            </div>
          </div>
        </main>
      </div>
    );
  }
  
  if (error || !event) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="rounded-md bg-red-50 p-4">
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">錯誤</h3>
                <div className="mt-2 text-sm text-red-700">
                  <p>{error || '無法載入活動詳情'}</p>
                </div>
              </div>
            </div>
          </div>
          
          <div className="mt-6">
            <Link 
              href="/events"
              className="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-500"
            >
              <ArrowLeftIcon className="h-4 w-4 mr-1" />
              回到活動列表
            </Link>
          </div>
        </main>
      </div>
    );
  }
  
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          <div>
            <Link 
              href="/events"
              className="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-500"
            >
              <ArrowLeftIcon className="h-4 w-4 mr-1" />
              回到活動列表
            </Link>
            <h1 className="mt-2 text-3xl font-bold text-gray-900">{event.title}</h1>
          </div>
          
          <div className="bg-white rounded-lg overflow-hidden shadow">
            <div className="relative h-80 w-full">
              <Image
                src={`https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80`}
                alt={event.title}
                fill
                style={{ objectFit: 'cover' }}
              />
            </div>
            
            <div className="p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h2 className="text-xl font-semibold text-gray-900 mb-4">活動詳情</h2>
                  
                  <div className="prose max-w-none">
                    <p className="text-gray-700">{event.description}</p>
                  </div>
                  
                  <div className="mt-6 space-y-3">
                    <div className="flex items-center">
                      <CalendarIcon className="h-5 w-5 text-gray-400 mr-2" />
                      <span className="text-gray-600">{formatDate(event.start_time)}</span>
                    </div>
                    
                    <div className="flex items-center">
                      <ClockIcon className="h-5 w-5 text-gray-400 mr-2" />
                      <span className="text-gray-600">
                        {formatTime(event.start_time)} - {formatTime(event.end_time)}
                      </span>
                    </div>
                    
                    <div className="flex items-center">
                      <MapPinIcon className="h-5 w-5 text-gray-400 mr-2" />
                      <span className="text-gray-600">{event.location}</span>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h2 className="text-xl font-semibold text-gray-900 mb-4">票券選擇</h2>
                  
                  {ticketTypes.length > 0 ? (
                    <div className="space-y-4">
                      {ticketTypes.map(ticketType => (
                        <TicketTypeCard 
                          key={ticketType.id}
                          ticketType={ticketType}
                          eventId={event.id}
                          onPurchase={refreshTicketTypes}
                        />
                      ))}
                    </div>
                  ) : (
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <p className="text-gray-700">此活動目前沒有可用的票券。</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
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
