import { Component } from '@angular/core';

@Component({
  selector: 'hx-search',
  standalone: true,
  template: `<main style="max-width:64rem;margin:2rem auto;font-family:system-ui;">
    <h1>Search</h1>
    <p>Hybrid search (Meilisearch + Qdrant + OpenSearch fan-out) lands in M7.</p>
  </main>`,
})
export class SearchComponent {}
