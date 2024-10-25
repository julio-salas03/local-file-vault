import type { Component } from 'solid-js';
import { Button } from '@/components/ui/button';
import { FileInput } from './components/ui/file-input';

const App: Component = () => {
  return (
    <div>
      <form class="mx-auto max-w-2xl space-y-4 py-5">
        <FileInput />
        <Button type="submit">save</Button>
      </form>
    </div>
  );
};

export default App;
