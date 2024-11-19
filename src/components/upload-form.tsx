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
        window.location.reload();
      }}
      class="space-y-4"
    >
      <FileInput />
      <Button class="px-8" type="submit">
        Save
      </Button>
    </form>
  );
};

export default UploadForm;
