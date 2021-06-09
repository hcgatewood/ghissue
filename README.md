# ghissue

Bulk-upload GitHub Issues as easy as

```bash
ghissue --token gh.token create issues.txt
```

## Example

Prerequisite: create a [GitHub personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) and save it to `./gh.token`

We'll create two issues: a self-assigned feature request and a bug report

```bash
$ cat issues.txt

hcgatewood23/test
---
Update the README | feature | hcgatewood23
Update readme with examples.

Add note that Issue body can contain multiple lines.
---
Fix CLI bug | bug
Bug report: CLI needs to be fixed. Someone please claim.
```

Next, bulk-create the Issues. This will also open your browser to view the newly-created Issues

```bash

$ ghissue --token gh.token create issues.txt --open

14,15
```

See `testdata/` for more examples.

## Notes

### Input format

The input file contains a repo target, followed by hyphen-separated issues.

The first line of each issue contains metadata, while all following lines comprise the Issue body. Labels, assignees, and body are optional.

```
repo_owner/repo_name
---
Title | Labels | Assignees
Body
---
Title | Labels | Assignees
Body
---
Title | Labels | Assignees
Body
```

### Additional options

CLI flags for the `create` command

```
--byline   Append hcgatewood/ghissue byline to Issue body (default true)
--dryrun   Don't actually create the Issues
--help     help for create
--info     Print more info about the Issues
--open     Open browser to view new Issues
```
