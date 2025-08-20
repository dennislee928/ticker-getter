import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { Toaster } from 'react-hot-toast'
import { TlsFingerprintProvider } from '@/components/TlsFingerprintProvider'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: '票務購買平台',
  description: '線上購買活動票券的便捷平台',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-Hant-TW">
      <body className={inter.className}>
        <TlsFingerprintProvider>
          <Toaster position="top-center" />
          {children}
        </TlsFingerprintProvider>
      </body>
    </html>
  )
}