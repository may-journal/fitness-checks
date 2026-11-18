# 0028 Checks Issue Checklist

## Why

An issue with unfinished items is never closed as done, so its open work stays visible.

## Measurement

issue bodies with unchecked items reported red
-
issue bodies with unchecked items

## Requirements

- 0028.1
    - Given Developer has an issue body with an unchecked item
        - When fitness runs issue-checklist on it
            - Then it reports red quoting the item
- 0028.2
    - Given Developer has an issue body with every item ticked
        - When fitness runs issue-checklist on it
            - Then it reports green
