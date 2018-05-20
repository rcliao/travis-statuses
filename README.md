# Travis CI Statuses

A quick CLI app to grab latest build statuses given organization name and a collection of repository names.

## Usage

```
go run main.go --orgName {orgname} --repoNamesFile {reponamesfiles} --token {travisci_auth_token}
```

An example of the `repoNamesFiles` for this repo:

```
travis-statuses
```
