![golucy](https://raw.github.com/philipsoutham/golucy/devel/artwork/golucy.png)

# golucy
Go bindings for [Apache Lucy][1]. The [Apache Lucy][1] search engine library provides full-text search for dynamic
programming languages. It is a "loose C" port of the [Apache Lucene™][2] search engine library for Java.


## Dependencies

### Go
Duh.

### Lucy
Works as of commit [e19687f9a6b0158308ac7bcafc663296635b107a][3].
+ `export BUILD_DIR=$HOME/build`
+ `export LUCY_HOME=$HOME/.local/lucy`
+ `cd $BUILD_DIR`
+ `git clone https://git-wip-us.apache.org/repos/asf/lucy.git`
+ `cd $BUILD_DIR/lucy/c`
+ `./configure --prefix=$LUCY_HOME`
+ `make && make test && make install`
+ `./install --prefix=$LUCY_HOME`
+ `cd $BUILD_DIR/lucy/clownfish/runtime/c`
+ `./configure --prefix=$LUCY_HOME`
+ `make && make test && make install`
+ `./install --prefix=$LUCY_HOME`

### Configuration
Add the following to your `.profile` or `.zshrc` or similar (you will also need to have your `GOHOME` and/or `GOPATH` set).
+ `export CGO_LDFLAGS="-L$LUCY_HOME/lib -llucy -lcfish ${CGO_LDFLAGS}"`
+ `export CGO_CFLAGS="-I$LUCY_HOME/include ${CGO_CFLAGS}"`
+ `export LD_LIBRARY_PATH=$LUCY_HOME/lib:$LD_LIBRARY_PATH`


[1]: http://lucy.apache.org/
[2]: http://lucene.apache.org/core/
[3]: https://git-wip-us.apache.org/repos/asf?p=lucy.git;a=commit;h=e19687f9a6b0158308ac7bcafc663296635b107a
