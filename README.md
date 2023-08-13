# Gaffer
 
## [Under Development]

### What's GoCondor's Gaffer?
`Gaffer` is `GoCondor's` `cli` tool, it helps you create new projects and run the `live reloading` development server.

![Build Status](https://github.com/gocondor/gaffer/actions/workflows/build.yml/badge.svg)
![Test Status](https://github.com/gocondor/gaffer/actions/workflows/test.yml/badge.svg)


### Install
To install run the following command:
```bash
go install github.com/gocondor/gaffer@latest
```

## Create a new project:
To create a new project run the following command:
```bash
gaffer new [my-project] [github.com/my-organization/my-project]
```

## Run the live reloading development server
To start the development server run the follwing command:
```bash
gaffer run:dev
```
