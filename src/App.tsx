import { createSignal, For, type Component } from 'solid-js';
import { Toaster } from 'solid-sonner';
import Cog from '@/components/icons/cog';
import UploadForm from '@/components/upload-form';
import { effect } from 'solid-js/web';
import { z } from 'zod';
import { filesizeToHumanReadable } from './lib/utils';
import dayjs from 'dayjs';
import CloudDownload from './components/icons/cloud-download';
import { buttonVariants } from './components/ui/button';

const FILE_LIST_SCHEMA = z.array(
  z.object({
    size: z.number(),
    name: z.string(),
    lastmod: z.string(),
    download: z.string(),
  })
);

const App: Component = () => {
  const [fileList, setFileList] = createSignal<
    z.infer<typeof FILE_LIST_SCHEMA>
  >([]);

  effect(async () => {
    const res = await fetch('/api/files');
    const data = await res.json();
    const parse = FILE_LIST_SCHEMA.safeParse(data);
    if (!parse.success) console.error(parse.error);
    else setFileList(parse.data);
  });

  return (
    <div class="font-inter antialiased">
      <nav class="bg-foreground px-5 py-4 text-background">
        <Cog class="ml-auto" />
      </nav>

      <main class="px-5 py-8">
        <ul class="grid grid-cols-1 gap-px border-y border-border bg-border">
          <For each={fileList()} fallback={<p>list of files here...</p>}>
            {item => (
              <li class="flex items-center justify-between gap-2 bg-background py-3">
                <div>
                  <p class="line-clamp-1 text-lg">{item.name}</p>
                  <div class="text-base opacity-80">
                    <span>{filesizeToHumanReadable(item.size)}</span>
                    {' - '}
                    <span>
                      {dayjs(item.lastmod).format('MMM DD, YYYY, hh:mm')}
                    </span>
                  </div>
                </div>
                <a
                  href={item.download}
                  download={item.name}
                  class={buttonVariants({
                    variant: 'ghost',
                    class: 'flex-shrink-0 !p-2',
                  })}
                >
                  <span class="sr-only">download file</span>
                  <CloudDownload width="2em" height="2em" />
                </a>
              </li>
            )}
          </For>
        </ul>

        <UploadForm />
      </main>
      <Toaster />
    </div>
  );
};

export default App;
