import Layout from '@theme/Layout';

export default function Trust(): JSX.Element {
  return (
    <Layout title="Trust Center">
      <main style={{padding: '3rem 2rem', maxWidth: 800, margin: '0 auto'}}>
        <h1>Trust Center</h1>
        <p>Status page, certifications, and compliance posture.</p>
        <ul>
          <li><a href="https://status.helixgitpx.io">status.helixgitpx.io</a></li>
          <li>SOC 2 Type I — report under NDA.</li>
          <li>ISO 27001 — gap analysis complete; certification in progress.</li>
          <li><a href="https://hackerone.com/helixgitpx">Bug bounty program</a></li>
        </ul>
      </main>
    </Layout>
  );
}
