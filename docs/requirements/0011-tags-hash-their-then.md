# 0011 Tags Hash Their THEN

## Why

I tie a test case to its THEN by a hash of its words, so a reworded THEN is noticed.

## Measurement

tagged cases whose tag holds the current hash

-

tagged cases

## Requirements

1. A tag holds the hash of its THEN
   1. GIVEN Developer has the words of a THEN
      1. WHEN the tag hashes them
         1. THEN the hash is the first 8 hex digits of their sha256
         2. THEN a change of whitespace leaves the hash as it was
2. A case makes an assertion
   1. GIVEN Developer has a Go test that calls `req.Then`
      1. WHEN its body makes no assertion
         1. THEN the test fails naming the THEN
      2. WHEN its body asserts
         1. THEN the test passes
3. A body runs as it was written
   1. GIVEN Developer has a Go test that calls `req.Then` with a body of its own `*testing.T`
      1. WHEN the body fails the test with `t.Errorf`
         1. THEN the test fails with that message
      2. WHEN the body asserts through a helper that calls `req.Asserted`
         1. THEN the test passes
