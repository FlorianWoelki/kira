import { useState } from 'react';
import { Checkbox } from './Checkbox';
import { CodeMirrorEditor } from './CodeMirrorEditor';

const useCodeMirrorEditor = () => {
  const [code, setCode] = useState<string>('');
  const [editorOptions, setEditorOptions] = useState<{
    line: number;
    column: number;
  }>();

  return { code, setCode, editorOptions, setEditorOptions };
};

interface Output {
  result: string;
  error: string;
  time: number;
}

interface CodeExecutionResult {
  compileOutput: Output;
  runOutput: Output;
  testOutput: {
    results: {
      name: string;
      received: string;
      actual: string;
      passed: boolean;
    }[];
    time: number;
  };
}

const App: React.FC = (): JSX.Element => {
  const [codeResult, setCodeResult] = useState<CodeExecutionResult>();
  const [bypassCache, setBypassCache] = useState<boolean>(false);

  const codeEditor = useCodeMirrorEditor();
  const testEditor = useCodeMirrorEditor();

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
          content: codeEditor.code,
          tests: JSON.parse(testEditor.code),
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

          <div className="flex flex-col justify-between col-span-2 gap-2">
            <div className="flex flex-col rounded-lg bg-white overflow-auto h-full">
              <div className="overflow-auto flex-1">
                <CodeMirrorEditor
                  language="python"
                  defaultValue={`print("Hello World")

def custom_multiply(a, b):
  return a * b

def custom_sum(a, b):
  return a + b`}
                  onChange={(v, options) => {
                    codeEditor.setCode(v);
                    codeEditor.setEditorOptions(options);
                  }}
                ></CodeMirrorEditor>
              </div>
              {codeEditor.editorOptions && (
                <div className="border-t p-2 text-sm text-gray-600">
                  Line: {codeEditor.editorOptions.line} Column:{' '}
                  {codeEditor.editorOptions.column}
                </div>
              )}
            </div>
            <div className="flex flex-col rounded-lg bg-white overflow-auto h-full">
              <div className="overflow-auto flex-1">
                <CodeMirrorEditor
                  language="json"
                  defaultValue={`[
  { "name": "Test 1", "actual": "2" },
  { "name": "Test 2", "actual": "2" }
]`}
                  onChange={(v, options) => {
                    testEditor.setCode(v);
                    testEditor.setEditorOptions(options);
                  }}
                ></CodeMirrorEditor>
              </div>
              {testEditor.editorOptions && (
                <div className="border-t p-2 text-sm text-gray-600">
                  Line: {testEditor.editorOptions.line} Column:{' '}
                  {testEditor.editorOptions.column}
                </div>
              )}
            </div>
          </div>

          <div className="rounded-lg bg-white p-4 overflow-auto h-full col-span-2">
            <p className="font-semibold">Program Output:</p>
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

            <div className="border my-6 border-gray-100"></div>

            <p className="font-semibold">Test Output:</p>
            {codeResult ? (
              <>
                <p className="italic text-sm mb-4">
                  Time: {codeResult.testOutput.time}ms
                </p>
                {codeResult.testOutput.results.map((result) => (
                  <p key={result.name} className="flex flex-col mb-4">
                    <span>
                      {result.passed ? '🟩' : '🟥'} Name: {result.name}
                    </span>
                    {!result.passed && (
                      <>
                        <span>Actual value: {result.actual}</span>
                        <span>Received value: {result.received}</span>
                      </>
                    )}
                  </p>
                ))}
              </>
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
