import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { z } from 'zod';
import { API_ERROR_CODES } from './errors';

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

/**
 * Allows to generate a `z.union` from dynamic `z.literal`'s.
 *
 * Code taken from https://github.com/colinhacks/zod/discussions/2790#discussioncomment-7096060
 */

export function unionOfLiterals<T extends string | number>(
  constants: readonly T[]
) {
  const literals = constants.map(x => z.literal(x)) as unknown as readonly [
    z.ZodLiteral<T>,
    z.ZodLiteral<T>,
    ...z.ZodLiteral<T>[],
  ];
  return z.union(literals);
}

export function ServerResponseSchema<T extends z.ZodType>(data: T) {
  const successfulResponse = z.object({
    data,
    message: z.string(),
    type: z.literal('success'),
  });

  const errorResponse = z.object({
    errorCode: unionOfLiterals(Object.values(API_ERROR_CODES)),
    message: z.string(),
    type: z.literal('error'),
  });

  return successfulResponse.or(errorResponse);
}
