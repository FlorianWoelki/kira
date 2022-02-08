def custom_sum(a, b):
	return a + b

def main():
  a = 42

  for i in range(0, 42):
    a += i

  print('Hello World', a);
  print("sum of 1 + 2 is", custom_sum(1, 2))

if __name__ == "__main__":
  main()
