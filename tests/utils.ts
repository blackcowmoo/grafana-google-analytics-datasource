import type { Page } from '@playwright/test';

// Grafana 13+ shows a "What's new in Grafana" announcement carousel on the
// first authenticated page load. It mounts into #grafana-portal-container
// with a viewport-sized scrim (z-index:1050) that intercepts pointer events,
// so any subsequent click — including saveAndTest — times out. Dismiss it
// before interacting with the page. Older Grafana versions don't render the
// dialog, so the visibility check returns quickly and is a no-op.
export const dismissWhatsNewModal = async (page: Page) => {
  const dialog = page.getByRole('dialog', { name: /What's new in Grafana/i });
  if (await dialog.isVisible({ timeout: 1500 }).catch(() => false)) {
    await dialog.getByRole('button', { name: 'Close' }).first().click();
    await dialog.waitFor({ state: 'hidden', timeout: 5000 });
  }
};
