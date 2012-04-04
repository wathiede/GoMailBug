#!/usr/bin/env python
"""Pull all To: and CC: addresses from gzipped Mailbox files.
"""

import codecs
import email
import email.header
import email.utils
import glob
import sys

# http://code.google.com/p/python-gflags/
import gflags

__author__ = "Bill Thiede"
__email__ = "python@xinu.tv"

FLAGS = gflags.FLAGS

gflags.DEFINE_string('email_pat', None, 'Glob pattern for email to parse')
gflags.DEFINE_string('output', None, 'File to store output')

gflags.MarkFlagAsRequired('email_pat')
gflags.MarkFlagAsRequired('output')

def main(argv):
    try:
        argv = FLAGS(argv)  # parse flags
    except gflags.FlagsError, e:
        print '%s\nUsage: %s ARGS\n%s' % (e, sys.argv[0], FLAGS)
        sys.exit(1)

    with codecs.open(FLAGS.output, 'w', 'utf-8') as output:
        for fn in sorted(glob.glob(FLAGS.email_pat)):
            try:
                with open(fn, 'rb') as fp:
                    msg = email.message_from_file(fp)
                    rcpts = email.utils.getaddresses(msg.get_all('to', []) +
                                                     msg.get_all('cc', []))
                    for rcpt in rcpts:
                        name, charset = email.header.decode_header(rcpt[0])[0]
                        try:
                            output.write('%s|%s|%s\n' % (name.decode(charset or 'latin-1'), rcpt[1], fn))
                        except LookupError:
                            print 'Unknown encoding %s for %r' % (charset, name)
                            output.write('%s|%s|%s\n' % (name.decode('latin-1'), rcpt[1], fn))
            except:
                continue


if __name__ == '__main__':
    main(sys.argv)

