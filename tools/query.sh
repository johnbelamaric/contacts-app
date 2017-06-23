kubectl -n contacts run -it --rm --image=infoblox/dnstools contacts-query --command -- curl http://api/contacts
