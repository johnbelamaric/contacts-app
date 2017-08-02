#!/bin/bash

kubectl delete ns db 2>/dev/null
while kubectl get ns 2>/dev/null db
do
	echo "$(date): Waiting for delete to finish..."
	sleep 5
done

echo "$(date): Creating database..."

kubectl apply -f manifest-db.yaml

STATUS=$(kubectl -n db get pods -o jsonpath={.items[0].status.phase})
while [ "$STATUS" != "Running" ]
do
	echo -n "$(date): Waiting for the database server container to come up ($STATUS)..."
	if [ -z "$REASON" ]; then
		echo
	else
		echo $REASON
	fi
	sleep 5
	INFO=$(kubectl -n db get pods -o jsonpath='{.items[0].status.phase},{.items[0].status.containerStatuses[0].state.waiting.reason}')
	STATUS=$(echo $INFO | cut -d , -f 1)
	REASON=$(echo $INFO | cut -d , -f 2)
done

echo $(date): Giving the database time to initialize ($STATUS)...

READY="no"
while [ "$READY" != "ready" ]
do
	READY=$(kubectl run mysql -i --rm --image=mysql --restart=Never -- mysql -hmysql.db -uroot -proot -e "select 'ready'" | grep ready | uniq)
	echo "$(date): Waiting for the database to be initialized..."
	sleep 5
done

echo "$(date): Configuring database..."

kubectl run mysql -i --rm --image=mysql --restart=Never -- mysql -hmysql.db -uroot -proot -e "create database contacts"
kubectl run mysql -i --rm --image=mysql --restart=Never -- mysql -hmysql.db -uroot -proot -e "grant all on contacts.* to 'api'@'%' identified by 'mypassword'"
