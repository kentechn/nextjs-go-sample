import path from "node:path";
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Docker では pnpm workspace ルートを含めて標準出力(standalone)をトレースする。
  output: "standalone",
  outputFileTracingRoot: path.join(import.meta.dirname, "../.."),
  // next の require-hook が参照する @swc/helpers の ESM 実装はトレースから漏れるため明示的に含める。
  outputFileTracingIncludes: {
    "**": ["../../node_modules/.pnpm/@swc+helpers@*/node_modules/@swc/helpers/esm/**"],
  },
};

export default nextConfig;
