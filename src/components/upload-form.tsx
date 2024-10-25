import { Component } from 'solid-js';
import { FileInput } from './ui/file-input';
import { Button } from './ui/button';
import { toast } from 'solid-sonner';

const UploadForm: Component = () => {
  return (
    <form
      onSubmit={async e => {
        e.preventDefault();
        const response = await fetch('/api/upload', {
          method: 'POST',
          body: new FormData(e.currentTarget),
        });
        const text = await response.text();
        toast(text);
      }}
      class="mx-auto max-w-2xl space-y-4 py-5"
    >
      <h1 class="text-xl">Backup a file</h1>
      <FileInput />
      <Button class="px-8" type="submit">
        Save
      </Button>
    </form>
  );
};

export default UploadForm;
