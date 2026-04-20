import { Component } from '@angular/core';

@Component({
  selector: 'hx-prs',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Pull Requests</h1>
    <p>List, detail, diff, review — bound to adapter-pool's CreatePR.</p>
  </main>`,
})
export class PrsComponent {}
