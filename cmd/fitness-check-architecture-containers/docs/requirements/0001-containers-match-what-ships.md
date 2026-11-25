# 0001 Containers Match What Ships

## Why

I get my containers view matched against the runtime products my repo builds and deploys.

## Measurement

repos whose containers view matches their products

-

repos with a containers view

## Requirements

- 0001.1
  - Given Developer has a containers view and products proven by build targets, entry points or deploy files
    1. When each product has one container naming its folder, a `*` or `.` path, or its name
       - Then it passes, counting the container views
    2. When a product has no container, or two containers claim it
       - Then it fails naming its entry point, target or deploy file, and both claims
    3. When a `Container` matches no product, or its paths hold only source layers
       - Then it fails naming the view line
    4. When a container's technology does not name its product's language
       - Then it fails naming the view line and the product
    5. When an owned store names no path, or no external system names a platform service used
       - Then it fails naming the store or the source line
    6. When a manifest cannot be parsed, or a target lacks a type or source folder
       - Then it fails naming that line
    7. When a container's span, label or exact entry point word names a package
       - Then it owns the package, ignoring trailing punctuation; other description words never do
- 0001.2
  - Given Developer has a repo, with or without a containers view
    1. When the repo has no architecture views
       - Then it passes with zero files
    2. When a manifest only lists packages, or a file sits in a test folder or under `scripts/`
       - Then it proves no product
    3. When `.gitattributes` marks a build manifest `linguist-generated`
       - Then it still reads that manifest
    4. When a scoped run changes only tests, changelog or docs
       - Then it passes with zero files
    5. When a scoped run changes a view or a product
       - Then it reads every view and every product
