# 0001 Containers Match What Ships

## Why

I get my containers view matched against the runtime products my repo builds and deploys.

## Measurement

repos whose containers view matches their products

-

repos with a containers view

## Requirements

1. Every product has one container
   1. GIVEN Developer has a containers view
   2. GIVEN it has products proven by build targets or entry points or deploy files
      1. WHEN each product has one container naming its folder
         1. THEN it passes counting the container views
         2. THEN a `*` or `.` path or its name works too
      2. WHEN a product has no container
         1. THEN it fails naming its entry point or target or deploy file
         2. THEN two containers claiming it fail naming both claims
      3. WHEN a `Container` matches no product
         1. THEN it fails naming the view line
         2. THEN paths holding only source layers fail the same way
      4. WHEN a container's technology does not name its product's language
         1. THEN it fails naming the view line
         2. THEN it fails naming the product
      5. WHEN an owned store names no path
         1. THEN it fails naming the store
         2. THEN an external system naming no platform service fails too
      6. WHEN a manifest cannot be parsed
         1. THEN it fails naming that line
         2. THEN a target lacking a type or source folder fails too
      7. WHEN a container's span names a package
         1. THEN it owns the package ignoring trailing punctuation
         2. THEN a label or exact entry point word owns it too
2. Scoped runs read what moved
   1. GIVEN Developer has a repo
   2. GIVEN it may hold a containers view
      1. WHEN the repo has no architecture views
         1. THEN it passes with zero files
      2. WHEN a manifest only lists packages
         1. THEN it proves no product
         2. THEN a test folder or `scripts/` file proves none
      3. WHEN `.gitattributes` marks a build manifest `linguist-generated`
         1. THEN it still reads that manifest
      4. WHEN a scoped run changes only tests or changelog or docs
         1. THEN it passes with zero files
      5. WHEN a scoped run changes a view or a product
         1. THEN it reads every view
         2. THEN it reads every product
