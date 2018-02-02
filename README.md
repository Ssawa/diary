# Diary

This is meant to be a very small and straight forward application that, when called, creates a file for the given day with a template of your choice and, if configured, backs it up to an S3 bucket.

## Running
```
diary
```

### Example File
```
# Friday February 2, 2018

## At 9:44am...
This is one update that I wrote blah blah blah

## At 10:29am...
I decided to make another update, this is it

## At 5:11pm...
Almost the end of the day now

## At 5:14pm...
I forgot something
```

## Configuration

**config** (`--config`) [`~/.diary.yml`] : Path to a config file to use (configuration is done using [Viper](https://github.com/spf13/viper))

**verbose** (`--verbose`, `-v`) [`false`]: Whether or not to show verbose logging

**editor** (`--editor`, `-v`, environment variable: `EDITOR`) [`vim`]: The editor to start for editing files. This command should stay alive until your done editing (for instance, if you're using Visual Studio Code's `code` commandline util use `code -w`)

**file.base** [`~/.diary`]: The base path to store diary entries on your local system

**file.template.path** [`2006/1/2-Mon-Jan-2006.md`]: The format to use when generating a new file (see https://golang.org/pkg/time/#Time.Format for how this formatting works)

**file.template.new** [`# Monday January 2, 2006\n`]: [Golang format template](https://golang.org/pkg/time/#Time.Format) that will be inserted at the beginning of new files

**file.template.append** [`\n## At 3:04pm...\n`]: [Golang format template](https://golang.org/pkg/time/#Time.Format) that will be inserted each time the file is opened

**s3.enable** [`false`]: Whether or not to use S3 for backing up

**s3.region**: The region to use for S3 uploading (required if `s3.enabled` is set)

**s3.id**: AWS API id to use for uploading (required if `s3.enabled` is set)

**s3.secret**: AWS API secret to use for uploading (required if `s3.enabled` is set)

**s3.bucket**: The S3 bucket to upload to (required if `s3.enabled` is set)

**s3.key_prefix**: An optional prefix to use for uploaded files

### Example `.diary.yml`
```
verbose: no
editor: code -wn
file:
  base: /Dropbox/diary
  template:
    path: 2006-Jan/2_Mon.txt
    new: "What I did today...\n"
    append: ""
s3:
  enabled: yes
  region: "us-east-1"
  id: "totally"
  secret: "a secret"
  bucket: my-diary-backup-bucket
```
