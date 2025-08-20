"use client";

import {
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  CurrencyDollarIcon,
  UserIcon,
  TicketIcon,
} from '@heroicons/react/24/outline';
import { formatDate } from '@/lib/utils';

interface Activity {
  id: string;
  type: 'order' | 'ticket' | 'user' | 'refund';
  status: 'success' | 'pending' | 'failed';
  message: string;
  timestamp: string;
  amount?: number;
  user?: string;
}

interface RecentActivitiesProps {
  activities: Activity[];
}

export default function RecentActivities({ activities }: RecentActivitiesProps) {
  const getIcon = (type: string) => {
    switch (type) {
      case 'order':
        return <CurrencyDollarIcon className="h-5 w-5 text-green-500" />;
      case 'ticket':
        return <TicketIcon className="h-5 w-5 text-blue-500" />;
      case 'user':
        return <UserIcon className="h-5 w-5 text-purple-500" />;
      case 'refund':
        return <CurrencyDollarIcon className="h-5 w-5 text-red-500" />;
      default:
        return <ClockIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircleIcon className="h-5 w-5 text-green-500" />;
      case 'pending':
        return <ClockIcon className="h-5 w-5 text-yellow-500" />;
      case 'failed':
        return <XCircleIcon className="h-5 w-5 text-red-500" />;
      default:
        return null;
    }
  };

  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:px-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900">
          最近活動
        </h3>
      </div>
      <div className="border-t border-gray-200 divide-y divide-gray-200">
        {activities.length > 0 ? (
          <ul className="divide-y divide-gray-200">
            {activities.map((activity) => (
              <li key={activity.id} className="px-4 py-4 sm:px-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      {getIcon(activity.type)}
                    </div>
                    <div className="ml-3">
                      <p className="text-sm font-medium text-gray-900">
                        {activity.message}
                      </p>
                      {activity.user && (
                        <p className="text-xs text-gray-500">
                          使用者: {activity.user}
                        </p>
                      )}
                      {activity.amount !== undefined && (
                        <p className="text-xs text-gray-500">
                          金額: ${activity.amount}
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center">
                    <div className="flex-shrink-0 mr-2">
                      {getStatusIcon(activity.status)}
                    </div>
                    <div className="text-sm text-gray-500">
                      {formatDate(activity.timestamp)}
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <div className="px-4 py-5 sm:px-6 text-center text-gray-500">
            暫無活動記錄
          </div>
        )}
      </div>
      {activities.length > 0 && (
        <div className="px-4 py-4 sm:px-6 border-t border-gray-200">
          <a
            href="#"
            className="text-sm font-medium text-primary-600 hover:text-primary-500"
          >
            查看所有活動
          </a>
        </div>
      )}
    </div>
  );
}
