// Pure utility for residency zone validation. Mirrors the Go
// domain.Residency enum in services/orgteam/internal/domain/residency.go.
// Intentionally framework-free so it is unit-testable in plain ts-jest
// without Angular's ESM dependencies.
export type Residency = 'EU' | 'UK' | 'US';

export const RESIDENCIES: readonly Residency[] = ['EU', 'UK', 'US'] as const;

export function isValidResidency(value: unknown): value is Residency {
  return typeof value === 'string' && (RESIDENCIES as readonly string[]).includes(value);
}

export class InvalidResidencyError extends Error {
  constructor(value: unknown) {
    super(`invalid residency: ${String(value)}`);
    this.name = 'InvalidResidencyError';
  }
}

export function normalizeResidency(value: string | null | undefined): Residency {
  const candidate = (value ?? '').toUpperCase();
  if (!isValidResidency(candidate)) {
    throw new InvalidResidencyError(value);
  }
  return candidate;
}
