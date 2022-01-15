package file

import "testing"

func TestExtractCodeOfFile(t *testing.T) {
	filepath := "../examples/python/example.py"

	expectedCode := `a = 42

for i in range(0, 42):
  a += i

print('Hello World', a);
`

	code, err := ExtractCodeOfFile(filepath)
	if err != nil {
		t.Errorf("something went wrong extracting the code of the file with path %s = %s", filepath, err)
	}

	if code != expectedCode {
		t.Fatalf("extracted code is not the same as expected code. expected=%s, got=%s", expectedCode, code)
	}
}
