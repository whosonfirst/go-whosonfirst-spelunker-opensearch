#!bin/sh

# Utitily tool / reference implementation for importing and then indexin
# multiple WOF repos in to a local OpenSearch index. This is not a general
# purpose tool so you may need to adjust details to suit your specific
# needs.

# https://github.com/whosonfirst/go-whosonfirst-github
LIST_REPOS=/usr/local/whosonfirst/go-whosonfirst-github/bin/wof-list-repos

# https://github.com/whosonfirst/go-whosonfirst-opensearch
OPENSEARCH_INDEX=/usr/local/whosonfirst/go-whosonfirst-opensearch/bin/wof-opensearch-index

# Maybe make this a CLI option...
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-a -prefix whosonfirst-data-postalcode-b -prefix whosonfirst-data-postalcode-c"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-d -prefix whosonfirst-data-postalcode-e -prefix whosonfirst-data-postalcode-f"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-g -prefix whosonfirst-data-postalcode-h -prefix whosonfirst-data-postalcode-i"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-j -prefix whosonfirst-data-postalcode-k -prefix whosonfirst-data-postalcode-l"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-m -prefix whosonfirst-data-postalcode-n -prefix whosonfirst-data-postalcode-o"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-p -prefix whosonfirst-data-postalcode-q -prefix whosonfirst-data-postalcode-r"
# LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-s -prefix whosonfirst-data-postalcode-t -prefix whosonfirst-data-postalcode-u"
LIST_REPOS_PREFIX="-prefix whosonfirst-data-postalcode-v -prefix whosonfirst-data-postalcode-w -prefix whosonfirst-data-postalcode-x -prefix whosonfirst-data-postalcode-y -prefix whosonfirst-data-postalcode-z"

# This is for debugging. Do not change this at your own risk.
# (That means you should change this.)
OPENSEARCH_PASSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

OPENSEARCH_WRITER_URI=constant://?val=opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Frequire-tls%3Dtrue%26insecure%3Dtrue%26debug%3Dfalse%26username%3Dadmin%26password%3D${OPENSEARCH_PASSWD}

# Pull straight from GitHub, write to /tmp, remove after indexing
ITERATOR_URI=git:///tmp

REPOS=`${LIST_REPOS} ${LIST_REPOS_PREFIX} -exclude whosonfirst-data-admin-alt`

for REPO in ${REPOS}
do
    echo "Start indexing ${REPO} (https://github.com/whosonfirst-data/${REPO}.git)"
    ${OPENSEARCH_INDEX} -writer-uri ${OPENSEARCH_WRITER_URI} -iterator-uri ${ITERATOR_URI} https://github.com/whosonfirst-data/${REPO}.git
    echo "Finished indexing ${REPO}"    
done
