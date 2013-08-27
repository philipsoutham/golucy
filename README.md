# golucy ![golucy](https://raw.github.com/philipsoutham/golucy/master/artwork/golucy.png)

Go bindings for Apache Lucy. The Apache Lucy search engine library provides full-text search for dynamic programming languages. It is a "loose C" port of the Apache Lucene™ search engine library for Java.


## Dependencies

### Go
Duh.

### Lucy
Works as of commit [e19687f9a6b0158308ac7bcafc663296635b107a](https://git-wip-us.apache.org/repos/asf?p=lucy.git;a=commit;h=e19687f9a6b0158308ac7bcafc663296635b107a).
`export BUILD_DIR=$HOME/build`
`export LUCY_HOME=$HOME/.local/lucy`
`cd $BUILD_DIR`
`git clone https://git-wip-us.apache.org/repos/asf/lucy.git`
`cd lucy/c`
`./configure --prefix=$LUCY_HOME`
`make && make test && make install`
`./install --prefix=$LUCY_HOME`
`cd $BUILD_DIR/lucy/clownfish/runtime/c`
`./configure --prefix=$LUCY_HOME`
`make && make test && make install`
`./install --prefix=$LUCY_HOME`

### Configuration
Add the following to your `.profile` or `.zshrc` or similar (you will also need to have your `GOHOME` and/or `GOPATH` set).
`export CGO_LDFLAGS="-L$LUCY_HOME/lib -llucy -lcfish ${CGO_LDFLAGS}"`
`export CGO_CFLAGS="-I$LUCY_HOME/include ${CGO_CFLAGS}"`
