galaxy-install-extended
========================
ansible-galaxy wrapper to import *latest release* role.

I confused always by version control of many ansible roles... Don't you think the latest release must be importd by `version: latest`?


Features
---------
- if `version: latest` is specified, latest release role is imported.


<br>
Installation
-------------
Download binary from [releases](https://github.com/tinoji/galaxy-install-extended/releases), then `chmod +x`, rename and move to your PATH.

Or, build yourself like the following:
```
$ git clone https://github.com/tinoji/galaxy-install-extended
$ cd galaxy-install-extended
$ GOOS=<os> GOARCH=<arch> go build main.go
$ mv main /to/your/PATH/galaxy-install-extended
```

<br>
Usage
--------
```
$ galaxy-install-extended -h
Usage: galaxy-install-extended -r FILE [options]

Options:
  -h, --help      show this help message and exit
  -r ROLE_FILE    A file containing a list of roles to be imported

  See 'ansible-galaxy install --help' for other options
```

<br>
Example
-------
Use `examples` dir.

```
$ cd examples
$ galaxy-install-extended -r requirements.yml -p roles -f
```

<br>
Note
-----
When "latest" branch exists too, the "latest" release overrides it.

<br>
License
--------
MIT
