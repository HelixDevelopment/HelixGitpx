import { Component } from '@angular/core';

@Component({
  selector: 'hx-conflicts',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Conflicts inbox</h1>
    <p>Streams from conflict-resolver's conflict.resolved topic via live-events-service.</p>
  </main>`,
})
export class ConflictsComponent {}
