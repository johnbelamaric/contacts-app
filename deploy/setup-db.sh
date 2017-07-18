#!/bin/bash

kubectl delete ns db 2>/dev/null
while kubectl get ns 2>/dev/null db
do
	echo "Waiting for delete to finish..."
	sleep 5
done

echo "Creating database..."

kubectl apply -f manifest-db.yaml

STATUS=$(kubectl -n db get pods -o jsonpath={.items[0].status.phase})
while [ "$STATUS" != "Running" ]
do
	echo "Waiting for the database server container to come up ($STATUS)..."
	sleep 5
	STATUS=$(kubectl -n db get pods -o jsonpath={.items[0].status.phase})
done

echo "Giving the database time to initialize ($STATUS)..."

sleep 10

echo "Configuring database..."

kubectl run mysql -i --rm --image=mysql --restart=Never -- mysql -hmysql.db -uroot -proot -e "create database contacts"
kubectl run mysql -i --rm --image=mysql --restart=Never -- mysql -hmysql.db -uroot -proot -e "grant all on contacts.* to 'api'@'%' identified by 'mypassword'"
