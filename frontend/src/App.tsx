import { useEffect, useState } from 'react';
import { Checkbox } from './Checkbox';
import { CodeMirrorEditor } from './CodeMirrorEditor';
import { Dropdown } from './Dropdown';
import { CodeTemplate, codeTemplates } from './codeTemplates';

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
  const [bypassCache, setBypassCache] = useState<boolean>(true);
  const [stdin, setStdin] = useState<string>('');
  const [template, setTemplate] = useState<CodeTemplate>(codeTemplates[0]);

  const codeEditor = useCodeMirrorEditor();
  const testEditor = useCodeMirrorEditor();

  const [isLoading, setIsLoading] = useState<boolean>(false);

  useEffect(() => {
    setStdin(template.defaultStdin);
  }, [template.defaultStdin]);

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
          stdin: [stdin],
          tests:
            testEditor.code.length === 0 ? [] : JSON.parse(testEditor.code),
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
        <div className="flex items-center justify-center p-2">
          <div className="flex items-center justify-center flex-1 space-x-4">
            <Checkbox
              id="bypass-cache"
              checked={bypassCache}
              onChange={() => setBypassCache((v) => !v)}
            >
              Bypass cache?
            </Checkbox>
            <button
              className="flex items-center px-4 py-2 font-semibold text-green-800 transition duration-100 ease-in-out bg-green-400 rounded-lg hover:bg-green-500 disabled:cursor-not-allowed disabled:bg-opacity-50"
              onClick={runCode}
              disabled={isLoading}
            >
              Run
            </button>
            <input
              type="text"
              className="block p-2 border border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              placeholder="Your input"
              onChange={(e) => setStdin(e.target.value)}
              value={stdin}
            />
          </div>

          <Dropdown
            title="Templates"
            items={codeTemplates.map((template) => template.name)}
            onSelect={(name) =>
              setTemplate(
                codeTemplates.find((template) => template.name === name)!,
              )
            }
          ></Dropdown>
        </div>

        <div
          className="grid grid-cols-5 gap-2 p-2 bg-gray-200"
          style={{ height: 'calc(100% - 56px)' }}
        >
          <div className="p-2 space-y-4 overflow-auto bg-white rounded-lg">
            <h2 className="text-xl font-bold">Files</h2>

            <ul className="space-y-1">
              <li
                aria-selected="true"
                className="px-2 py-1 transition duration-100 ease-in-out bg-gray-200 rounded-lg hover:bg-gray-200"
              >
                <span>main.py</span>
              </li>
            </ul>
          </div>

          <div className="flex flex-col justify-between col-span-2 gap-2">
            <div className="flex flex-col h-full overflow-auto bg-white rounded-lg">
              <div className="flex-1 overflow-auto">
                <CodeMirrorEditor
                  language="python"
                  defaultValue={template.code}
                  onChange={(v, options) => {
                    codeEditor.setCode(v);
                    codeEditor.setEditorOptions(options);
                  }}
                ></CodeMirrorEditor>
              </div>
              {codeEditor.editorOptions && (
                <div className="p-2 text-sm text-gray-600 border-t">
                  Line: {codeEditor.editorOptions.line} Column:{' '}
                  {codeEditor.editorOptions.column}
                </div>
              )}
            </div>
            <div className="flex flex-col h-full overflow-auto bg-white rounded-lg">
              <div className="flex-1 overflow-auto">
                <CodeMirrorEditor
                  language="json"
                  defaultValue={template.testCode}
                  onChange={(v, options) => {
                    testEditor.setCode(v);
                    testEditor.setEditorOptions(options);
                  }}
                ></CodeMirrorEditor>
              </div>
              {testEditor.editorOptions && (
                <div className="p-2 text-sm text-gray-600 border-t">
                  Line: {testEditor.editorOptions.line} Column:{' '}
                  {testEditor.editorOptions.column}
                </div>
              )}
            </div>
          </div>

          <div className="h-full col-span-2 p-4 overflow-auto bg-white rounded-lg">
            <p className="font-semibold">Program Output:</p>
            {codeResult ? (
              codeResult.compileOutput.error || codeResult.runOutput.error ? (
                <>
                  <p>Compile Error: {codeResult.compileOutput.error}</p>
                  <p>Run Error: {codeResult.runOutput.error}</p>
                </>
              ) : (
                <>
                  <p className="mb-4 text-sm italic">
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

            <div className="my-6 border border-gray-100"></div>

            {codeResult?.testOutput.results?.length ? (
              <>
                <p className="font-semibold">Test Output:</p>
                <p className="mb-4 text-sm italic">
                  Time: {codeResult.testOutput.time}ms
                </p>
                {codeResult.testOutput.results.map((result) => (
                  <p key={result.name} className="flex flex-col mb-4">
                    <span>
                      {result.passed ? '🟩' : '🟥'} Name: {result.name}
                    </span>
                    <span>Actual value: {result.actual}</span>
                    <span>Received value: {result.received}</span>
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
