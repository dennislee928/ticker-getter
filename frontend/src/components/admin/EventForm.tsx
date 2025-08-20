"use client";

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { Event } from '@/api/events';
import { createEvent, updateEvent } from '@/api/admin';

interface EventFormProps {
  event?: Event;
  onSuccess: () => void;
  onCancel?: () => void;
}

interface EventFormData {
  title: string;
  description: string;
  location: string;
  start_time: string;
  end_time: string;
}

export default function EventForm({ event, onSuccess, onCancel }: EventFormProps) {
  const [isLoading, setIsLoading] = useState(false);
  const isEditing = !!event;

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<EventFormData>({
    defaultValues: event
      ? {
          title: event.title,
          description: event.description,
          location: event.location,
          start_time: new Date(event.start_time).toISOString().slice(0, 16),
          end_time: new Date(event.end_time).toISOString().slice(0, 16),
        }
      : undefined,
  });

  const onSubmit = async (data: EventFormData) => {
    try {
      setIsLoading(true);
      
      if (isEditing && event) {
        await updateEvent(event.id, {
          ...data,
          start_time: new Date(data.start_time).toISOString(),
          end_time: new Date(data.end_time).toISOString(),
        });
        toast.success('活動已更新');
      } else {
        await createEvent({
          ...data,
          start_time: new Date(data.start_time).toISOString(),
          end_time: new Date(data.end_time).toISOString(),
        });
        toast.success('活動已創建');
      }
      
      onSuccess();
    } catch (error) {
      console.error('保存活動失敗:', error);
      toast.error('操作失敗，請稍後再試');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label htmlFor="title" className="block text-sm font-medium text-gray-700">
          活動名稱
        </label>
        <input
          type="text"
          id="title"
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
          {...register('title', {
            required: '請輸入活動名稱',
            maxLength: {
              value: 100,
              message: '活動名稱不能超過 100 個字符',
            },
          })}
        />
        {errors.title && (
          <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="description" className="block text-sm font-medium text-gray-700">
          活動描述
        </label>
        <textarea
          id="description"
          rows={3}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
          {...register('description', {
            required: '請輸入活動描述',
          })}
        />
        {errors.description && (
          <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="location" className="block text-sm font-medium text-gray-700">
          活動地點
        </label>
        <input
          type="text"
          id="location"
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
          {...register('location', {
            required: '請輸入活動地點',
          })}
        />
        {errors.location && (
          <p className="mt-1 text-sm text-red-600">{errors.location.message}</p>
        )}
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="start_time" className="block text-sm font-medium text-gray-700">
            開始時間
          </label>
          <input
            type="datetime-local"
            id="start_time"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            {...register('start_time', {
              required: '請選擇開始時間',
            })}
          />
          {errors.start_time && (
            <p className="mt-1 text-sm text-red-600">{errors.start_time.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="end_time" className="block text-sm font-medium text-gray-700">
            結束時間
          </label>
          <input
            type="datetime-local"
            id="end_time"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            {...register('end_time', {
              required: '請選擇結束時間',
            })}
          />
          {errors.end_time && (
            <p className="mt-1 text-sm text-red-600">{errors.end_time.message}</p>
          )}
        </div>
      </div>

      <div className="flex justify-end space-x-3">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            取消
          </button>
        )}
        <button
          type="submit"
          disabled={isLoading}
          className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
        >
          {isLoading ? '處理中...' : isEditing ? '更新活動' : '創建活動'}
        </button>
      </div>
    </form>
  );
}
