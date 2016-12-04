# appdescs

Coming up with project ideas can be difficult&mdash;so difficult that I want my computer to do it for me.

The goal of this project is to gather a ton of textual app descriptions. I can then train a language model on these descriptions in the hopes that it will learn to generate amusing app ideas. The language model probably won't be part of this repository, since I can just use my [char-rnn](https://github.com/unixpickle/char-rnn) to start off. However, I will still hopefully post some results back to this repository.

# Usage

Currently, this is simply a tool for downloading Play Store app descriptions.

First, install and configure [Go](https://golang.org/doc/install). Once you have Go, you can download the code like so:

```
$ go get github.com/unixpickle/appdescs/play_store
```

Afterwards, you can switch into the directory you just fetched:

```
$ cd $GOPATH/src/github.com/unixpickle/appdescs/play_store
```

Finally, to run the command, you will need to pick an apps list page and an output directory. The output directory will be created if it does not already exist, and will be filled with HTML files (one per app). Here is an example:

```
$ go run *.go https://play.google.com/store/apps/collection/topselling_free?hl=en out_dir
```

This will run until you terminate it or it runs out of apps (most likely the former).
