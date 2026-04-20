import { Component, inject } from '@angular/core';
import { AuthSession } from '../core/auth.service';

@Component({
  selector: 'hx-login',
  standalone: true,
  template: `
    <main style="max-width: 24rem; margin: 4rem auto; text-align: center; font-family: system-ui;">
      <h1>HelixGitpx</h1>
      <p>Sign in to continue.</p>
      <button (click)="signIn()" style="padding: 0.75rem 1.5rem; font-size: 1rem;">
        Sign in with HelixGitpx
      </button>
    </main>
  `,
})
export class LoginComponent {
  private auth = inject(AuthSession);
  signIn() { this.auth.startOIDC(); }
}
