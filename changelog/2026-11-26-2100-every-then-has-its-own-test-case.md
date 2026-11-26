### 2026.11.26.2100

> Every THEN in a requirement doc has exactly one test case named with its full ID and its exact text.

- Feat: `requirements` fails a THEN with no case, naming the THEN and the test to add it to. The case opens with the THEN's ID, such as `0005.3.1.2.1`, then the THEN's words without `THEN `.
- Feat: `requirements` fails a case whose words differ from its THEN, a THEN with two cases, a case naming a path no doc holds, and a case with no assertion call.
- Feat: Swift cases count through XCTest activities. Node cases stay exempt behind the `nodeReportsCases` switch until fitness-runner#1108 merges.
- Docs: 0005.3 is a `Removed` line and 0005.4 states the rule. The README template shows a test with one case per THEN.
- Docs: every requirement test here names a case for each THEN of its docs.
