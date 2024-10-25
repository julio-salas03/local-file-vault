import type { Component, ComponentProps } from 'solid-js';
import { splitProps } from 'solid-js';
import { cn } from '@/lib/utils';
import { Label } from './label';

const FileInput: Component = () => {
  return (
    <Label class="relative block border-4 border-dashed border-border py-5 text-center">
      Drag your file here! <br />
      or click to select
      <input
        type="file"
        name="file"
        class="absolute inset-0 cursor-pointer opacity-0"
      />
    </Label>
  );
};

export { FileInput };
