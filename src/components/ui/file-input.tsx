import type { Component, ComponentProps } from 'solid-js';
import { splitProps } from 'solid-js';
import { cn } from '@/lib/utils';
import { Label } from './label';

const FileInput: Component = () => {
  return (
    <div class="relative">
      <input
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
  );
};

export { FileInput };
