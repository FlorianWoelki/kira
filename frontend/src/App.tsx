import { useRef } from 'react';
import { MonacoEditor } from './Editor';

const App: React.FC = (): JSX.Element => {
  const codeEditorRef = useRef<any>(null);

  const runCode = async (): Promise<void> => {
    const result = await fetch('http://localhost:9090/execute', {
      method: 'POST',
      body: JSON.stringify({
        language: 'python',
        content: codeEditorRef.current.getValue(),
      }),
    });

    const jsonResult = await result.json();
    console.log(jsonResult);
  };

  return (
    <div className="relative antialiased h-screen">
      <MonacoEditor
        value="print('Hello World')"
        ref={codeEditorRef}
      ></MonacoEditor>

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
