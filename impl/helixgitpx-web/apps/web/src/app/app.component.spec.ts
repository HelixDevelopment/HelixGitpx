import { AppComponent } from './app.component';
import { routes } from './routes';

describe('AppComponent', () => {
  it('is constructable as a standalone component class', () => {
    const instance = new AppComponent();
    expect(instance).toBeDefined();
    expect(instance instanceof AppComponent).toBe(true);
  });
});

describe('App routing structure', () => {
  it('defines exactly the expected number of routes', () => {
    expect(routes.length).toBe(12);
  });

  it('redirects root path to dashboard', () => {
    const root = routes.find(r => r.path === '');
    expect(root).toBeDefined();
    expect(root!.redirectTo).toBe('/dashboard');
    expect(root!.pathMatch).toBe('full');
  });

  it('has unauthenticated routes without guards', () => {
    const unauthenticated = routes.filter(r => r.path === 'login' || r.path === 'auth/callback' || r.path === 'trust');
    expect(unauthenticated.length).toBe(3);
    for (const route of unauthenticated) {
      expect(route.canActivate).toBeUndefined();
    }
  });

  it('guards all authenticated routes with authGuard', () => {
    const guarded = routes.filter(r =>
      r.path !== '' &&
      r.path !== 'login' &&
      r.path !== 'auth/callback' &&
      r.path !== 'trust',
    );
    expect(guarded.length).toBe(7);
    for (const route of guarded) {
      expect(route.canActivate).toBeDefined();
      expect(route.canActivate!.length).toBeGreaterThan(0);
    }
  });

  it('includes all required feature routes', () => {
    const paths = routes.filter(r => r.path !== '').map(r => r.path);
    const requiredPaths = ['login', 'auth/callback', 'dashboard', 'orgs', 'repos', 'prs', 'issues', 'conflicts', 'search', 'settings', 'trust'];
    for (const required of requiredPaths) {
      expect(paths).toContain(required);
    }
  });
});
