yarastore
==========

``yarastore`` is a command-line tool for managing and compiling YARA rulesets. It allows users to
collect, organize, and compile YARA rules from multiple sources into a single ruleset. The tool
supports recursive file scanning, pattern matching, and exclusion rules to make the compilation
and scanning more flexible and efficient.

Features
========

* Compiling YARA rules from files and directories
* Exclude specific files and directories from being processed
* Include and Exclude globbing patterns
* Recursive file scanning
* Configurable options via TOML configuration file and command-line arguments
* JSON report output for matched rules

Installation
============

Requirements
------------

Make sure you have ``libyara`` installed to build ``go-yara``. 
Check out `go-yara <https://github.com/hillu/go-yara?tab=readme-ov-file#buildinstallation>`_ for more information.


Build from Source
-----------------

Clone the repository and build the binary using the following commands:


.. code-block:: bash

    git clone https://github.com/RyugaXhypeR/yarastore
    cd yarastore


Build the binary:

.. code-block:: bash

    go build -o yarastore


Usage
=====

Command-line Options
--------------------

``yarastore`` provides two subcommands: ``compile`` and ``match`` to compile and match YARA rules, respectively,
with an optional configuration file which can be used to specify the ruleset and other options.

To see the available options for each subcommand, run:

.. code-block:: bash

    ./yarastore compile --help
    ./yarastore match --help


Configuration
-------------

You can specify options for configuration via a TOML configuration file. The configuration file can be passed
as a command-line argument to the ``yarastore`` tool. The configuration file should have the following format:

.. code-block:: toml

    [rules]
    dirs = ["/path/to/rules/dir1", "/path/to/rules/dir2"]
    files = ["/path/to/rules/file1.yar", "/path/to/rules/file2.yar"]
    exclude = ["dir1/", "file1.yar"]
    include_pattern = "*.yar"
    exclude_pattern = "_*"
    recursive = true

    [target]
    dirs = ["/path/to/target/dir1", "/path/to/target/dir2"]
    files = ["/path/to/target/file1", "/path/to/target/file2"]
    exclude = ["dir1/", "file1"]
    include_pattern = "*.exe"
    exclude_pattern = "_*"
    recursive = true


The configuration file can be passed as a command-line argument using the ``--config`` option:

.. code-block:: bash

    ./yarastore compile --config /path/to/config.toml


These configurations can also be modified using command-line arguments, the command-line arguments will take precedence.


Example Usage
-------------

.. code-block:: bash

    ./yarastore compile --dirs "rules1/ rules2/" -r -o rules.yar
    ./yarastore match rules.yar --dirs "target1/ target2/" -r -o report.json


Future Plans
============

* Add support to download YARA rules from the internet
* Add support to scan compressed files
* Better scan reports
