# pgmgr - Manage pg_hba files with more sophistication

# Introduction

pghba is a tool meant to manage pg_hba files with a bit more sphistication.
pg_hba.conf files are usually created by initdb, and managed with vi.
Alternatively, they are generated by puppet / Ansible with a template option.
Both options work, but also are error prone and require experience and manual labor.
This repo delivers admins with a tool to manage pg_hba files with a bit more sophistication.
Alternatively, this module can also be imported in other GoLang projects that what to read, manipulate and write pg_hba.conf files internally.

## The origin
Somewhere in 2015, a DBA of the Dutch Government built a python script to manage pg_hba.conf files with a bit more sophistication.
The script could read and parse a pg_hba file in pg_hba rules, add or remove phba rules and write out the result.
The result was ordered in a fashion that made sense (larger scopes, like larger networks, all users, all database, etc. ordered later in the file).
And some common inconveniences where cleaned, like:
* Networks where the network address was not the first addres in the pool (e.a. 192.168.0.192/24, should actually be 192.168.0.0/24)
* Multiple rows with the same identifiers (e.a. 'host all all 127.0.0.1/32 trust', and 'host all all 127.0.0.1/32 md5', first row wins).
* Convert ip/netmask addresses in ip/CIDR notation (whenever possible)
* Add a CIDR /32 for IP addresses without CIDR notation
* Reformat ill formulated P addresses (like 192.168.00.001, to be formulated as 192.168.0.1)...
* Capture and error on eroneous ip adresses (like 192.1168.0.1)
* Multiple ways to format, like spacing the columns where all brought back to one format.

The script was idempotent, and very intuitive.
In the end, the script was changed into [Ansible's postgresql_pg_hba module](https://docs.ansible.com/ansible/2.9/modules/postgresql_pg_hba_module.html) and upstreamed to the project.

Somewhere in 2018, the ambition grew to write the similar functionality as a GoLang module, and now that ambition has grown into an actual module / tool.

# Downloading
The most straight forward option is to download the [pghba](https://github.com/MannemSolutions/pghba) binary directly from the [github release page](https://github.com/MannemSolutions/pghba/releases).
The other option would be to build from source (if you feel you must)

# Using
**Note** please refer to our [commandline arguments](CLI.md) to learn about all capabilities and configuration features of pghba.

After downloading the binary to a folder in your path, you can run pghba with a command like:
```bash
pghba add -h localhost -d all -U all -m md5 -f $PGDATA/pg_hba.conf 
```
to allow access with md5 authetication for all tcp connections from localhost for connections on any database by any user in the $PGDATA/pg_hba.conf.
