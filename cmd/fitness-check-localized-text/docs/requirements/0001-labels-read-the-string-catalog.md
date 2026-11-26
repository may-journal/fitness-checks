# 0001 Labels Read The String Catalog

## Why

I get every label in my app and its tests read from the String Catalog so no language shows a missing entry.

## Measurement

labels read through the String Catalog

-

labels in SwiftUI views and requirement tests

## Requirements

1. View literals skip no String Catalog entry
   1. GIVEN Developer has an app that tracks a `.xcstrings` file
   2. GIVEN the app has SwiftUI views
      1. WHEN a literal goes to a text view such as `Text`
         1. THEN it fails naming the file with its line
         2. THEN it names the literal
         3. THEN it names the call
         4. THEN it names the entry or reference to use
      2. WHEN a literal goes to a modifier such as `navigationTitle`
         1. THEN it fails naming that modifier
      3. WHEN a literal is a branch of a ternary or `??` expression
         1. THEN it fails
      4. WHEN a literal is a ternary branch inside an interpolation
         1. THEN it fails
         2. THEN a reference in the interpolation passes
      5. WHEN a literal spans several lines
         1. THEN it is named whole
      6. WHEN a view passes a reference or a value
         1. THEN it passes
         2. THEN a literal without letters passes
2. Test literals never retype a catalog label
   1. GIVEN Developer has an app with a String Catalog
   2. GIVEN the app has requirement tests
      1. WHEN a test types a literal equal to a catalog entry
         1. THEN it fails naming the file with its line
         2. THEN it names the literal
         3. THEN it names the entry
         4. THEN it names the reference to read it through
      2. WHEN a test types a caption
         1. THEN it passes
         2. THEN a price passes alike
         3. THEN a literal assigned with `=` passes alike
         4. THEN a stale entry passes alike
3. The app configures what the check allows
   1. GIVEN Developer has an app with a String Catalog
   2. GIVEN the repo configures `localizedText`
      1. WHEN an allow pattern with a reason matches a whole literal
         1. THEN that literal passes
      2. WHEN an allow pattern has no reason
         1. THEN the check fails naming the pattern
         2. THEN a pattern that does not compile fails alike
      3. WHEN the repo tracks no `.xcstrings` file
         1. THEN the check passes with zero files
