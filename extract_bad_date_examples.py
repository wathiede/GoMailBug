#!/usr/bin/env python
"""One-line description

In-depth description
"""

import csv
import pprint
import os.path
import sys

__author__ = "Bill Thiede"
__copyright__ = "Copyright 2011"
__email__ = "python@xinu.tv"


def main(args):
    mailers = {}
    fp = open(args[0])
    csv_fp = csv.reader(fp)
    for row in csv_fp:
        if len(row) != 2:
            continue

        mailer, date = row
        if 'Maildir' in mailer:
            # Skip unknown mailers
            continue
        mailers[mailer] = date
    fp.close()

    fp = open(args[1], 'wb')
    csv_fp = csv.writer(fp)
    csv_fp.writerows(sorted(mailers.items()))
    fp.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print "Usage: %s <input.csv> <output.csv>" % os.path.basename(sys.argv[0])
        raise SystemExit

    main(sys.argv[1:])

