import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";

const config: Config = {
  title: "HelixGitpx",
  tagline: "Helix Git Proxy eXtended — federated, privacy-preserving Git proxy",
  favicon: "img/favicon.ico",
  url: "https://docs.helixgitpx.dev",
  baseUrl: "/",
  organizationName: "helixgitpx",
  projectName: "helixgitpx",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  presets: [
    ["classic", {
      docs: {
        sidebarPath: "./sidebars.ts",
        editUrl: "https://github.com/helixgitpx/helixgitpx/edit/main/docs/specifications/main/main_implementation_material/HelixGitpx/",
      },
      theme: { customCss: "./src/css/custom.css" },
    }],
  ],
  themeConfig: {
    navbar: {
      title: "HelixGitpx",
      items: [
        { to: "/docs/intro", label: "Docs", position: "left" },
        { to: "/docs/roadmap/17-milestones", label: "Roadmap", position: "left" },
        { href: "https://github.com/helixgitpx/helixgitpx", label: "GitHub", position: "right" },
      ],
    },
    prism: { theme: prismThemes.github, darkTheme: prismThemes.dracula },
  },
};
export default config;
