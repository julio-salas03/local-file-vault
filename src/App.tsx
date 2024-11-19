import { createSignal, For, type Component } from 'solid-js';
import { Toaster } from 'solid-sonner';
import Cog from '@/components/icons/cog';
import UploadForm from '@/components/upload-form';
import { effect } from 'solid-js/web';
import { z } from 'zod';
import { filesizeToHumanReadable } from './lib/utils';
import dayjs from 'dayjs';
import { Button } from '@/components/ui/button';
import Dots from '@/components/icons/dots';
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';
import CloudDownload from './components/icons/cloud-download';
import { Separator } from './components/ui/separator';
import Trash from './components/icons/trash';

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
                <Drawer>
                  <DrawerTrigger
                    as={Button<'button'>}
                    class="p-2"
                    variant="ghost"
                  >
                    <span class="sr-only">download file</span>
                    <Dots width="2em" height="2em" />
                  </DrawerTrigger>
                  <DrawerContent class="px-5">
                    <div class="mx-auto w-full max-w-sm pb-5">
                      <DrawerHeader class="px-0 text-left">
                        <DrawerTitle class="line-clamp-1 !text-lg">
                          {item.name}
                        </DrawerTitle>
                        <DrawerDescription>
                          <div class="text-base opacity-80">
                            <span>{filesizeToHumanReadable(item.size)}</span>
                            {' - '}
                            <span>
                              {dayjs(item.lastmod).format(
                                'MMM DD, YYYY, hh:mm'
                              )}
                            </span>
                          </div>
                        </DrawerDescription>
                      </DrawerHeader>
                      <Separator />
                      <ul class="mt-5 grid grid-cols-1 gap-2 text-lg">
                        <li class="bg-background">
                          <DrawerClose
                            as="a"
                            href={item.download}
                            download={item.name}
                            class="grid grid-cols-12 font-semibold"
                          >
                            <CloudDownload
                              class="col-span-2 mx-auto"
                              width="2rem"
                              height="2rem"
                            />
                            <span>Download</span>
                          </DrawerClose>
                        </li>
                        <li class="bg-background text-destructive">
                          <DrawerClose
                            as="button"
                            class="grid grid-cols-12 font-semibold opacity-70"
                            disabled
                          >
                            <Trash
                              class="col-span-2 mx-auto"
                              width="1.5rem"
                              height="1.5rem"
                            />
                            Delete
                          </DrawerClose>
                        </li>
                      </ul>
                    </div>
                  </DrawerContent>
                </Drawer>
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
