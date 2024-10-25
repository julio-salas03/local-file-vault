import type { Component, ComponentProps } from 'solid-js';
import { createSignal, splitProps } from 'solid-js';
import { cn, getFileName } from '@/lib/utils';
import { Label } from './label';

const FileInput: Component = () => {
  const [filename, setFileName] = createSignal('');
  return (
    <div>
      <div class="relative">
        <input
          onChange={e => setFileName(getFileName(e.currentTarget.value))}
          required
          id="file"
          type="file"
          name="file"
          class="peer absolute inset-0 z-10 cursor-pointer opacity-0"
        />
        <Label
          for="id"
          class="relative block border-4 border-dashed border-border py-10 text-center text-lg peer-focus:outline peer-focus:outline-2 peer-focus:outline-offset-2"
        >
          Drag your file here! <br />
          or click to select
        </Label>
      </div>
      <span class="mt-2 block">
        {filename().length ? (
          <span>
            Selected: <strong>{filename()}</strong>
          </span>
        ) : (
          'No file selected'
        )}
      </span>
    </div>
  );
};

export { FileInput };
