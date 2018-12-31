Contributing
===

### prerequisites

* [Opening a Pull Request](https://help.github.com/articles/about-pull-requests/)

### as a technical writer

Like many projects, technical documentation is often needed.
Instead of restricting contributions strictly to engineers, we want to encourage contributions from all walks of life.
Documentation is open to contributions from the public.
They can be found under the `docs` folder and hosted on github here: https://mjpitz.github.io/gitfs
Changes to documentation do not require tickets in advanced, but do require pull requests.

### as a developer / engineer

All work items should have a corresponding issue under the github issue manager.
The first commit associated tickets must be prefixed with the github issue number as such: `gh-{issue}: <message>`.
All following commits to the branch do not require that prefix as they will be squashed before merge.
Your first commit will be used as the commit message, unless project owners desire a more clear message.

**Dependency Management**

```
go get -u <dependency>
go mod tidy
```

**Testing**

```
go vet ./...
go test -v -race ./...
```

**Building / Installing**

```
go build
go install
```

**Issue Reporting**

If you believe you've found an issue, first check existing issues for a similar issue.
If no issue has been opened, then please submit a ticket.
After a ticket is opened, you are welcome to open a PR for review.
All development on this issue must follow the guidelines above.

### as a project owner

As a project owner, you will be responsible for managing releases of the application.
To make life easy, the project uses `npm` and `package.json` for version management.
This is because managing versions using this method is trivial.

```
npm version patch|minor|major
git push
git push --tags
```

A passing build is required for publishing.
There are currently no hooks that prevent this, so it is your responsibility to check the build. 
