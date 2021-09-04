# All commandline options

- Optionally convert dns to ip
- Convert ipv4 addresses with a CIDR, with a netmask, 
- Use ipv4 addresses with a CIDR, with a netmask, 
- Expansion patterns e.a.:
  - 'server-{0..2}' -> 'server-0', 'server-1', 'server-2')
  - 'server-{a,b,c,d}' -> 'server-a', 'server-b', 'server-c', 'server-d',
- Query postgres using a regexp and using the results as the list of databases, or users
- Read errors from postgresql.log and generate the required rules
- Read pg_hba_file_rules table and generate a pg_hba.conf file from it
