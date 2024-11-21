import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { z } from 'zod';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const getFileName = (path: string) => path.split(/(\\|\/)+/).at(-1);

/**
 * Converts file size in bytes to a human-readable format.
 * @param size - The size of the file in bytes.
 * @param decimals - Number of decimal places to include (default is 2).
 * @returns A string representing the size in a human-readable format.
 */
export const filesizeToHumanReadable = (
  size: number,
  decimals: number = 2
): string => {
  const BINARY_BASE = 1024;

  const UNITS = [
    { unit: 'GB', scale: BINARY_BASE ** 3 },
    { unit: 'MB', scale: BINARY_BASE ** 2 },
    { unit: 'KB', scale: BINARY_BASE },
  ];

  for (const { unit, scale } of UNITS) {
    if (size >= scale) {
      return (size / scale).toFixed(decimals) + unit;
    }
  }

  return size + 'B';
};

export function ServerResponseSchema<T extends z.ZodType>(data: T) {
  const successfulResponse = z.object({
    data,
    message: z.string(),
    type: z.literal('success'),
  });

  const errorResponse = z.object({
    errorCode: z.string(),
    message: z.string(),
    type: z.literal('error'),
  });

  return successfulResponse.or(errorResponse);
}
