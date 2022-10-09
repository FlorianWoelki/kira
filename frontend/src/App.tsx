import { useRef, useState } from 'react';
import { MonacoEditor } from './Editor';

interface CodeExecutionResult {
  output: string;
  error: string;
}

const App: React.FC = (): JSX.Element => {
  const codeEditorRef = useRef<any>(null);
  const [codeResult, setCodeResult] = useState<string>('');

  const [isLoading, setIsLoading] = useState<boolean>(false);

  const runCode = async (): Promise<void> => {
    if (isLoading) {
      return;
    }

    setIsLoading(true);
    const result = await fetch('http://localhost:9090/execute', {
      method: 'POST',
      body: JSON.stringify({
        language: 'python',
        content: codeEditorRef.current.getValue(),
      }),
    });

    const jsonResult: CodeExecutionResult = await result.json();
    if (jsonResult.error) {
      setCodeResult(`Error: ${jsonResult.error}`);
    } else {
      setCodeResult(`${jsonResult.output}`);
    }

    console.log(jsonResult);
    setIsLoading(false);
  };

  return (
    <div className="relative antialiased h-screen">
      <MonacoEditor
        value="print('Hello World')"
        ref={codeEditorRef}
      ></MonacoEditor>
      <div className="absolute left-0 bottom-0 px-4 py-2 mb-4 ml-4 bg-gray-200 rounded">
        <p>Output: {codeResult.length === 0 ? '/' : codeResult}</p>
      </div>
      <button
        className="flex items-center absolute right-0 bottom-0 px-4 py-2 bg-green-600 rounded text-white mb-4 mr-4 transition duration-100 ease-in-out hover:bg-green-500 disabled:cursor-not-allowed disabled:bg-opacity-50"
        onClick={runCode}
        disabled={isLoading}
      >
        {isLoading ? (
          <>
            <svg
              className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              ></circle>
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <span>Running...</span>
          </>
        ) : (
          <span>Run</span>
        )}
      </button>
    </div>
  );
};

export default App;
