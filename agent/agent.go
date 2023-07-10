package agent

import (
	"fmt"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func MemberList(gcfg GlobalConfig, eps []string, options ...clientv3.OpOption) (*clientv3.MemberListResponse, error) {
	cfgSpec := clientConfigWithoutEndpoints(gcfg)
	var err error
	if len(eps) == 0 {
		eps, err = endpointsFromCmd(gcfg)
		if err != nil {
			return nil, err
		}
	}

	cfgSpec.Endpoints = eps

	c, err := createClient(cfgSpec)
	if err != nil {
		return nil, err
	}

	ctx, cancel := commandCtx(gcfg.CommandTimeout)
	defer func() {
		c.Close()
		cancel()
	}()

	members, err := c.MemberList(ctx, options...)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func EndpointStatus(gcfg GlobalConfig, ep string) (*clientv3.StatusResponse, error) {
	cfgSpec := clientConfigWithoutEndpoints(gcfg)
	cfgSpec.Endpoints = []string{ep}
	c, err := createClient(cfgSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to createClient: %w", err)
	}

	ctx, cancel := commandCtx(gcfg.CommandTimeout)
	defer func() {
		c.Close()
		cancel()
	}()
	return c.Status(ctx, ep)
}

func Read(gcfg GlobalConfig, ep string, options ...clientv3.OpOption) error {
	cfgSpec := clientConfigWithoutEndpoints(gcfg)
	cfgSpec.Endpoints = []string{ep}
	c, err := createClient(cfgSpec)
	if err != nil {
		return fmt.Errorf("failed to createClient: %w", err)
	}

	ctx, cancel := commandCtx(gcfg.CommandTimeout)
	defer func() {
		c.Close()
		cancel()
	}()

	_, err = c.Get(ctx, "health", options...)
	if err == nil || err == rpctypes.ErrPermissionDenied {
		return nil
	}
	return err
}
