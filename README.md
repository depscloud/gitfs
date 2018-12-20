## gitfs

_Better repository management for the everyday developer._

When starting a job at any company, you often spend a fair amount setting up your development environment.
This involves installing dependent tooling, cloning repositories, amongst many other things.
`gitfs` was developed to ease some of that process by removing the need to manually clone down every repository.
Instead, it creates a virtual directory structure around git urls, making it easy to navigate sources for projects.

## Support

Support on this library is limited right now, but open for active development.

### Remotes

Currently, it supports a generic endpoint, but are open to supporting new remotes for various integrations.
We are planning to support Github, Gitlab, and Bitbucket in the short term.
PR's for each of these remotes based on their [config](pkg/config/config.proto) definitions are welcome.

### In Memory Storage

For the proof of concept, this project leverages an in memory file system to store the cloned repositories.
This leverages the billy library, making it really easy to swap the underlying git filesystem for an alternative implementation.
Support for on persistent stores will be added for better long term support.

## Sample File System Mount

```
{mount}/
  {host}/
    {user}/
      {project}/
        .git/
        ...
    {group}/
      {project}/
        .git/
        ...
  {host}/
    {user}/
      {project}/
        .git/
        ...
    {group}/
      {project}/
        .git/
        ...
```
