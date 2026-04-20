import { Component, inject, OnInit, signal } from '@angular/core';
import { OrgTeamApi } from '../core/orgteam.service';
import { AuthSession } from '../core/auth.service';

type Org = { id: string; slug: string; name: string };

@Component({
  selector: 'hx-orgs',
  standalone: true,
  template: `
    <main style="max-width: 48rem; margin: 2rem auto; font-family: system-ui;">
      <h1>Your organisations</h1>
      <p>Signed in as {{ auth.user()?.email }}</p>
      <ul>
        @for (org of orgs(); track org.id) {
          <li>{{ org.slug }} — {{ org.name }}</li>
        } @empty {
          <li>No organisations yet.</li>
        }
      </ul>
      <div style="margin-top: 2rem;">
        <h2>Create organisation</h2>
        <input #slug placeholder="slug" />
        <input #name placeholder="Name" />
        <button (click)="create(slug.value, name.value)">Create</button>
      </div>
    </main>
  `,
})
export class OrgsComponent implements OnInit {
  private api = inject(OrgTeamApi);
  protected auth = inject(AuthSession);
  readonly orgs = signal<Org[]>([]);

  async ngOnInit() {
    await this.refresh();
  }

  async refresh() {
    const resp = await this.api.orgs.list({});
    this.orgs.set(resp.orgs as Org[]);
  }

  async create(slug: string, name: string) {
    if (!slug || !name) return;
    await this.api.orgs.create({ slug, name });
    await this.refresh();
  }
}
