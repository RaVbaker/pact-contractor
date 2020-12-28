# Pact Contractor

A piece of software that helps to store Pact contracts on AWS S3 storage. Coordinates the retrieval and pushes them in appropriate structure. The promise it that it can replace `pact broker` software. 

## Basic usage

* `pact-contractor pull [path]`        Pulls pact contracts from configured S3 bucket
* `pact-contractor push`        Push generated pact contracts to configured S3 bucket, (default path="pacts/provider/*/consumer/*/*.json")

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

The `make` command can be used for basic usage like `build` and `run` ( `make run cmd="pull --help"`).



