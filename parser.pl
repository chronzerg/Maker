#!/usr/bin/perl
use warnings;
use strict;

while(<STDIN>) {
    # Remove comments
    s/(?<!\\)#.*\n/\n/g;

    # Remove newlines not
    # preceeded by a semi
    # colon
    # TODO: Count the number
    # of replaced lines and
    # add blanks afterwards
    # to preserve line numbers
    # of other commands.
    s/(?<!;)\s*\n//g;

    # Remove semicolons
    s/;\s*\n/\n/g;

    print;
}
