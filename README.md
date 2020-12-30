# Pact Contractor

A piece of software that helps to store Pact contracts on AWS S3 storage. Coordinates the retrieval and pushes them in appropriate structure. The promise it that it can replace `pact broker` software. 

--- 
## Idea

Pact contracts are built to be used with PactBroker software. Which when managed is making some troubles. You have to manage the app, keep the DB backups, update, patch and remember to match version of broker and your clients in various apps. This app is meant to replace it with reliable and solid solution.

Since what's most important is a shared storage the PactContractor is storing them in AWS S3 object store. In a structure that is ready for the version/branch-based lookups. 

In the future, it should also help in configuration of a "hook" that could run either locally when the new contract version is pushed or with an AWS Lambda function/ GitHub Action called to trigger variety of builds.

### Features:

1. When `push` called it stores contracts from `pacts/{provider-name}/{consumer-name}/spec.json` files in S3 appropriate tag-based files. E.g. `pacts/{provider-name}/{consumer-name}/main.json` or `pacts/{provider-name}/{consumer-name}/develop.json`. 
2. Objects in S3 should use S3 versioning and metadata and tagging. The versioning would allow accessing uploaded in past contracts. Metadata would contain some meta information of the contract. Like: Author, CommitSHA, Branch + Origin context, a freeform value for BuildID etc.
3. The  Verification Status will be stored in Object Tags or separate S3 object. There should be a helper to push the verification status to S3. [@TODO]
4. It assumes it runs in Git environment where it can interfere tag from branch name (with `main` as default, so if legacy `master` branch is present it maps it to `main` name), extracts details like Author, CommitSHA and Branch name from Git data. All options can also be provided on a command using flags (useful in CI environment).
5. Config file can be kept not only globally per-user in $HOMEDIR but also in local `pacts/` directory and interfered if present.
6. Pulling of the contracts should allow some dynamic specTag detection.
7. Pulling is allows to download multiple files at once, when paths are separated by comma, version can be then provided after # sign in path. E.g. "paths/foo/bar/test.json#some-v3rsion-1D,paths/foo/baz/{branch}.json#oth3r-v3rsion-1D" 
8. There can be 3 "hook" modes. [@TODO]
    * Local - runs verification on provider immediately after `push` (might require extra configuration in `config.yaml`).
    * AWS Lambda - that is based on S3 hooks executed when new file is uploaded/updated. The tool helps in marking the status of verification which is stored in S3 Object Tags and setting up the Lambdas in AWS.
    * GitHub Actions - similar to Local but runs the code/verification using GitHub Action. Also here it is planned to support such setup also on preparation part.
9. Contracts deep-merge might be needed if existing version of contract is defined as same commitsha. Especially in usecase like we have with split build runs.

--- 
## Basic usage

* `pact-contractor pull [path]`        Pulls pact contracts from configured S3 bucket
* `pact-contractor push`        Push generated pact contracts to configured S3 bucket, (default path="pacts/*/*/spec.json")

Configuration flags are: 

* `-b, --bucket string`   AWS S3 Bucket name
* `-r, --region string`   AWS S3 Region name (default "us-east-1")
* `--config string`   config file (default is $HOME/.pact-contractor.yaml)

The `bucket` and `region` can be also configured in the config file. Either globally or provided with `--config` flag when executed.

The file sample (`~/config.yaml`):
```yaml
---
region: us-east-1
bucket: mybucket
```

## Makefile

The `make` command can be used for basic usage like `make build`, `make run` (e.g. `make run cmd="pull --help"`) and `make help` to get list of available commands.

---




