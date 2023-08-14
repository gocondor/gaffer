# Gaffer
 
## [Under Development]

### What's Gaffer?
`Gaffer` is `GoCondor's` `cli` tool, it helps you create new projects and perform other tasks like run your app in the `live reloading` mode for development purpose.

![Build Status](https://github.com/gocondor/gaffer/actions/workflows/build.yml/badge.svg)
![Test Status](https://github.com/gocondor/gaffer/actions/workflows/test.yml/badge.svg)


### Install
To install `Gaffer` run the following command:
```bash
go install github.com/gocondor/gaffer@latest
```

## Create a new project:
To create a new project run the following command:
```bash
gaffer new [project-name] [project-remote-repository]
```
example:
```bash
gaffer new myapp github.com/gocondor/myapp
```

## Run the app in the live reloading mode
To start the app in the `live reloading` mode for development, first cd into the project directory, then run the following command:
```bash
gaffer run:dev
```
