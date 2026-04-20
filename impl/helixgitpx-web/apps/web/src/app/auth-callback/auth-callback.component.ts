import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthSession } from '../core/auth.service';

@Component({
  selector: 'hx-auth-callback',
  standalone: true,
  template: `<p style="font-family: system-ui; padding: 2rem;">Signing you in…</p>`,
})
export class AuthCallbackComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private auth = inject(AuthSession);

  async ngOnInit() {
    const code = this.route.snapshot.queryParamMap.get('code');
    if (!code) {
      this.router.navigate(['/login']);
      return;
    }
    await this.auth.exchangeCode(code);
  }
}
