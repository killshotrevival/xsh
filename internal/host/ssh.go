package host

import (
	"database/sql"
	"fmt"
	"xsh/internal/identity"
)

func buildSSHConnectionString(cHost *Host, dbConnection *sql.DB) (string, error) {
	var (
		cjumpHost                   *Host
		cIdentity, cJumhostIdentity *identity.Identity
		identityString              string
		jumpHostString              string
		err                         error
	)

	// Adding identity string to ssh connection only if the identity id attached is different then default identity
	if cHost.IdentityID != identity.DefaultIdentityID {
		cIdentity, err = identity.GetIdentityByID(dbConnection, cHost.IdentityID)
		if err != nil {
			return "", err
		}
		identityString = fmt.Sprintf("-i %s", cIdentity.Path)
	}

	if cHost.JumphostID.Valid {
		cjumpHost, err = GetHostByID(dbConnection, cHost.JumphostID.UUID.String())
		if err != nil {
			return "", err
		}

		proxyIdentityString := ""

		if cjumpHost.IdentityID != identity.DefaultIdentityID {
			cJumhostIdentity, err = identity.GetIdentityByID(dbConnection, cjumpHost.IdentityID)
			if err != nil {
				return "", err
			}
			proxyIdentityString = fmt.Sprintf("-i %s", cJumhostIdentity.Path)

		}

		jumpHostString = fmt.Sprintf(`-o ProxyCommand="ssh %s -W %s:%d %s@%s -p %d"`,
			proxyIdentityString,

			cHost.Address,
			cHost.Port,

			cjumpHost.User,
			cjumpHost.Address,
			cjumpHost.Port)

	}

	connectionString := fmt.Sprintf("ssh -p %d %s %s %s@%s",
		cHost.Port,
		identityString,
		jumpHostString,
		cHost.User,
		cHost.Address,
	)

	return connectionString, nil
}
