import { type Component } from 'solid-js';
import { Toaster } from 'solid-sonner';
import UploadForm from '@/components/upload-form';
import { Separator } from './components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import Account from './components/account';
import { AuthProvider } from './lib/auth';
import Files from './components/files';

const App: Component = () => {
  return (
    <AuthProvider>
      <div class="font-inter antialiased">
        <main class="px-5 py-8">
          <Tabs defaultValue="account">
            <TabsList class="grid w-full grid-cols-3">
              <TabsTrigger value="upload">Upload</TabsTrigger>
              <TabsTrigger value="files">Files</TabsTrigger>
              <TabsTrigger value="account">Account</TabsTrigger>
            </TabsList>
            <Separator class="my-5" />
            <TabsContent value="upload">
              <UploadForm />
            </TabsContent>
            <TabsContent value="files">
              <Files />
            </TabsContent>
            <TabsContent value="account">
              <Account />
            </TabsContent>
          </Tabs>
        </main>
        <Toaster />
      </div>
    </AuthProvider>
  );
};

export default App;
