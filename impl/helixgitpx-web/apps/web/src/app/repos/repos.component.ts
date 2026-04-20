import { Component } from '@angular/core';

@Component({
  selector: 'hx-repos',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Repositories</h1>
    <p>List + code browser (connects to repo-service in M6 full build).</p>
  </main>`,
})
export class ReposComponent {}
