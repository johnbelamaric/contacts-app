#!/bin/bash

kubectl run mysql -it --rm --image=mysql -- mysql -hmysql.db -uroot -proot <<EOF
create database contacts;
grant all on contacts.* to 'api'@'%' identified by 'mypassword';
EOF
