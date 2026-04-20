import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  mainSidebar: [
    'intro',
    { type: 'category', label: 'Getting Started', items: ['getting-started/install', 'getting-started/first-repo'] },
    { type: 'category', label: 'Concepts', items: ['concepts/federation', 'concepts/conflicts', 'concepts/ai'] },
    { type: 'category', label: 'Operations', items: ['operations/runbooks', 'operations/slo'] },
    { type: 'category', label: 'Reference', items: ['reference/api', 'reference/cli'] },
    { type: 'category', label: 'Legal', items: ['legal/terms', 'legal/privacy', 'legal/dpa'] },
  ],
};

export default sidebars;
