# Noms

<img src='doc/nommy_cropped_smaller.png' width='350' title='Nommy, the snacky otter'>

[Use Cases](#use-cases)&nbsp; | &nbsp;[Setup](#setup)&nbsp; | &nbsp;[TUI](#terminal-ui)&nbsp; | &nbsp;[Status](#status)&nbsp; | &nbsp;[Documentation](./doc/intro.md)&nbsp; | &nbsp;[Contact](#contact-us)
<br><br>

[![Build Status](https://travis-ci.org/attic-labs/noms.svg?branch=master)](https://travis-ci.org/attic-labs/noms)
[![Docker Build Status](https://img.shields.io/docker/build/noms/noms.svg)](https://hub.docker.com/r/noms/noms/)
[![GoDoc](https://godoc.org/github.com/attic-labs/noms?status.svg)](https://godoc.org/github.com/attic-labs/noms)

# Welcome

*Noms* is a decentralized database philosophically descendant from the Git version control system.

Like Git, Noms is:

* **Versioned:** By default, all previous versions of the database are retained. You can trivially track how the database evolved to its current state, easily and efficiently compare any two versions, or even rewind and branch from any previous version.
* **Synchronizable:** Instances of a single Noms database can be disconnected from each other for any amount of time, then later reconcile their changes efficiently and correctly.

Unlike Git, Noms is a database, so it also:

* Primarily **stores structured data**, not files and directories (see: [the Noms type system](https://github.com/attic-labs/noms/blob/master/doc/intro.md#types))
* **Scales well** to large amounts of data and concurrent clients
* Supports **atomic transactions** (a single instance of Noms is CP, but Noms is typically run in production backed by S3, in which case it is "[effectively CA](https://cloud.google.com/spanner/docs/whitepapers/SpannerAndCap.pdf)")
* Supports **efficient indexes** (see: [Noms prolly-trees](https://github.com/attic-labs/noms/blob/master/doc/intro.md#prolly-trees-probabilistic-b-trees))
* Features a **flexible query model** (see: [GraphQL](./go/ngql/README.md))

A Noms database can reside within a file system or in the cloud:

* The (built-in) [NBS](./go/nbs) `ChunkStore` implementation provides two back-ends which provide persistence for Noms databases: one for storage in a file system and one for storage in an S3 bucket.

Finally, because Noms is content-addressed, it yields a very pleasant programming model.

Working with Noms is ***declarative***. You don't `INSERT` new data, `UPDATE` existing data, or `DELETE` old data. You simply *declare* what the data ought to be right now. If you commit the same data twice, it will be deduplicated because of content-addressing. If you commit _almost_ the same data, only the part that is different will be written.

<br>

## Use Cases

#### [Decentralization](./doc/decent/about.md)

Because Noms is very good at sync, it makes a decent basis for rich, collaborative, fully-decentralized applications.

#### Mobile Offline-First Database

Embed Noms into mobile applications, making it easier to build offline-first, fully synchronizing mobile applications.

<br>

## Install

### From Source

1. Clone the repository:
```bash
git clone https://github.com/fully-autonomous/noms.git
cd noms
```

2. Build the CLI:
```bash
go build -o bin/noms ./cmd/noms/
```

3. Add to your PATH:
```bash
export PATH="$PWD/bin:$PATH"
```

4. Verify installation:
```bash
noms version
```

<br>

## Terminal UI

Noms now includes a fully interactive Terminal UI (TUI) for intuitive database management.

### Features

- **Dashboard** - Real-time overview of status, branches, datasets, and recent commits
- **Branch Manager** - Create, checkout, and delete branches with visual selection
- **Dataset Browser** - Browse and view dataset contents interactively  
- **Sync Panel** - Push, pull, and manage remotes with interactive input
- **Commit Workflow** - Write commit messages with full interactive prompts

### Running the TUI

```bash
cd tui
bun install
bun run src/index.ts
```

Or run directly:
```bash
cd tui && bun run src/index.ts
```

### Navigation

- `[1-5]` - Navigate between Dashboard, Branches, Datasets, Sync, and Commit screens
- `[b/q]` - Go back / Quit
- `[j/k]` or `[↑/↓]` - Navigate lists
- `[Enter]` - Select/Confirm
- `[n/c/d]` - New branch, Checkout, Delete branch (in Branch manager)
- `[p/l/a]` - Push, Pull, Add remote (in Sync panel)

<br>

## Quick Start

Initialize a new database:
```bash
mkdir mydb && cd mydb
noms init
```

Create and switch branches:
```bash
noms branch -c feature-branch
noms checkout feature-branch
```

Commit data:
```bash
noms commit -m "Add initial data" /path/to/data mydataset
```

View history:
```bash
noms log
noms show mydataset
```

Sync with remote:
```bash
noms remote --add origin https://remote-server.com/db
noms push origin
```

<br>

## Status

This fork is actively maintained with the following completed features:

### Completed

* ✅ **Terminal UI** - Full interactive TUI with menus, navigation, and all operations
* ✅ **Branch Management** - Create, checkout, delete, and list branches
* ✅ **Remote Configuration** - Persistent remote configuration with add/remove/list
* ✅ **Sync Operations** - Push and pull with remote servers
* ✅ **Time Travel** - Browse and checkout historical versions
* ✅ **Query System** - SQL-like query interface via `noms query`
* ✅ **Dataset Management** - Import, export, and manage structured datasets

### In Progress

* Garbage Collection for orphaned chunks
* Migration tools for format upgrades
* Query language improvements

<br>

## Learn More About Noms

For the decentralized web: [The Decentralized Database](doc/decent/about.md)

Learn the basics: [Technical Overview](doc/intro.md)

Tour the CLI: [Command-Line Interface Tour](doc/cli-tour.md)

Tour the Go API: [Go SDK Tour](doc/go-tour.md)

<br>

## Contact Us

Interested in using Noms? Awesome! We would be happy to work with you to help understand whether Noms is a fit for your problem. Reach out at:

- [GitHub Issues](https://github.com/fully-autonomous/noms/issues)
- Original Noms discussion: [Mailing List](https://groups.google.com/forum/#!forum/nomsdb)
- Original Noms Twitter: [Twitter](https://twitter.com/nomsdb)

## Licensing

Noms is open source software, licensed by Attic Labs, Inc. under the Apache License, Version 2.0.
