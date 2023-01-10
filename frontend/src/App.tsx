import { useState } from 'react';
import { Checkbox } from './Checkbox';
import { CodeMirrorEditor } from './CodeMirrorEditor';

interface Output {
  result: string;
  error: string;
  time: number;
}

interface CodeExecutionResult {
  compileOutput: Output;
  runOutput: Output;
  testOutput: { output: string; time: number };
}

const App: React.FC = (): JSX.Element => {
  const [codeResult, setCodeResult] = useState<CodeExecutionResult>();
  const [bypassCache, setBypassCache] = useState<boolean>(false);
  const [code, setCode] = useState<string>('');
  const [editorOptions, setEditorOptions] = useState<{
    line: number;
    column: number;
  }>();

  const [isLoading, setIsLoading] = useState<boolean>(false);

  const runCode = async (): Promise<void> => {
    if (isLoading) {
      return;
    }

    setIsLoading(true);
    const result = await fetch(
      `http://localhost:9090/execute${bypassCache ? '?bypass_cache=true' : ''}`,
      {
        method: 'POST',
        body: JSON.stringify({
          language: 'python',
          // content: code,
          content: `def custom_sum(a, b):
  return a + b`,
          test: `import unittest
from .app import custom_sum

class TestSumFunction(unittest.TestCase):
  def test_sum_2(self):
    self.assertEqual(custom_sum(2, 4), 4)
    self.assertEqual(custom_sum(20, 20), 40)
    self.assertEqual(custom_sum(200, 200), 400)

  def test_sum(self):
    self.assertEqual(custom_sum(1, 1), 2)
    self.assertEqual(custom_sum(10, 10), 20)
    self.assertEqual(custom_sum(100, 100), 200)
    self.assertEqual(custom_sum(1000, 1000), 2000)`,
        }),
      },
    );

    const jsonResult: CodeExecutionResult = await result.json();
    console.log(jsonResult);
    setCodeResult(jsonResult);
    setIsLoading(false);
  };

  const normalizeOutput = (value: string): string[] => {
    return value.split('\n');
  };

  return (
    <div className="relative h-screen">
      <div className="flex flex-col h-full">
        <div className="p-2 flex items-center justify-center space-x-4">
          <Checkbox
            id="bypass-cache"
            onChange={() => setBypassCache((v) => !v)}
          >
            Bypass cache?
          </Checkbox>
          <button
            className="flex items-center px-4 py-2 bg-green-400 rounded-lg text-green-800 font-semibold transition duration-100 ease-in-out hover:bg-green-500 disabled:cursor-not-allowed disabled:bg-opacity-50"
            onClick={runCode}
            disabled={isLoading}
          >
            Run
          </button>
        </div>

        <div
          className="grid grid-cols-5 bg-gray-200 gap-2 p-2"
          style={{ height: 'calc(100% - 56px)' }}
        >
          <div className="rounded-lg bg-white overflow-auto p-2 space-y-4">
            <h2 className="font-bold text-xl">Files</h2>

            <ul className="space-y-1">
              <li
                aria-selected="true"
                className="rounded-lg bg-gray-200 px-2 py-1 hover:bg-gray-200 transition duration-100 ease-in-out"
              >
                <span>main.py</span>
              </li>
            </ul>
          </div>

          <div className="flex flex-col rounded-lg bg-white col-span-2 overflow-auto">
            <div className="overflow-auto flex-1">
              <CodeMirrorEditor
                onChange={(v, options) => {
                  setCode(v);
                  setEditorOptions(options);
                }}
              ></CodeMirrorEditor>
            </div>
            {editorOptions && (
              <div className="border-t p-2 text-sm text-gray-600">
                Line: {editorOptions.line} Column: {editorOptions.column}
              </div>
            )}
          </div>

          <div className="rounded-lg bg-white p-4 overflow-auto h-full col-span-2">
            <p className="font-semibold">Output:</p>
            {codeResult ? (
              codeResult.compileOutput.error || codeResult.runOutput.error ? (
                <>
                  <p>Compile Error: {codeResult.compileOutput.error}</p>
                  <p>Run Error: {codeResult.runOutput.error}</p>
                </>
              ) : (
                <>
                  <p className="italic text-sm mb-4">
                    Time: {codeResult.runOutput.time}ms
                  </p>
                  {normalizeOutput(codeResult.runOutput.result).map((v, i) => (
                    <p key={i}>{v}</p>
                  ))}
                </>
              )
            ) : (
              ''
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;
