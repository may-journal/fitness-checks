# 0001 Subjects Carry Type And Scope

## Why

I read a branch's history and see what each change was, because every commit subject names its type and scope.

## Measurement

commit subjects that name a listed type and a scope

-

commit subjects judged

## Requirements

1. Subjects follow Conventional Commits
   1. GIVEN Developer has a commit message or a pull request title
      1. WHEN the first line is `type(scope): description` with a listed type or a `Merge` line
         1. THEN it passes
         2. THEN only the first line is judged
      2. WHEN the subject has no scope
         1. THEN it fails quoting the subject
         2. THEN it gives the suggested format
         3. THEN it gives the allowed types
      3. WHEN the subject has an unlisted type or a breaking mark
         1. THEN it fails naming the subject with the format to use
      4. WHEN the message is empty or no repo can be read
         1. THEN it fails with a no-message error
2. A full run judges only the branch's own commits
   1. GIVEN Developer has a branch over history that `origin` holds
      1. WHEN the repo has no origin refs
         1. THEN the HEAD commit subject is judged
      2. WHEN the default branch is `main` or `master` or one Fitness cannot guess
         1. THEN only commits origin lacks are judged
         2. THEN a bad inherited subject is never quoted
      3. WHEN the branch has no commits of its own
         1. THEN it passes
3. The pre-commit hook judges earlier commits
   1. GIVEN Developer has the installer's pre-commit hook alone on a fork branch
      1. WHEN a commit has a valid subject
         1. THEN the commit succeeds
         2. THEN the inherited parent subject is not judged
      2. WHEN an earlier commit on the branch has a bad subject
         1. THEN the next commit fails quoting that subject
