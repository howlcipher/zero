const test = require("node:test");
const assert = require("node:assert");
const { add } = require("./main.js");

test("add function returns correct sum", () => {
  assert.strictEqual(add(2, 3), 5);
});
