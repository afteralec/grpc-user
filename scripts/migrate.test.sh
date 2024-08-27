#!/bin/bash

# TODO: Get a path
touch ./test.db
chmod a+rw ./test.db
sqlx migrate run -D sqlite://./test.db
