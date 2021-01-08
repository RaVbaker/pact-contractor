# Pact Contractor

A piece of software that helps to store Pact contracts on AWS S3 storage. Coordinates the retrieval and pushes them in appropriate structure. The promise is that it can replace [pact broker](https://github.com/pact-foundation/pact_broker) software. 

--- 
## Basics

Pact contracts are built to be used with PactBroker software. Which when managed is making some troubles. You have to manage the app, keep the DB backups, update, patch and remember to match version of broker and your clients in various apps. This app is meant to replace it with reliable and solid solution.

Since what's most important is a shared storage the PactContractor is storing them in AWS S3 object store. In a structure that is ready for the version/branch-based lookups. 

In the future, it should also help in configuration of a "hook" that could run either locally when the new contract version is pushed or with an AWS Lambda function/ GitHub Action called to trigger variety of builds.

### Features and future ideas:

1. When `push` called it stores contracts from `pacts/{provider-name}/{consumer-name}/spec.json` files in S3 appropriate tag-based files. E.g. `pacts/{provider-name}/{consumer-name}/main.json` or `pacts/{provider-name}/{consumer-name}/develop.json`. 
2. Objects in S3 should use S3 versioning and metadata and tagging. The versioning would allow accessing uploaded in past contracts. Metadata would contain some meta information of the contract. Like: Author, CommitSHA, Branch + Origin context, a freeform value for BuildID etc.
3. The  Verification Status will be stored in Object Tags or separate S3 object. There should be a helper to push the verification status to S3.
4. It assumes it runs in Git environment where it can interfere tag from branch name (with `main` as default, so if legacy `master` branch is present it maps it to `main` name), extracts details like Author, CommitSHA and Branch name from Git data. All options can also be provided on a command using flags (useful in CI environment).
5. Config file can be kept not only globally per-user in $HOMEDIR but also in local `pacts/` directory and interfered if present.
6. Pulling of the contracts should allow some dynamic specTag detection.
7. Pulling is allows to download multiple files at once, when paths are separated by comma, version can be then provided after # sign in path. E.g. "paths/foo/bar/test.json#some-v3rsion-1D,paths/foo/baz/{branch}.json#oth3r-v3rsion-1D" 
8. There can be 3 "hook" modes. [@TODO] in version 1.1.0
    * Local - runs arbitrary code immediately after `push` (might require extra configuration in `config.yaml`).
    * AWS Lambda - that is based on S3 hooks executed when new file is uploaded/updated. The tool helps in marking the status of verification which is stored in S3 Object Tags and setting up the Lambdas in AWS.
    * GitHub Actions - similar to Local but runs the code/verification using GitHub Action. Also here it is planned to support such setup also on preparation part.
9. Contracts deep-merge might be needed if existing version of contract is defined as same commitsha. Especially in usecase like we have with split build runs.
10. `verify [path]` helper command that accepts as argument a command to run for verification and does the `pull [path]` & `submit [path] [verification-status]` around it.
11. `can-i-deploy` checks verification status on a path#version and prints details of the object and verification.
--- 
## Basic usage

* `pact-contractor pull [path]` Pulls pact contracts from configured S3 bucket
* `pact-contractor push` Push generated pact contracts to configured S3 bucket, (default path="pacts/*/*/spec.json")
* `pact-contractor submit [path] [verification-status]` To submit status of contract verification and store it in S3 Object Tag
* `pact-contractor list [path]` To show list of contracts. The path has to be valid glob path.
* `pact-contractor get-version [path]` Shows the versionID of the latest version of the stored contract
* `pact-contractor can-i-deploy [path]` Returns exitCode 0 if verification status is `"success"` and displays details of the path. Is also aliased with `get` keyword.
* `pact-contractor verify [path] --cmd "command to verify {path}` Helper to run pull, then the command to verify the path and finally submits new verification result and removes pulled local file.  

The `--help` flag says even more about every command. 

Configuration flags are: 

* `-b, --bucket string`   AWS S3 Bucket name
* `-r, --region string`   AWS S3 Region name (default from `$AWS_REGION` env variable)
* `--config string`   config file (default is $HOME/.pact-contractor.yaml)
* `--cmd string`  Command to execute during verification, {path} is replaced with provided [path]. (default `"bundle exec rake pact:verify:at[{path}]"`)

The `bucket` and `region` can be also configured in the config file. Either globally or provided with `--config` flag when executed.

You can also provide AWS variables with in the file for values like:
`AWS_PROFILE, AWS_ACCESS_KEY, AWS_ACCESS_KEY_ID, AWS_SECRET_KEY, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN, AWS_CONFIG_FILE, AWS_SHARED_CREDENTIALS_FILE, AWS_ROLE_ARN, AWS_CA_BUNDLE, AWS_REGION, AWS_DEFAULT_REGION, AWS_SDK_LOAD_CONFIG`. Or provide prefixed with `PACT_` values so they will be overwritten for the run of the pact-contractor app only.

The file sample (`~/config.yaml`):
```yaml
---
aws_region: us-east-1
bucket: mybucket
cmd: echo "{path}" && bundle exec rake pact:verify:at[{path}]
hooks:
- type: local
  path_filter: pacts/*/*.json
  spec:
     command: echo "foo {path} OK"
```

## Makefile

The `make` command can be used for basic usage like `make build`, `make run` (e.g. `make run cmd="pull --help"`) and `make help` to get list of available commands.

---

## How to use it?

1. Run your tests as usual, let them generate updated pact contracts. Let's say the pact contracts were stored under path `spec/pacts/consumerA-*.json` files.
2. After successful test run execute the push command and provide bucket/region params: 

```bash
   pact-contractor push spec/pacts/consumerA-* -b my-bucket
2021/01/04 13:40:47 For path "spec/pacts/consumerA-*.json" detected files: [spec/pacts/consumerA-producerA.json spec/pacts/consumerA-producerB.json spec/pacts/consumerA-producerC.json]
Successfully uploaded "spec/pacts/consumerA-producerA/main.json" [version: "zmhkM4VNFv6BD9lolilHps_ODxkY5eX_"] to "my-bucket"
Successfully uploaded "spec/pacts/consumerA-producerB/main.json" [version: "uIhyCu9rMQCTvX_s3AcRw9gehuOdObhO"] to "my-bucket"
Successfully uploaded "spec/pacts/consumerA-producerC/main.json" [version: "zGzJPy3qAZW6MC6sFIMhe9LihsmEWV9l"] to "my-bucket"   
```

3. Now under producer app producerB or any other you can run the verification step:

```bash
pact-contractor verify -b  my-bucket spec/pacts/consumerA-producerB/{branch}.json  --provider-version "1.2.3" --provider-context  "MY LOCAL BUILD test 123"
Successfully downloaded "spec/pacts/consumerA-producerB/main.json" from bucket "my-bucket" to file "spec/pacts/consumerA-producerB.json", 1732 bytes
Executing command: `bundle exec rake pact:verify:at[spec/pacts/consumerA-producerB.json]`
SPEC_OPTS='' /Users/ravbaker/.asdf/installs/ruby/2.5.1/bin/ruby -S pact verify --pact-helper /Users/ravbaker/Code/producerB/spec/service_consumers/pact_helper.rb --pact-uri spec/pacts/consumerA-producerB.json
INFO: Reading pact at spec/pacts/consumerA-producerB.json
Verifying a pact between ConsumerA and ProducerB
  Given Something happend
    an event with something changed
      has matching content
1 interaction, 0 failures
Marked as "success" all paths: "spec/pacts/consumerA-producerB/main.json" in bucket my-bucket
```

4. You can verify the state of the run with `can-i-deploy` command:

```bash
pact-contractor can-i-deploy spec/pacts/consumerA-producerB/main.json -b my-bucket --provider-version 1.2.3
Examinating path: "spec/pacts/consumerA-producerB/main.json", version ID: "uIhyCu9rMQCTvX_s3AcRw9gehuOdObhO"

Provider Version: "1.2.3"
Pact Verification: "success"
Branch: "main"
Commitsha: "49fc9285db54dc582d91267397bf84f5353a3ae8"
Author: "My Name"
VersionID: "uIhyCu9rMQCTvX_s3AcRw9gehuOdObhO"
Last Modified: "2021-01-04T12:40:48Z"
Provider Context: "MY LOCAL BUILD test 123"
```

---

(c) Copyright Rafal "RaVbaker" Piekarski 2020-2021
