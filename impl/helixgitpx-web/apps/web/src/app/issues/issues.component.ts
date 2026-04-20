import { Component } from '@angular/core';

@Component({
  selector: 'hx-issues',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Issues</h1>
    <p>Issue flows backed by collab-service CRDT metadata (M5).</p>
  </main>`,
})
export class IssuesComponent {}
