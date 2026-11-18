# 0081 Checks No Contrastive Reframing

## Why

I want docs that state the point directly, without first rejecting a claim nobody made.

## Measurement

markdown files with no reframed sentence
-
markdown files checked

## Requirements

- 0081.1
    - Given Developer has a doc with a negated sentence followed by a restating sentence
        - When no-contrastive-reframing checks the doc
            - Then it fails naming the file and quoting both sentences
- 0081.2
    - Given Developer has a doc with one sentence that negates, then pivots with a comma
        - When no-contrastive-reframing checks the doc
            - Then it fails quoting that sentence
- 0081.3
    - Given Developer has a doc with ordinary negations that restate nothing
        - When no-contrastive-reframing checks the doc
            - Then it passes
- 0081.4
    - Given Developer has the pattern in front matter, fenced code, a heading, or inline code
        - When no-contrastive-reframing checks the doc
            - Then it passes, since that text is not prose
- 0081.5
    - Given Developer has the pattern inside straight or curly quotes
        - When no-contrastive-reframing checks the doc
            - Then it passes, so a doc can cite the pattern
- 0081.6
    - Given Developer has a reframed sentence longer than 80 characters
        - When no-contrastive-reframing checks the doc
            - Then it quotes the first 80 characters and an ellipsis
