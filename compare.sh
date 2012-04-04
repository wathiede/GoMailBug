#!/bin/sh
#
# compare.sh '/pat/for/mail/messages/*'

# No unset vars
set -u

FN=/tmp/email-compare-$$
PYFN=${FN}-py
GOFN=${FN}-go

echo "Running python extraction to ${PYFN}"
./extract.py --email_pat=$1 --output=${PYFN}
echo "Running go extraction to ${GOFN}"
go run extract.go -email_pat=$1 -output=${GOFN}

echo "Comparing"
diff -u ${GOFN} ${PYFN} | less
#rm -f ${PYFN} ${GOFN}
