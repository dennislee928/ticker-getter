"use client";

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';

interface RegisterFormData {
  name: string;
  email: string;
  phone: string;
  password: string;
  confirmPassword: string;
}

export default function Register() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<RegisterFormData>();
  
  const password = watch('password', '');
  
  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true);
    
    try {
      // 模擬 API 呼叫
      // 實際實現中應替換為真實的 API 呼叫
      console.log('註冊資料:', data);
      
      // 模擬延遲以示範載入狀態
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      toast.success('註冊成功！請登入您的帳戶。');
      router.push('/auth/login');
    } catch (error) {
      console.error('註冊錯誤:', error);
      toast.error('註冊失敗，請稍後再試');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <div className="flex min-h-screen flex-col justify-center py-12 sm:px-6 lg:px-8 bg-gray-50">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
          創建新帳戶
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          或{' '}
          <Link href="/auth/login" className="font-medium text-primary-600 hover:text-primary-500">
            登入您的帳戶
          </Link>
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <label htmlFor="name" className="block text-sm font-medium leading-6 text-gray-900">
                姓名
              </label>
              <div className="mt-2">
                <input
                  id="name"
                  type="text"
                  autoComplete="name"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                  {...register('name', {
                    required: '請輸入姓名',
                    minLength: {
                      value: 2,
                      message: '姓名至少需要 2 個字元',
                    },
                  })}
                />
                {errors.name && (
                  <p className="mt-2 text-sm text-red-600">{errors.name.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">
                電子郵件
              </label>
              <div className="mt-2">
                <input
                  id="email"
                  type="email"
                  autoComplete="email"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                  {...register('email', {
                    required: '請輸入電子郵件',
                    pattern: {
                      value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                      message: '無效的電子郵件地址',
                    },
                  })}
                />
                {errors.email && (
                  <p className="mt-2 text-sm text-red-600">{errors.email.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="phone" className="block text-sm font-medium leading-6 text-gray-900">
                電話號碼
              </label>
              <div className="mt-2">
                <input
                  id="phone"
                  type="tel"
                  autoComplete="tel"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                  {...register('phone', {
                    required: '請輸入電話號碼',
                    pattern: {
                      value: /^[0-9]{10}$/,
                      message: '請輸入有效的10位電話號碼',
                    },
                  })}
                />
                {errors.phone && (
                  <p className="mt-2 text-sm text-red-600">{errors.phone.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium leading-6 text-gray-900">
                密碼
              </label>
              <div className="mt-2">
                <input
                  id="password"
                  type="password"
                  autoComplete="new-password"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                  {...register('password', {
                    required: '請輸入密碼',
                    minLength: {
                      value: 8,
                      message: '密碼至少需要 8 個字元',
                    },
                    pattern: {
                      value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/,
                      message: '密碼必須包含至少一個大寫字母、一個小寫字母和一個數字',
                    },
                  })}
                />
                {errors.password && (
                  <p className="mt-2 text-sm text-red-600">{errors.password.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium leading-6 text-gray-900">
                確認密碼
              </label>
              <div className="mt-2">
                <input
                  id="confirmPassword"
                  type="password"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                  {...register('confirmPassword', {
                    required: '請確認密碼',
                    validate: value => value === password || '密碼不匹配',
                  })}
                />
                {errors.confirmPassword && (
                  <p className="mt-2 text-sm text-red-600">{errors.confirmPassword.message}</p>
                )}
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <input
                  id="agree-terms"
                  name="agree-terms"
                  type="checkbox"
                  className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                  required
                />
                <label htmlFor="agree-terms" className="ml-2 block text-sm text-gray-900">
                  我同意 <Link href="/terms" className="text-primary-600 hover:text-primary-500">使用條款</Link> 和 <Link href="/privacy" className="text-primary-600 hover:text-primary-500">隱私政策</Link>
                </label>
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={isLoading}
                className="flex w-full justify-center rounded-md bg-primary-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? '處理中...' : '註冊'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
