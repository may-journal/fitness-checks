### 2026.11.26.2000

> Old-style requirements fail wherever they hide, and every THEN has its own test case.

- Feat: `readme-behavior` fails a GIVEN, WHEN or THEN list item in any case outside a code fence, naming the line.
- Feat: `requirements` fails a Gherkin word or block keyword not in caps inside a step, 0008.
- Feat: `requirements` fails a THEN with no case named with its ID, such as `0005.3.1.2.1`, plus its words, a stray case, or a case with no assertion call.
- Feat: Swift cases count through XCTest activities, and Node stays exempt behind the `nodeReportsCases` switch.
- Docs: 0005.3 is a `Removed` line, 0005.4 states the rule, and every requirement doc and test here follows the new rules.
