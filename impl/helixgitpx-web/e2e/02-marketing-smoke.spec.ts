import { test, expect } from '@playwright/test';

// These tests exercise the marketing website and the web app's `/trust`
// route — neither requires a backend. They run out of the box when the
// marketing site is served locally (`npm --prefix ../helixgitpx-website run dev`).
// Use HELIXGITPX_MARKETING_URL to target a deployed site instead.

const marketingUrl = process.env.HELIXGITPX_MARKETING_URL ?? 'http://localhost:4321';

test.describe('marketing site — no backend required', () => {
  test('home renders the hero copy', async ({ page }) => {
    const response = await page.goto(marketingUrl).catch(() => null);
    test.skip(!response || !response.ok(), 'marketing site not reachable; start `npm run dev` first');

    await expect(page.getByRole('heading', { name: /one namespace/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /start free/i }).first()).toBeVisible();
  });

  test('trust center disclaims non-existent audits', async ({ page }) => {
    const response = await page.goto(marketingUrl + '/trust/').catch(() => null);
    test.skip(!response || !response.ok(), 'trust page not reachable');

    // The honest copy change mentions "no report issued".
    await expect(page.getByText(/no report issued, no auditor engaged/i)).toBeVisible();
  });

  test('customers page tells the truth', async ({ page }) => {
    const response = await page.goto(marketingUrl + '/customers/').catch(() => null);
    test.skip(!response || !response.ok(), 'customers page not reachable');

    await expect(page.getByText(/no hosted service today|no hosted customers today/i)).toBeVisible();
  });

  test('changelog lists v1.0.0', async ({ page }) => {
    const response = await page.goto(marketingUrl + '/changelog/').catch(() => null);
    test.skip(!response || !response.ok(), 'changelog not reachable');

    await expect(page.getByText(/v1\.0\.0/)).toBeVisible();
  });
});
