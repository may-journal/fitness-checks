# 0001 Containers Match What Ships

## Why

I get my containers view matched against the runtime products my repo builds and deploys.

## Measurement

repos whose containers view matches their products

-

repos with a containers view

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a containers view AND products proven by build targets or entry points or deploy files
  WHEN each product has one container naming its folder or a `*` or `.` path or its name
  THEN it passes
    AND it counts the container views
  WHEN a product has no container
    AND two containers claim a product
  THEN it fails naming its entry point or target or deploy file
    AND it names both claims
  WHEN a `Container` matches no product
    AND a container's paths hold only source layers
  THEN it fails naming the view line
  WHEN a container's technology does not name its product's language
  THEN it fails naming the view line
    AND it names the product
  WHEN an owned store names no path
    AND no external system names a platform service used
  THEN it fails naming the store or the source line
  WHEN a manifest cannot be parsed
    AND a target lacks a type
    AND a target lacks a source folder
  THEN it fails naming that line
  WHEN a container's span or label or exact entry point word names a package
  THEN it owns the package
    AND it ignores trailing punctuation
    BUT other description words never do

SCENARIO: 0001.2
  GIVEN Developer has a repo with or without a containers view
  WHEN the repo has no architecture views
  THEN it passes with zero files
  WHEN a manifest only lists packages
    AND a file sits in a test folder
    AND a file sits under `scripts/`
  THEN it proves no product
  WHEN `.gitattributes` marks a build manifest `linguist-generated`
  THEN it still reads that manifest
  WHEN a scoped run changes only tests or changelog or docs
  THEN it passes with zero files
  WHEN a scoped run changes a view or a product
  THEN it reads every view
    AND it reads every product
```
