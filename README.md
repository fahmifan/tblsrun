# dockertbls

Generate database documentation from migration files. Uses [tbls](github.com/k1LoW/tbls) internally to generate the docs.

It will spawn a database, run the migrations files, and run the `tbls` to generate the db doc.

## Example Database

Here is the example database documentation [link](/example/dbdoc/README.md)

### Snippets

![db documentation snipped](doc/dbdoc-snippet.png)

## How To Use
- Set env variables, you can use `.env.example` as a reference.
- Create tbls config (see `example.tbls.yml`)
- Then run `tbslrun`, e.g.: `tblsrun postgres docker`
  - Currently it only support `postgres` with 2 modes `embedded` or `docker`
    - `embedded` will run postgres binary as child process
    - `docker` will spawn a postgres docker container, you need docker installed to use this
  - When running, it will install the latest `tbls` automatically if not exists in the `$PATH`

### Available Commands
```
Generate database documentation from migration files

Usage:
  tblsrun [command]

Available Commands:
  help        Help about any command
  postgres    Run tbls with postgres

Flags:
  -h, --help   help for tblsrun
```
