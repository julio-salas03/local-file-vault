import { Component, createEffect, createSignal, For } from 'solid-js';
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from './ui/drawer';
import Trash from './icons/trash';
import { Separator } from './ui/separator';
import { filesizeToHumanReadable, ServerResponseSchema } from '@/lib/utils';
import dayjs from 'dayjs';
import { Button } from './ui/button';
import Dots from './icons/dots';
import CloudDownload from './icons/cloud-download';
import { z } from 'zod';
import { APIError } from '@/lib/errors';

const FILE_LIST_SCHEMA = z.array(
  z.object({
    size: z.number(),
    name: z.string(),
    lastmod: z.string(),
    download: z.string(),
    owner: z.string(),
  })
);

const Files: Component = () => {
  const [fileList, setFileList] = createSignal<
    z.infer<typeof FILE_LIST_SCHEMA>
  >([]);

  createEffect(async () => {
    try {
      const res = await fetch('/api/files');
      const data = await res.json();
      const schema = ServerResponseSchema(
        z.object({ files: FILE_LIST_SCHEMA })
      );
      const parse = schema.parse(data);

      if (parse.type === 'error') {
        throw new APIError(parse.message, parse.errorCode);
      }

      setFileList(parse.data.files);
    } catch (error) {
      console.error(error);
    }
  });

  return (
    <ul class="grid grid-cols-1 gap-px border-y border-border bg-border">
      <For each={fileList()} fallback={<p>list of files here...</p>}>
        {item => (
          <li class="flex items-center justify-between gap-2 bg-background py-3">
            <div>
              <p class="line-clamp-1 text-lg">{item.name}</p>
              <div class="text-base opacity-80">
                <span>{filesizeToHumanReadable(item.size)}</span>
                {' - '}
                <span>{dayjs(item.lastmod).format('MMM DD, YYYY, hh:mm')}</span>
              </div>
            </div>
            <Drawer>
              <DrawerTrigger as={Button<'button'>} class="p-2" variant="ghost">
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
                        <span>Owner: {item.owner}</span>
                        {' - '}
                        <span>{filesizeToHumanReadable(item.size)}</span>
                        <br />
                        <span>
                          {dayjs(item.lastmod).format('MMM DD, YYYY, hh:mm')}
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
  );
};

export default Files;
