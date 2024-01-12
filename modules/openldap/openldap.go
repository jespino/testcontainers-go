package openldap

import (
	"context"
	"fmt"
	"net"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultUser     = "admin"
	defaultPassword = "adminpassword"
	defaultRoot     = "dc=example,dc=org"
	defaultAdminDn  = "cn=admin,dc=example,dc=org"
)

// OpenLDAPContainer represents the OpenLDAP container type used in the module
type OpenLDAPContainer struct {
	testcontainers.Container
}

func (c *OpenLDAPContainer) ConnectionString(ctx context.Context, args ...string) (string, error) {
	containerPort, err := c.MappedPort(ctx, "1389/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	connStr := fmt.Sprintf("ldap://%s", net.JoinHostPort(host, containerPort.Port()))
	return connStr, nil
}

// WithAdminUsername sets the initial admin username to be created when the container starts
// It is used in conjunction with WithAdminPassword to set a username and its password.
// It will create the specified user with admin power.
func WithAdminUsername(username string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Env["LDAP_ADMIN_USERNAME"] = username
	}
}

// WithAdminPassword sets the initial admin password of the user to be created when the container starts
// It is used in conjunction with WithAdminUsername to set a username and its password.
// It will set the admin password for OpenLDAP.
func WithAdminPassword(password string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Env["LDAP_ADMIN_PASSWORD"] = password
	}
}

// WithRoot sets the root of the OpenLDAP instance
func WithRoot(root string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Env["LDAP_ROOT"] = root
	}
}

// RunContainer creates an instance of the OpenLDAP container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*OpenLDAPContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "bitnami/openldap:2.6.6",
		Env: map[string]string{
			"LDAP_ADMIN_USERNAME": defaultUser,
			"LDAP_ADMIN_PASSWORD": defaultPassword,
			"LDAP_ROOT":           defaultRoot,
		},
		ExposedPorts: []string{"1389/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("** Starting slapd **"),
			wait.ForListeningPort("1389/tcp"),
		),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &OpenLDAPContainer{Container: container}, nil
}
