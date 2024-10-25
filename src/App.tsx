import type { Component } from 'solid-js';
import { Button } from '@/components/ui/button';
import { FileInput } from './components/ui/file-input';

const App: Component = () => {
  return (
    <div>
      <form class='max-w-2xl mx-auto py-5 space-y-4'>
      <FileInput/>
      <Button type='submit'>save</Button>
      </form>
    </div>
  );
};

export default App;
