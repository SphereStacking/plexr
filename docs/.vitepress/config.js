import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Plexr',
  description: 'Developer-friendly CLI tool for automating local development environment setup',
  base: '/plexr/', // This should match your repository name
  
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/getting-started' },
          { text: 'API', link: '/api/' },
          { text: 'Examples', link: '/examples/' }
        ],
        
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'What is Plexr?', link: '/' },
              { text: 'Getting Started', link: '/guide/getting-started' }
            ]
          },
          {
            text: 'Guide',
            items: [
              { text: 'Installation', link: '/guide/installation' },
              { text: 'Configuration', link: '/guide/configuration' },
              { text: 'Commands', link: '/guide/commands' },
              { text: 'Executors', link: '/guide/executors' },
              { text: 'State Management', link: '/guide/state-management' }
            ]
          },
          {
            text: 'API Reference',
            items: [
              { text: 'CLI Commands', link: '/api/cli-commands' },
              { text: 'Configuration Schema', link: '/api/configuration-schema' },
              { text: 'Executors API', link: '/api/executors' }
            ]
          },
          {
            text: 'Examples',
            items: [
              { text: 'Basic Setup', link: '/examples/basic-setup' },
              { text: 'Advanced Patterns', link: '/examples/advanced-patterns' },
              { text: 'Real World', link: '/examples/real-world' }
            ]
          }
        ]
      }
    },
    ja: {
      label: '日本語',
      lang: 'ja',
      link: '/ja/',
      themeConfig: {
        nav: [
          { text: 'ホーム', link: '/ja/' },
          { text: 'ガイド', link: '/ja/guide/getting-started' },
          { text: 'API', link: '/ja/api/' },
          { text: '例', link: '/ja/examples/' }
        ],
        
        sidebar: [
          {
            text: 'はじめに',
            items: [
              { text: 'Plexrとは？', link: '/ja/' },
              { text: 'はじめる', link: '/ja/guide/getting-started' }
            ]
          },
          {
            text: 'ガイド',
            items: [
              { text: 'インストール', link: '/ja/guide/installation' },
              { text: '設定', link: '/ja/guide/configuration' },
              { text: 'コマンド', link: '/ja/guide/commands' },
              { text: 'エグゼキューター', link: '/ja/guide/executors' },
              { text: 'ステート管理', link: '/ja/guide/state-management' }
            ]
          },
          {
            text: 'APIリファレンス',
            items: [
              { text: 'CLIコマンド', link: '/ja/api/cli-commands' },
              { text: '設定スキーマ', link: '/ja/api/configuration-schema' },
              { text: 'エグゼキューターAPI', link: '/ja/api/executors' }
            ]
          },
          {
            text: '例',
            items: [
              { text: '基本セットアップ', link: '/ja/examples/basic-setup' },
              { text: '高度なパターン', link: '/ja/examples/advanced-patterns' },
              { text: '実世界の例', link: '/ja/examples/real-world' }
            ]
          }
        ]
      }
    }
  },
  
  themeConfig: {
    logo: '🚀',
    
    socialLinks: [
      { icon: 'github', link: 'https://github.com/SphereStacking/plexr' }
    ],
    
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2023-present SphereStacking'
    },
    
    search: {
      provider: 'local'
    },
    
    editLink: {
      pattern: 'https://github.com/SphereStacking/plexr/edit/main/docs/:path'
    }
  }
})