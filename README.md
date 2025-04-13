# passutils
utilities for pass password manager

## exporter

**Exports a password store to a user-specified output directory**

build like

```shell
make exporter
```

run like

```shell
./bin/exporter --outdir "somewhere"
```

to create export file `$HOME/somewhere/password-export.json`