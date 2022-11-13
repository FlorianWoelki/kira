import { useRef, useState } from 'react';
import { MonacoEditor } from './MonacoEditor';

interface CodeExecutionResult {
  compileOutput: string;
  compileError: string;
  compileTime: number;
  runOutput: string;
  runError: string;
  runTime: number;
}

const App: React.FC = (): JSX.Element => {
  const codeEditorRef = useRef<any>(null);
  const [codeResult, setCodeResult] = useState<CodeExecutionResult>();
  const [bypassCache, setBypassCache] = useState<boolean>(false);

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
          content: codeEditorRef.current.getValue(),
        }),
      },
    );

    const jsonResult: CodeExecutionResult = await result.json();
    setCodeResult(jsonResult);
    setIsLoading(false);
  };

  const normalizeOutput = (value: string): string[] => {
    return value.split('\n');
  };

  return (
    <div className="relative antialiased h-screen">
      <div className="flex flex-col h-full">
        <div className="p-2 flex items-center justify-center">
          <button
            className="flex items-center px-4 py-2 bg-green-400 rounded-lg text-green-800 font-semibold transition duration-100 ease-in-out hover:bg-green-500 disabled:cursor-not-allowed disabled:bg-opacity-50"
            onClick={runCode}
            disabled={isLoading}
          >
            Run
          </button>
        </div>

        <div
          className="grid grid-cols-2 bg-gray-200 gap-2 p-2"
          style={{ height: 'calc(100% - 56px)' }}
        >
          <div className="rounded-lg bg-white">
            <MonacoEditor
              value="print('Hello World')"
              onCtrlCmdEnter={runCode}
              ref={codeEditorRef}
            ></MonacoEditor>
          </div>
          <div className="rounded-lg bg-white p-4 overflow-auto h-full">
            <p className="font-semibold">Output:</p>
            {codeResult ? (
              codeResult.compileError || codeResult.runError ? (
                <>
                  <p>Compile Error: {codeResult.compileError}</p>
                  <p>Run Error: {codeResult.runError}</p>
                </>
              ) : (
                <>
                  <p className="italic text-sm mb-4">
                    Time: {codeResult.runTime}ms
                  </p>
                  {normalizeOutput(codeResult.runOutput).map((v, i) => (
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
