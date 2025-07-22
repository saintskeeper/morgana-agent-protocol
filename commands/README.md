# My thoughts

This repo is designed to create simple and reusable commands for users and code bases


## assumptions
- you have zen mcp set up
- you have linear integratrions (mcp)
- you have  a basic undestanding of docker

## Concept
1. Human workbench setup
1. Context
2. Plan
3. Work
4. Review
5. Test
6. Update
7. Repeat.

### Command
set_project (sets the current project in linear so we dont have to search) wills store this in ~/.claude/linear/data
prep_issue(given a set of requirements an issue can be created from a template and preped for development )
prep_branch(preps a branch with the naming schema required for ci/cd and sets context)
prep_context(given the feature, codebase and available information review with self(deepthink, then cal zen:analyzie to prep the work plan) then write this to a file, and  a comment in the issue labeled :CONTEXT:
prompt
work_context: given the context fill call zen