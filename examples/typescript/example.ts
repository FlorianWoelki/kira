let a: number = 42;

for (let i: number = 0; i < 42; i++) {
  a += i;
}

console.log('Hello World' as string, a as number);

export const sum = (a: number, b: number): number => {
  return a + b;
};

console.log('1 + 2 =', sum(1, 2));
