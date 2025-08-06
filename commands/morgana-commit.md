# MORGANA-COMMIT Command

Add all changes to staging, create a commit, and push to remote.

Follow this checklist for writing your commit message: <type>[optional scope]: <description>
[optional body] [optional footer(s)]

- commit SHOULD contain the following structural elements to communicate intent:
  fix: a commit of the type fix patches a bug in your codebase (this correlates
  with PATCH in Semantic Versioning). feat: a commit of the type feat introduces
  a new feature to the codebase (this correlates with MINOR in Semantic
  Versioning). BREAKING CHANGE: a commit that has a footer BREAKING CHANGE:, or
  appends a ! after the type/scope, introduces a breaking API change
  (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part
  of commits of any type. types other than fix: and feat: are allowed, for
  example @commitlint/config-conventional (based on the Angular convention)
  recommends build:, chore:, ci:, docs:, style:, refactor:, perf:, test:, and
  others. footers other than BREAKING CHANGE: <description> may be provided and
  follow a convention similar to git trailer format.

DO NOT add any information about this being committed by claude code
