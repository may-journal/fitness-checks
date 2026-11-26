### 2026.11.26.2100

> Every THEN in a requirement doc has exactly one test case named with its full ID and its exact text.

- Feat: `requirements` fails a THEN with no case, naming the THEN and the test to add it to.
- Feat: a case opens with the THEN's ID, such as `0005.3.1.2.1`, then the THEN's words, or it fails.
- Feat: `requirements` fails a case naming a path no doc holds, or holding no assertion call.
- Feat: Swift cases count through XCTest activities, and Node stays exempt behind the `nodeReportsCases` switch.
- Docs: 0005.3 is a `Removed` line, 0005.4 states the rule, and every requirement test here names a case per THEN.
