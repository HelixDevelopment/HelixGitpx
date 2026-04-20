import { Component, inject, OnInit, signal } from '@angular/core';
import { OrgTeamApi } from '../core/orgteam.service';

@Component({
  selector: 'hx-dashboard',
  standalone: true,
  template: `
    <main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
      <h1>Dashboard</h1>
      <section><h2>Recent activity</h2><p>(activity timeline — wire to live-events-service in M7)</p></section>
      <section><h2>Your organisations</h2>
        <ul>@for (o of orgs(); track o.id) { <li>{{o.slug}}</li> }</ul>
      </section>
    </main>
  `,
})
export class DashboardComponent implements OnInit {
  private api = inject(OrgTeamApi);
  readonly orgs = signal<any[]>([]);
  async ngOnInit() { const r = await this.api.orgs.list({}); this.orgs.set(r.orgs as any[]); }
}
