import { useRef, useState } from 'react';
import { MonacoEditor } from './Editor';

interface CodeExecutionResult {
  output: string;
  error: string;
}

const App: React.FC = (): JSX.Element => {
  const codeEditorRef = useRef<any>(null);
  const [codeResult, setCodeResult] = useState<string>('');

  const runCode = async (): Promise<void> => {
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
        className="absolute right-0 bottom-0 px-4 py-2 bg-green-600 rounded text-white mb-4 mr-4 transition duration-100 ease-in-out hover:bg-green-500"
        onClick={runCode}
      >
        Run
      </button>
    </div>
  );
};

export default App;
