"use client";

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { TicketType } from '@/api/events';
import { createTicketType, updateTicketType } from '@/api/admin';

interface TicketTypeFormProps {
  eventId: string;
  ticketType?: TicketType;
  onSuccess: () => void;
  onCancel?: () => void;
}

interface TicketTypeFormData {
  name: string;
  price: number;
  total_quantity: number;
  available_quantity?: number;
  sale_start: string;
  sale_end: string;
}

export default function TicketTypeForm({ eventId, ticketType, onSuccess, onCancel }: TicketTypeFormProps) {
  const [isLoading, setIsLoading] = useState(false);
  const isEditing = !!ticketType;

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<TicketTypeFormData>({
    defaultValues: ticketType
      ? {
          name: ticketType.name,
          price: ticketType.price,
          total_quantity: ticketType.total_quantity,
          available_quantity: ticketType.available_quantity,
          sale_start: new Date(ticketType.sale_start).toISOString().slice(0, 16),
          sale_end: new Date(ticketType.sale_end).toISOString().slice(0, 16),
        }
      : {
          available_quantity: undefined,
        },
  });

  const totalQuantity = watch('total_quantity');

  const onSubmit = async (data: TicketTypeFormData) => {
    try {
      setIsLoading(true);
      
      // 若未設定可用數量，則使用總數量
      if (!data.available_quantity) {
        data.available_quantity = data.total_quantity;
      }

      if (isEditing && ticketType) {
        await updateTicketType(ticketType.id, {
          ...data,
          sale_start: new Date(data.sale_start).toISOString(),
          sale_end: new Date(data.sale_end).toISOString(),
        });
        toast.success('票種已更新');
      } else {
        await createTicketType(eventId, {
          ...data,
          sale_start: new Date(data.sale_start).toISOString(),
          sale_end: new Date(data.sale_end).toISOString(),
        });
        toast.success('票種已創建');
      }
      
      onSuccess();
    } catch (error) {
      console.error('保存票種失敗:', error);
      toast.error('操作失敗，請稍後再試');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">
          票種名稱
        </label>
        <input
          type="text"
          id="name"
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
          {...register('name', {
            required: '請輸入票種名稱',
            maxLength: {
              value: 100,
              message: '票種名稱不能超過 100 個字符',
            },
          })}
        />
        {errors.name && (
          <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="price" className="block text-sm font-medium text-gray-700">
          票價 (NTD)
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span className="text-gray-500 sm:text-sm">$</span>
          </div>
          <input
            type="number"
            id="price"
            className="pl-7 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            min="0"
            step="1"
            {...register('price', {
              required: '請輸入票價',
              valueAsNumber: true,
              min: {
                value: 0,
                message: '票價不能為負數',
              },
            })}
          />
        </div>
        {errors.price && (
          <p className="mt-1 text-sm text-red-600">{errors.price.message}</p>
        )}
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="total_quantity" className="block text-sm font-medium text-gray-700">
            總數量
          </label>
          <input
            type="number"
            id="total_quantity"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            min="1"
            step="1"
            {...register('total_quantity', {
              required: '請輸入總數量',
              valueAsNumber: true,
              min: {
                value: 1,
                message: '總數量必須至少為 1',
              },
            })}
          />
          {errors.total_quantity && (
            <p className="mt-1 text-sm text-red-600">{errors.total_quantity.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="available_quantity" className="block text-sm font-medium text-gray-700">
            可用數量 (留空則等於總數量)
          </label>
          <input
            type="number"
            id="available_quantity"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            min="0"
            max={totalQuantity}
            step="1"
            {...register('available_quantity', {
              valueAsNumber: true,
              min: {
                value: 0,
                message: '可用數量不能為負數',
              },
              max: {
                value: totalQuantity || Number.MAX_SAFE_INTEGER,
                message: '可用數量不能超過總數量',
              },
            })}
          />
          {errors.available_quantity && (
            <p className="mt-1 text-sm text-red-600">{errors.available_quantity.message}</p>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="sale_start" className="block text-sm font-medium text-gray-700">
            開售時間
          </label>
          <input
            type="datetime-local"
            id="sale_start"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            {...register('sale_start', {
              required: '請選擇開售時間',
            })}
          />
          {errors.sale_start && (
            <p className="mt-1 text-sm text-red-600">{errors.sale_start.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="sale_end" className="block text-sm font-medium text-gray-700">
            結束銷售時間
          </label>
          <input
            type="datetime-local"
            id="sale_end"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
            {...register('sale_end', {
              required: '請選擇結束銷售時間',
            })}
          />
          {errors.sale_end && (
            <p className="mt-1 text-sm text-red-600">{errors.sale_end.message}</p>
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
          {isLoading ? '處理中...' : isEditing ? '更新票種' : '創建票種'}
        </button>
      </div>
    </form>
  );
}
