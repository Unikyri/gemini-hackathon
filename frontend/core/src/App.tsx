import { PathGenerator } from '@/features/path-generator';
import './App.css';

function App() {
  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8">
          Learning Path Generator
        </h1>
        <PathGenerator />
      </div>
    </div>
  );
}

export default App;
