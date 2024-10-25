import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import { cn } from "@/lib/utils"
import { Label } from "./label"

const FileInput: Component = () => {
  return (
    <Label class="block text-center border-border border-dashed border-4 py-5 relative">
        Drag your file here! <br /> 
        or click to select
        <input type="file" name="file" class="absolute inset-0 opacity-0 cursor-pointer"/>
    </Label>
  )
}

export { FileInput }
