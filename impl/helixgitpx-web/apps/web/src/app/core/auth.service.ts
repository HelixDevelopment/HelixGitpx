import { Injectable, signal } from '@angular/core';
import { createConnectTransport } from '@connectrpc/connect-web';
import { createClient } from '@connectrpc/connect';
import { AuthService as AuthSvc } from '../../../../../libs/proto/src/helixgitpx/auth/v1/auth_connect';

const KEYCLOAK_URL = 'https://keycloak.helix.local/realms/helixgitpx';
const CLIENT_ID = 'helixgitpx-web';
const REDIRECT_URI = `${window.location.origin}/auth/callback`;

@Injectable({ providedIn: 'root' })
export class AuthSession {
  readonly user = signal<{ email: string } | null>(null);

  readonly client = createClient(AuthSvc, createConnectTransport({
    baseUrl: 'https://auth.helix.local',
    credentials: 'include',
  }));

  startOIDC() {
    const params = new URLSearchParams({
      response_type: 'code',
      client_id: CLIENT_ID,
      redirect_uri: REDIRECT_URI,
      scope: 'openid email profile',
      code_challenge_method: 'S256',
      code_challenge: 'placeholder',
    });
    window.location.href = `${KEYCLOAK_URL}/protocol/openid-connect/auth?${params}`;
  }

  async exchangeCode(code: string): Promise<void> {
    // The auth-service REST callback at /v1/auth/callback performs the
    // OIDC exchange + sets cookies. We just navigate there.
    window.location.href = `https://auth.helix.local/v1/auth/callback?code=${encodeURIComponent(code)}`;
  }

  async whoAmI(): Promise<{ email: string } | null> {
    try {
      const u = await this.client.whoAmI({});
      this.user.set({ email: u.email });
      return { email: u.email };
    } catch {
      this.user.set(null);
      return null;
    }
  }
}
