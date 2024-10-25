import type { Component } from 'solid-js';
import UploadForm from './components/upload-form';
import { Toaster } from 'solid-sonner';

const App: Component = () => {
  return (
    <main class="font-inter px-5 py-8 antialiased">
      <UploadForm />
      <Toaster />
    </main>
  );
};

export default App;
