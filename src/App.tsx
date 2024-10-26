import type { Component } from 'solid-js';
import { Toaster } from 'solid-sonner';
import Cog from './components/icons/cog';
import UploadForm from './components/upload-form';

const App: Component = () => {
  return (
    <div class="font-inter antialiased">
      <nav class="bg-foreground px-5 py-4 text-background">
        <Cog class="ml-auto" />
      </nav>
      <main class="px-5 py-8">
        <p>list of files here...</p>
        <UploadForm />
      </main>
      <Toaster />
    </div>
  );
};

export default App;
