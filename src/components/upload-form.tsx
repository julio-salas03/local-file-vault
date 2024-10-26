import { Component, createSignal } from 'solid-js';
import { FileInput } from './ui/file-input';
import { Button, buttonVariants } from './ui/button';
import { toast } from 'solid-sonner';
import { Dialog, DialogContent, DialogTitle, DialogTrigger } from './ui/dialog';

const UploadForm: Component = () => {
  const [isOpen, setIsOpen] = createSignal(false);
  return (
    <Dialog open={isOpen()} onOpenChange={open => setIsOpen(open)}>
      <DialogTrigger
        class={buttonVariants({
          class: 'fixed bottom-5 right-5 items-center gap-1',
        })}
      >
        Backup a file <span class="text-[1.5em]">+</span>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle class="text-2xl">Backup a file</DialogTitle>
        <form
          onSubmit={async e => {
            e.preventDefault();
            const response = await fetch('/api/upload', {
              method: 'POST',
              body: new FormData(e.currentTarget),
            });
            const text = await response.text();
            toast(text);
            setIsOpen(false);
          }}
          class="space-y-4 py-5"
        >
          <FileInput />
          <Button class="px-8" type="submit">
            Save
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default UploadForm;
