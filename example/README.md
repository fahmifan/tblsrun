# Example

In this example we configure 2 schemas, `bar` & `foo`. Each schema have different migration files. We want to generate the documentation for each schema into different folder.

The migrations files are located in `example/migrations` folder, and the generated documentation will be located in `example/dbdoc` folder.

We create two `.env` files, each one will specify what schema it will generate the documentation for, and where the migration files are located.
The schema and migrations is separated by comma (`,`). And each schema & migrations must be in the same order.

Then we create two `.tbls.yml` files, each one will specify the output folder for the documentation and to filter which schemas to include in the documentation using `include` option.

Finally we run `tblsrun` with the `postgres` driver and `docker` mode, check the Makefile to understand how it runs.
