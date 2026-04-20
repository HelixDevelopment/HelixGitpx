import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  standalone: true,
  imports: [CommonModule],
  selector: 'app-trust-center',
  template: `
    <section class="trust">
      <h1>HelixGitpx Trust Center</h1>
      <p>Live status, security attestations, and compliance posture.</p>

      <h2>Status</h2>
      <p><a href="https://status.helixgitpx.io">status.helixgitpx.io</a></p>

      <h2>Certifications &amp; audits</h2>
      <ul>
        <li>SOC 2 Type I — report available under NDA.</li>
        <li>ISO 27001 — gap analysis complete; certification in progress.</li>
        <li>Annual third-party pen-test; latest executive summary on request.</li>
      </ul>

      <h2>Policies</h2>
      <ul>
        <li><a href="https://helixgitpx.io/privacy">Privacy policy</a></li>
        <li><a href="https://helixgitpx.io/security">Security policy</a></li>
        <li><a href="https://helixgitpx.io/dpa">Data Processing Agreement (DPA)</a></li>
        <li><a href="https://helixgitpx.io/subprocessors">Subprocessors</a></li>
      </ul>

      <h2>Report a vulnerability</h2>
      <p><a href="https://hackerone.com/helixgitpx">Bug bounty on HackerOne</a> or
      <a href="mailto:security&#64;helixgitpx.io">security&#64;helixgitpx.io</a>.</p>
    </section>
  `,
  styles: [`.trust { max-width: 800px; padding: 2rem; }`],
})
export class TrustCenterComponent {}
