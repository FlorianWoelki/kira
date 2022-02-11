const app = require('./app');

test('sum test', () => {
  expect(app.sum(1, 2)).toBe(3);
});
