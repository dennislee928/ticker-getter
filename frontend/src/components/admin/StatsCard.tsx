"use client";

import { ReactNode } from 'react';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: ReactNode;
  change?: {
    value: number;
    positive: boolean;
  };
  textColor?: string;
  bgColor?: string;
}

export default function StatsCard({ 
  title, 
  value, 
  icon,
  change,
  textColor = 'text-primary-600',
  bgColor = 'bg-primary-50'
}: StatsCardProps) {
  return (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className={`flex-shrink-0 rounded-md ${bgColor} p-3`}>
            {icon}
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">
                {title}
              </dt>
              <dd>
                <div className={`text-lg font-medium ${textColor}`}>
                  {value}
                </div>
              </dd>
            </dl>
          </div>
        </div>
      </div>
      {change && (
        <div className="bg-gray-50 px-5 py-3">
          <div className="text-sm">
            <span 
              className={`font-medium ${change.positive ? 'text-green-600' : 'text-red-600'} mr-1`}
            >
              {change.positive ? '+' : ''}{change.value}%
            </span>
            <span className="text-gray-500">
              較上月
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
