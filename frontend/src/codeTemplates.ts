export interface CodeTemplate {
  name: string;
  code: string;
  testCode: string;
}

export const codeTemplates = [
  {
    name: 'Without Tests',
    code: `print(2)

def custom_multiply(a, b):
  return a * b

def custom_sum(a, b):
  return a + b`,
    testCode: '',
  },
  {
    name: 'With Tests',
    code: `import sys

value = sys.argv[1]
print(value + 1)`,
    testCode: `[
  { "name": "Test 1", "stdin": "2", "actual": "3" },
  { "name": "Test 2", "stdin": "3", "actual": "4" }
]`,
  },
];
