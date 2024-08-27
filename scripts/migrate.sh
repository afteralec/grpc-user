#!/bin/bash

touch /var/db/user.db
chgrp -R noroot /var/db
chmod g+rw /var/db
chmod g+rw /var/db/user.db
sqlx migrate run -D sqlite:///var/db/user.db
