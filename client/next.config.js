/** @type {import('next').NextConfig} */
const nextConfig = {
    images: {
        domains: ["lh3.googleusercontent.com"],
      },
    async redirects() {
        return [
            {
                source: '/',
                missing : [
                    {
                        type: 'cookie',
                        key: 'task-together-session'
                    }
                ],
                permanent: false,
                destination: '/login',
            },
        ]
    }
}

module.exports = nextConfig
