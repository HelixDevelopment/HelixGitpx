export default {
  displayName: 'web',
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(ts|mjs|js|html)$': [
      'ts-jest',
      { tsconfig: '<rootDir>/tsconfig.spec.json', isolatedModules: true },
    ],
  },
  moduleFileExtensions: ['ts', 'mjs', 'js'],
  transformIgnorePatterns: ['node_modules/(?!(@angular|@bufbuild|@connectrpc|rxjs|zone\\.js|tslib)/)'],
  coverageDirectory: '../../coverage/apps/web',
  collectCoverage: true,
  coverageReporters: ['text', 'lcov', 'html', 'json-summary'],
  coverageThreshold: {
    global: { branches: 0, functions: 0, lines: 0, statements: 0 },
  },
  testMatch: ['<rootDir>/src/**/*.spec.ts'],
};
