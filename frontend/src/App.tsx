import { useRef } from 'react';
import { MonacoEditor } from './Editor';

const App: React.FC = (): JSX.Element => {
  const codeEditorRef = useRef<any>(null);

  return (
    <div className="antialiased h-screen">
      <MonacoEditor
        value="print('Hello World')"
        ref={codeEditorRef}
      ></MonacoEditor>
    </div>
  );
};

export default App;
