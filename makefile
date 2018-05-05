PREFIX="/usr/local/bin/"
NAME="base58"
all:
	CGO_ENABLED=0 go build -ldflags='-w -s' -o ${NAME}
	strip ${NAME}

install:
	install ${NAME} ${PREFIX}

package:
	gzip -k ${NAME}
