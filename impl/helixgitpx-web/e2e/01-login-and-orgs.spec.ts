import { test, expect } from '@playwright/test';

// Constitution §II §2: e2e tests NEVER use mocks. This spec assumes a live
// HelixGitpx stack reachable at HELIXGITPX_WEB_URL with a seeded user whose
// credentials live in environment variables.

const USER = process.env.HELIXGITPX_E2E_USER ?? '';
const PASSWORD = process.env.HELIXGITPX_E2E_PASSWORD ?? '';

test.beforeAll(() => {
  test.skip(!USER || !PASSWORD, 'HELIXGITPX_E2E_USER and HELIXGITPX_E2E_PASSWORD must be set');
});

test('redirects anonymous visitor to /login', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveURL(/\/login/);
});

test('login → dashboard → org list', async ({ page }) => {
  await page.goto('/login');

  // Click the Keycloak-backed "Continue" button; fill in on the Keycloak page.
  await page.getByRole('button', { name: /continue/i }).click();

  await page.getByLabel(/email/i).fill(USER);
  await page.getByLabel(/password/i).fill(PASSWORD);
  await page.getByRole('button', { name: /sign in/i }).click();

  await expect(page).toHaveURL(/\/dashboard/);

  await page.getByRole('link', { name: /orgs/i }).click();
  await expect(page).toHaveURL(/\/orgs/);
  // At minimum the user sees their own seeded "default" org.
  await expect(page.getByText(/default/i)).toBeVisible();
});

test('trust center loads without auth', async ({ page }) => {
  await page.goto('/trust');
  await expect(page.getByRole('heading', { name: /trust/i })).toBeVisible();
});
