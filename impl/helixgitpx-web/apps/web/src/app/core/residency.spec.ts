import { isValidResidency, normalizeResidency, InvalidResidencyError, RESIDENCIES } from './residency';

describe('residency', () => {
  it('exports exactly the three GA zones', () => {
    expect(RESIDENCIES).toEqual(['EU', 'UK', 'US']);
  });

  it('accepts each valid zone', () => {
    for (const z of RESIDENCIES) {
      expect(isValidResidency(z)).toBe(true);
    }
  });

  it('rejects lowercase', () => {
    expect(isValidResidency('eu')).toBe(false);
  });

  it('rejects non-string inputs', () => {
    expect(isValidResidency(null)).toBe(false);
    expect(isValidResidency(undefined)).toBe(false);
    expect(isValidResidency(123)).toBe(false);
    expect(isValidResidency({})).toBe(false);
  });

  it('normalizes case-insensitive input', () => {
    expect(normalizeResidency('eu')).toBe('EU');
    expect(normalizeResidency('Uk')).toBe('UK');
  });

  it('throws InvalidResidencyError for unknown values', () => {
    expect(() => normalizeResidency('DE')).toThrow(InvalidResidencyError);
    expect(() => normalizeResidency(null)).toThrow(InvalidResidencyError);
  });
});
