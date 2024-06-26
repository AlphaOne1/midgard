package htpasswd_auth

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type HTPassWDAuth struct {
	Auths    map[string]string
	HTpassWD string
}

func (h *HTPassWDAuth) ReadFile(fileName string) error {
	file, openErr := os.Open(fileName)

	if openErr != nil {
		return openErr
	}

	scanner := bufio.NewScanner(file)
	line := 1

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")

		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			return errors.New("invalid htpasswd file, line " + strconv.Itoa(line))
		}

		line++

		h.Auths[parts[0]] = parts[1]
	}
}

func (h *HTPassWDAuth) Authorize(username, password string) (bool, error) {
	entry, entryFound := h.Auths[username]

	if !entryFound {
		return false, nil
	}

	if strings.HasPrefix(entry, "$2y$") {
		/* bcrypt */
	} else if strings.HasPrefix(entry, "$apr1$") {
		/* md5 */
	} else if strings.HasPrefix(entry, "{SHA}") {
		/* sha1 */
	} else {
		/* crypt */
	}

}
