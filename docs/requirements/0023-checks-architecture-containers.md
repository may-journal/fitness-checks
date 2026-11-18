# 0023 Architecture Containers

## Why

My container view matches what I ship. A new server, a missing native target, or a layer drawn as a runtime fails.

## Measurement

container views whose containers match build targets, entry points, and deploy files
-
container views checked

## Requirements

- 0023.1
    - Given Developer has a native app whose view declares its app and extension targets beside Node tools
        - When I check its containers
            - Then it passes
- 0023.2
    - Given Developer has a Go command, Java application, Node service, or deploy file the view does not declare
        - When I check its containers
            - Then the error names the entry point
- 0023.4
    - Given Developer has a container whose technology is not its product's language
        - When I check its containers
            - Then the error names the container and the evidence
- 0023.6
    - Given Developer has an owned store naming no path, or a platform service no external system names
        - When I check its containers
            - Then the error names the store line or the import line
- 0023.7
    - Given Developer has a target or package manifest the check cannot read
        - When I check its containers
            - Then the error names the manifest
- 0023.9
    - Given Developer has a container naming the root path `.`, or a target manifest marked `linguist-generated`
        - When I check its containers
            - Then it matches the products built there
- 0023.10
    - Given Developer has a container whose build target no longer exists, or whose paths hold only source layers
        - When I check its containers
            - Then the error names the container line and says to show source layers in the component view
- 0023.11
    - Given Developer has a published npm library whose entry points are tracked or built, beside a private or npm init package
        - When I check its containers
            - Then a span or label naming the package or an entry passes, a description needs the exact entry, and the others need none
