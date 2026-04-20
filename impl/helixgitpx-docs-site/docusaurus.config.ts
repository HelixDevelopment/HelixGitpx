import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'HelixGitpx',
  tagline: 'Federated Git proxy with AI, search, and policy.',
  url: 'https://docs.helixgitpx.io',
  baseUrl: '/',
  organizationName: 'HelixGitpx',
  projectName: 'HelixGitpx',
  favicon: 'img/favicon.ico',
  presets: [
    ['classic', {
      docs: { sidebarPath: './sidebars.ts', editUrl: 'https://github.com/HelixGitpx/HelixGitpx/tree/main/impl/helixgitpx-docs-site/' },
      blog: { showReadingTime: true },
      theme: { customCss: './src/css/custom.css' },
    } satisfies Preset.Options],
  ],
  themeConfig: {
    navbar: {
      title: 'HelixGitpx',
      items: [
        { type: 'docSidebar', sidebarId: 'mainSidebar', position: 'left', label: 'Docs' },
        { to: '/trust', label: 'Trust', position: 'left' },
        { href: 'https://github.com/HelixGitpx/HelixGitpx', label: 'GitHub', position: 'right' },
      ],
    },
    prism: { theme: prismThemes.github, darkTheme: prismThemes.dracula },
  } satisfies Preset.ThemeConfig,
};

export default config;
