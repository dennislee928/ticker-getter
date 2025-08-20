"use client";

import { Fragment, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Dialog, Disclosure, Popover, Transition } from '@headlessui/react';
import {
  Bars3Icon,
  XMarkIcon,
  UserIcon,
  TicketIcon,
  CalendarIcon,
} from '@heroicons/react/24/outline';
import { ChevronDownIcon } from '@heroicons/react/20/solid';
import { cn } from '@/lib/utils';

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}

interface NavbarProps {
  isLoggedIn?: boolean;
  isAdmin?: boolean;
  userName?: string;
}

export default function Navbar({ isLoggedIn = false, isAdmin = false, userName }: NavbarProps) {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const pathname = usePathname();
  
  const navigation = [
    { name: '活動列表', href: '/events' },
    { name: '關於我們', href: '/about' },
    { name: '常見問題', href: '/faq' },
    { name: '聯絡我們', href: '/contact' },
  ];

  const userNavigation = isLoggedIn
    ? [
        { name: '我的訂單', href: '/account/orders', icon: TicketIcon },
        { name: '個人設定', href: '/account/settings', icon: UserIcon },
        { name: '登出', href: '/auth/logout', icon: XMarkIcon },
      ]
    : [
        { name: '登入', href: '/auth/login', icon: UserIcon },
        { name: '註冊', href: '/auth/register', icon: UserIcon },
      ];

  if (isLoggedIn && isAdmin) {
    userNavigation.unshift({
      name: '管理後台',
      href: '/admin',
      icon: CalendarIcon,
    });
  }

  return (
    <header className="bg-white shadow">
      <nav className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8" aria-label="Top">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Link href="/" className="flex-shrink-0">
              <span className="text-2xl font-bold text-primary-600">票務平台</span>
            </Link>
            <div className="hidden md:ml-10 md:block">
              <div className="flex space-x-4">
                {navigation.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={classNames(
                      pathname === item.href
                        ? 'bg-primary-50 text-primary-700'
                        : 'text-gray-600 hover:text-primary-700',
                      'rounded-md px-3 py-2 text-sm font-medium'
                    )}
                  >
                    {item.name}
                  </Link>
                ))}
              </div>
            </div>
          </div>
          
          <div className="hidden md:block">
            <div className="ml-4 flex items-center md:ml-6">
              {isLoggedIn ? (
                <Popover className="relative">
                  <Popover.Button className="flex items-center gap-x-1 text-sm font-semibold leading-6 text-gray-900">
                    <span className="text-gray-600">你好，{userName || '用戶'}</span>
                    <ChevronDownIcon className="h-5 w-5 flex-none text-gray-400" aria-hidden="true" />
                  </Popover.Button>

                  <Transition
                    as={Fragment}
                    enter="transition ease-out duration-200"
                    enterFrom="opacity-0 translate-y-1"
                    enterTo="opacity-100 translate-y-0"
                    leave="transition ease-in duration-150"
                    leaveFrom="opacity-100 translate-y-0"
                    leaveTo="opacity-0 translate-y-1"
                  >
                    <Popover.Panel className="absolute -right-8 top-full z-10 mt-3 w-56 rounded-xl bg-white p-2 shadow-lg ring-1 ring-gray-900/5">
                      {userNavigation.map((item) => (
                        <Link
                          key={item.name}
                          href={item.href}
                          className="block rounded-lg px-3 py-2 text-sm font-semibold leading-6 text-gray-900 hover:bg-gray-50"
                        >
                          <div className="flex items-center">
                            <item.icon className="h-5 w-5 mr-2" />
                            {item.name}
                          </div>
                        </Link>
                      ))}
                    </Popover.Panel>
                  </Transition>
                </Popover>
              ) : (
                <div className="flex space-x-4">
                  <Link
                    href="/auth/login"
                    className="text-gray-600 hover:text-primary-700 px-3 py-2 rounded-md text-sm font-medium"
                  >
                    登入
                  </Link>
                  <Link
                    href="/auth/register"
                    className="bg-primary-600 text-white hover:bg-primary-700 px-3 py-2 rounded-md text-sm font-medium"
                  >
                    註冊
                  </Link>
                </div>
              )}
            </div>
          </div>
          
          <div className="flex md:hidden">
            {/* Mobile menu button */}
            <button
              type="button"
              className="relative inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
              onClick={() => setMobileMenuOpen(true)}
            >
              <span className="absolute -inset-0.5" />
              <span className="sr-only">開啟選單</span>
              <Bars3Icon className="block h-6 w-6" aria-hidden="true" />
            </button>
          </div>
        </div>
      </nav>

      {/* Mobile menu */}
      <Dialog
        as="div"
        className="md:hidden"
        open={mobileMenuOpen}
        onClose={setMobileMenuOpen}
      >
        <div className="fixed inset-0 z-50" />
        <Dialog.Panel className="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10">
          <div className="flex items-center justify-between">
            <Link href="/" className="-m-1.5 p-1.5">
              <span className="text-2xl font-bold text-primary-600">票務平台</span>
            </Link>
            <button
              type="button"
              className="rounded-md p-2.5 text-gray-700"
              onClick={() => setMobileMenuOpen(false)}
            >
              <span className="sr-only">關閉選單</span>
              <XMarkIcon className="h-6 w-6" aria-hidden="true" />
            </button>
          </div>
          <div className="mt-6 flow-root">
            <div className="-my-6 divide-y divide-gray-500/10">
              <div className="space-y-2 py-6">
                {navigation.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    {item.name}
                  </Link>
                ))}
              </div>
              <div className="py-6">
                {userNavigation.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    className="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    <div className="flex items-center">
                      <item.icon className="h-5 w-5 mr-2" />
                      {item.name}
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </Dialog.Panel>
      </Dialog>
    </header>
  );
}
