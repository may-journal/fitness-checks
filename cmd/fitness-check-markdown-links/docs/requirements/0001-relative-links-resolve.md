# 0001 Relative Links Resolve

## Why

I follow every relative link in my docs to a file that exists.

## Measurement

relative links that resolve

-

relative links in Markdown files

## Requirements

1. Broken relative links fail naming file and line
   1. GIVEN Developer has a repo with Markdown files
      1. WHEN a link targets a missing file
         1. THEN it fails as `path:line: broken relative link: target`
      2. WHEN a link in a subfolder targets a missing sibling
         1. THEN it fails resolving the target from that subfolder
      3. WHEN an image targets a missing file
         1. THEN it fails naming the image
         2. THEN a reference definition fails alike
      4. WHEN a `.markdown` file holds a broken link
         1. THEN it fails
         2. THEN an `.mdx` file fails alike
2. Links to what exists pass
   1. GIVEN Developer has a repo with Markdown files
      1. WHEN a link targets a file
         1. THEN it passes
         2. THEN a folder passes
         3. THEN a heading passes since a `#fragment` is stripped and not validated
      2. WHEN a link is a web address
         1. THEN it passes untouched
         2. THEN a mail link passes alike
         3. THEN a protocol-relative link passes alike
         4. THEN a same-page link passes alike
      3. WHEN a broken link sits in inline code or a fenced block
         1. THEN it is ignored
      4. WHEN a link holds a percent-encoded space
         1. THEN it is decoded and resolves
