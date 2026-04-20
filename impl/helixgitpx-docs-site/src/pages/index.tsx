import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';

export default function Home(): JSX.Element {
  return (
    <Layout title="HelixGitpx" description="Federated Git proxy with AI, search, and policy.">
      <main style={{padding: '4rem 2rem', maxWidth: 800, margin: '0 auto'}}>
        <h1>HelixGitpx</h1>
        <p>One namespace across many Git hosts. AI-assisted. Policy-driven.</p>
        <p>
          <Link className="button button--primary button--lg" to="/docs/intro">Get started →</Link>
        </p>
      </main>
    </Layout>
  );
}
