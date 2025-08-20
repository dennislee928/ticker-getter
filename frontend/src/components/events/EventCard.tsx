"use client";

import Link from 'next/link';
import Image from 'next/image';
import { formatDate, formatPrice } from '@/lib/utils';
import { CalendarIcon, MapPinIcon, CurrencyDollarIcon } from '@heroicons/react/24/outline';
import { Event } from '@/api/events';

interface EventCardProps {
  event: Event;
  minPrice?: number;
}

export default function EventCard({ event, minPrice }: EventCardProps) {
  return (
    <div className="flex flex-col rounded-lg shadow-md overflow-hidden bg-white">
      <div className="relative h-48 w-full">
        <Image
          src={`https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80`}
          alt={event.title}
          fill
          style={{ objectFit: 'cover' }}
        />
      </div>

      <div className="flex-1 p-4 flex flex-col">
        <h3 className="text-xl font-semibold text-gray-900">{event.title}</h3>
        
        <div className="mt-2 flex items-center text-sm text-gray-600">
          <CalendarIcon className="mr-1.5 h-5 w-5 text-gray-400" />
          {formatDate(event.start_time)}
        </div>
        
        <div className="mt-1 flex items-center text-sm text-gray-600">
          <MapPinIcon className="mr-1.5 h-5 w-5 text-gray-400" />
          {event.location}
        </div>
        
        {minPrice && (
          <div className="mt-1 flex items-center text-sm font-medium text-primary-600">
            <CurrencyDollarIcon className="mr-1.5 h-5 w-5 text-primary-500" />
            起價 {formatPrice(minPrice)}
          </div>
        )}
        
        <div className="mt-2 text-sm text-gray-500 line-clamp-2">
          {event.description}
        </div>
        
        <div className="mt-auto pt-4">
          <Link
            href={`/events/${event.id}`}
            className="w-full inline-flex justify-center items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            查看詳情
          </Link>
        </div>
      </div>
    </div>
  );
}
