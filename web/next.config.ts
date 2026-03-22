import createNextIntlPlugin from 'next-intl/plugin';
import type { NextConfig } from "next";

const withNextIntl = createNextIntlPlugin();

const nextConfig: NextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/v1/auth/:path*',
        destination: 'http://localhost:8081/api/v1/auth/:path*',
      },
      {
        source: '/api/v1/invite/:path*',
        destination: 'http://localhost:8081/api/v1/invite/:path*',
      },
      {
        source: '/api/v1/jobs/:path*',
        destination: 'http://localhost:8082/api/v1/jobs/:path*',
      },
      {
        source: '/api/v1/candidates/:path*',
        destination: 'http://localhost:8082/api/v1/candidates/:path*',
      },
      {
        source: '/api/v1/departments/:path*',
        destination: 'http://localhost:8082/api/v1/departments/:path*',
      },
      {
        source: '/api/v1/chat/:path*',
        destination: 'http://localhost:8082/api/v1/chat/:path*',
      },
      {
        source: '/api/v1/interview/:path*',
        destination: 'http://localhost:8082/api/v1/interview/:path*',
      },
      {
        source: '/api/v1/dashboard/:path*',
        destination: 'http://localhost:8082/api/v1/dashboard/:path*',
      },
      {
        source: '/api/v1/ai-settings/:path*',
        destination: 'http://localhost:8082/api/v1/ai-settings/:path*',
      },
      {
        source: '/api/v1/ai/:path*',
        destination: 'http://localhost:8082/api/v1/ai/:path*',
      },
    ];
  },
};

export default withNextIntl(nextConfig);
