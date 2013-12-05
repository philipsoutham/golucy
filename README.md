![golucy](https://raw.github.com/philipsoutham/golucy/master/_artwork/golucy.png)

# golucy
Go bindings for [Apache Lucy][1]. The [Apache Lucy][1] search engine library provides full-text search for dynamic
programming languages. It is a "loose C" port of the [Apache Lucene™][2] search engine library for Java.


## Dependencies

### Go
Duh.

### Lucy
Works as of commit [e19687f9a6b0158308ac7bcafc663296635b107a][3].
```shell
$ export BUILD_DIR=$HOME/build
$ export LUCY_HOME=$HOME/.local/lucy
$ cd $BUILD_DIR
$ git clone https://git-wip-us.apache.org/repos/asf/lucy.git
$ cd $BUILD_DIR/lucy/c
$ git checkout e19687f9a6b0158308ac7bcafc663296635b107a
$ ./configure --prefix=$LUCY_HOME
$ make && make test && make install
$ ./install.sh --prefix $LUCY_HOME
$ cd $BUILD_DIR/lucy/clownfish/runtime/c
$ ./configure --prefix=$LUCY_HOME
$ make && make test && make install
$ ./install.sh --prefix $LUCY_HOME
```
### Configuration
Add the following to your `.profile` or `.zshrc` or similar (you will also need to have your `GOHOME` and/or `GOPATH` set).
```bash
export LUCY_HOME=$HOME/.local/lucy
export CGO_LDFLAGS="-L$LUCY_HOME/lib -llucy -lcfish ${CGO_LDFLAGS}"
export CGO_CFLAGS="-I$LUCY_HOME/include ${CGO_CFLAGS}"
export LD_LIBRARY_PATH=$LUCY_HOME/lib:$LD_LIBRARY_PATH
```

## Installation
Provided you have the dependencies in order add this:
```go
import (
  golucy "github.com/philipsoutham/golucy/v0.0.1"
)
```
then do:
```shell
$ go get 
```


## Example
See [this example][4], inspired by [this][5] one in C. 
If you're running [docker][7] you can see it in action like so. 


```shell
$ docker pull psoutham/golucy
$ docker run psoutham/golucy
```
Details on the docker image can be found [here][8].

[1]: http://lucy.apache.org/
[2]: http://lucene.apache.org/core/
[3]: https://git-wip-us.apache.org/repos/asf?p=lucy.git;a=commit;h=e19687f9a6b0158308ac7bcafc663296635b107a
[4]: https://gist.github.com/philipsoutham/6395791
[5]: https://git-wip-us.apache.org/repos/asf?p=lucy.git;a=blob;f=c/sample/getting_started.c;h=efdf050d546d7fbfaadd6882f2f6b91033c8be69;hb=HEAD
[6]: https://gist.github.com/philipsoutham/6371770
[7]: http://docker.io
[8]: https://index.docker.io/u/psoutham/golucy/
