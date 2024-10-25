import type { Component } from 'solid-js';
import { Button } from '@/components/ui/button';
import { effect } from 'solid-js/web';

const App: Component = () => {
  effect(()=> {
    const ping =async ()=> {
      const response = await fetch("/api/ping")
      const text = await response.text()
      console.log(text)
    }
    ping()
  })
  return (
    <div>
      <Button>click me</Button>
    </div>
  );
};

export default App;
