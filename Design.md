# The idea behind pghba as a tool3

Describing the process that lead to the creation of pghba and the design
choices made.

## The problem(s)

PostgreSQL has a sort of 'Bouncer' that guards the entry to the database. It is
called 'pg\_hba.conf'. Or at least, that is the configuration file that defines
the rules for access. The pg\_hba.conf file holds a collection of lines that
are read in order from first to last and the first matching line decides the
faith of the incoming connection. So far, so good. Unfortunately, the format of
the lines is not uniform. Depending on the type of access, the line can have
two, four, fivex six or seven fields. Additionally, one has to be careful about
the ordering of the lines. There are multiple criteria to take into
consideration and if a more broadly specified line allows something that should
be forbidden, it can be difficult to detect, negating the additional security
that pg\_hba should offer. Conversely it can be difficult to figure out why
certain access is disallowed (although the verbose error message PostgreSQL
produces certainly helps.)

## The answer

The answer to this issue could be to rewrite PostgreSQL to have a better way of
providing these rules to the system, but this is somewhat complicated, to put
it mildly.

The alternative, chosen here, is to create a tool that helps with some of the
less intuitive aspects of the file itself. The pghba tool does this by making
it easier to add and delete rules in an idiomatic way instead of editing a files
by hand by providing 'shorthand' to create multiple entries simply from words
combined with 'generators' that expand in a way reminiscent of shell expansion.

## The syntax

The pghba command has two sub-commands: add and delete to add or delete line(s)
with rules to/from pg_hba.conf. Combined with the generator patterns, using a
simple CLI command can produce a lot of rules with little effort.

TODO Insert a 'man page' like description of the command line.

## What may have to change
Currently this package uses the 'net' package for the network address type but
the latest tech in this area is much better: 'net/netip' allows to more easily
convert and compare IPv4 and IPv6 addresses, including masks.

The project seems to have been cut in too many pieces and there is a mix of
types, methods and functions in some others that may have to be redistributed.

## The config file

### Format


# Development

The idea behind the development of this tool was to use what is already there.
The lead programmer and designer adopted several widely used libraries among
which is Cobra. This choice has led to the structure of the code and the
programming style which is as of writing an amalgam of procedural, functional
and object oriented.

## Architecture

The code is split up in several parts:
- The command line parser
- A single file internal package for the versioning
- The generator package creates the expansions of fields
- The pg_hba.conf file is described and instrumented in the hba package

## Data structures

The main datastructures center around the concept of a pg_hba.conf line
containing the rules and a 'Generator' that describes a pattern to use to
generate lines according to the patterns. The patterns are defined as regular
expressions and described in pkg/gnrtr/main.go.


### Why Go?
Like Python, Go is flexible in the programming styles it supports. Also like
Python, the syntax is legible, and I like Go more because it uses delimeters
instead of whitespace, but that's an opinion :). Go is easier than Rust in the
variable declaration department, but not as loose as Python. The most
important reason for Go here is building experience and using the large eco-
system/community around it.

### Why do this at all?
The main question that I personally have is: shouldn't PostgreSQL get a better
way to provide these rules? The answer should be yes, given the error-prone
nature of the current system, but in all honesty, the lines with rules are a
pretty efficient way of storing those rules and editing the file is often fine
as changes are rare in most environments. However, in more complex and modern
environments that are mostly governed and administered by software, a nicer API
may well be in place.

On a personal note, this piece of software is giving me flash-backs to ye olde
Sendmail.cf. (sorry!)