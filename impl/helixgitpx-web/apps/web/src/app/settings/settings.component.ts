import { Component } from '@angular/core';

@Component({
  selector: 'hx-settings',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Settings</h1>
    <p>Org admin, members, upstream config (wired to upstream-service from M4).</p>
  </main>`,
})
export class SettingsComponent {}
