# consync
consync is a tool to help sync or diff Consul key/value data across multiple datacenters.

Please note that consync is in early development and does not have full test coverage yet.  It
seems to be working correctly but until tests are implemented completely use with caution!


## Installation

consync uses [glide](https://github.com/Masterminds/glide) to handle dependencies, so be sure to
have that installed first.  Once you have glide installed, install the dependencies required to
build consync by typing:

```bash
$ glide install
```

__Note:__ If you have a problem installing packages from `gopkg.in` when installing dependencies,
see [gopkg #50](https://github.com/niemeyer/gopkg/issues/50) for more information and a workaround.

After the dependencies have been installed, consync can be built and installed using:

```bash
$ go install
```

This will put the consync binary in `$GOPATH/bin` so you can use it as a normal command if you have
that directory on your execution path.


## Usage

One thing to note with consync is that the URLs it expects on the command line are not the same as
the URLs used if you were to connect to the Consul REST API.  consync URLs only require the path
from the beginning of the key/value path instead of the full REST path.  For example, if you
want to see the difference between `/` on one server and `/my/path` on another you would use
something like `http://10.10.10.10:8500/` for the "from" and `http://10.10.10.11:8500/my/path` for
the "to".

### Displaying Differences

consync provides a `diff` command to display the differences between two paths in Consul.  These paths
can be in the same data center at different paths, different data centers at the same path or even
different data centers at different paths.

For example, to view the difference between the same path for two data centers in the same cluster,
you could do:

```bash
$ consync diff \
    -f http://consul.host/app/settings --from-dc dc1 --from-token b3e...9a2 \
    -t http://consul.host/app/settings --to-dc dc2 --to-token b3e...9a2
```

### Syncing Changes

consync also provides a `sync` command to modify a target path to match a source path.  This includes
adding any new keys, updating any modified keys and removing any keys that exist in the target
but not in the source.  This command functions similar to the `diff` command as far as the parameters
you can provide but it also includes an additional `--execute` (or `-e`) parameter. If the
execute parameter is not provided, consync will only display the changes it would make instead of
actually applying them.  This is to make sure changes aren't accidentally applied to a target.

An example of using the `sync` command is as follows:

```bash
$ consync sync \
    -f http://consul.host/app/settings --from-dc dc1 --from-token b3e...9a2 \
    -t http://consul.host/app/settings --to-dc dc2 --to-token b3e...9a2 \
    --execute
```
