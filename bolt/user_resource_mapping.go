package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	bolt "github.com/coreos/bbolt"
	platform "github.com/influxdata/influxdb"
)

var (
	userResourceMappingBucket = []byte("userresourcemappingsv1")
)

func (c *Client) initializeUserResourceMappings(ctx context.Context, tx *bolt.Tx) error {
	if _, err := tx.CreateBucketIfNotExists([]byte(userResourceMappingBucket)); err != nil {
		return err
	}
	return nil
}

func filterMappingsFn(filter platform.UserResourceMappingFilter) func(m *platform.UserResourceMapping) bool {
	return func(mapping *platform.UserResourceMapping) bool {
		return (!filter.UserID.Valid() || (filter.UserID == mapping.UserID)) &&
			(!filter.ResourceID.Valid() || (filter.ResourceID == mapping.ResourceID)) &&
			(filter.UserType == "" || (filter.UserType == mapping.UserType)) &&
			(filter.ResourceType == "" || (filter.ResourceType == mapping.ResourceType))
	}
}

// FindUserResourceMappings returns a list of UserResourceMappings that match filter and the total count of matching mappings.
func (c *Client) FindUserResourceMappings(ctx context.Context, filter platform.UserResourceMappingFilter, opt ...platform.FindOptions) ([]*platform.UserResourceMapping, int, error) {
	ms := []*platform.UserResourceMapping{}
	err := c.db.View(func(tx *bolt.Tx) error {
		mappings, err := c.findUserResourceMappings(ctx, tx, filter)
		if err != nil {
			return err
		}
		ms = mappings
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return ms, len(ms), nil
}

func (c *Client) findUserResourceMappings(ctx context.Context, tx *bolt.Tx, filter platform.UserResourceMappingFilter) ([]*platform.UserResourceMapping, error) {
	ms := []*platform.UserResourceMapping{}
	filterFn := filterMappingsFn(filter)
	err := c.forEachUserResourceMapping(ctx, tx, func(m *platform.UserResourceMapping) bool {
		if filterFn(m) {
			ms = append(ms, m)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (c *Client) findUserResourceMapping(ctx context.Context, tx *bolt.Tx, filter platform.UserResourceMappingFilter) (*platform.UserResourceMapping, error) {
	ms, err := c.findUserResourceMappings(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	if len(ms) == 0 {
		return nil, fmt.Errorf("userResource mapping not found")
	}

	return ms[0], nil
}

func (c *Client) CreateUserResourceMapping(ctx context.Context, m *platform.UserResourceMapping) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		if err := c.createUserResourceMapping(ctx, tx, m); err != nil {
			return err
		}

		if m.ResourceType == platform.OrgsResourceType {
			return c.createOrgDependentMappings(ctx, tx, m)
		}

		return nil
	})
}

func (c *Client) createUserResourceMapping(ctx context.Context, tx *bolt.Tx, m *platform.UserResourceMapping) error {
	unique := c.uniqueUserResourceMapping(ctx, tx, m)

	if !unique {
		return fmt.Errorf("mapping for user %s already exists", m.UserID.String())
	}

	v, err := json.Marshal(m)
	if err != nil {
		return err
	}

	key, err := userResourceKey(m)
	if err != nil {
		return err
	}

	if err := tx.Bucket(userResourceMappingBucket).Put(key, v); err != nil {
		return err
	}

	return nil
}

// This method creates the user/resource mappings for resources that belong to an organization.
func (c *Client) createOrgDependentMappings(ctx context.Context, tx *bolt.Tx, m *platform.UserResourceMapping) error {
	bf := platform.BucketFilter{OrganizationID: &m.ResourceID}
	bs, err := c.findBuckets(ctx, tx, bf)
	if err != nil {
		return err
	}
	for _, b := range bs {
		m := &platform.UserResourceMapping{
			ResourceType: platform.BucketsResourceType,
			ResourceID:   b.ID,
			UserType:     m.UserType,
			UserID:       m.UserID,
		}
		if err := c.createUserResourceMapping(ctx, tx, m); err != nil {
			return err
		}
		// TODO(desa): add support for all other resource types.
	}

	return nil
}

func userResourceKey(m *platform.UserResourceMapping) ([]byte, error) {
	encodedResourceID, err := m.ResourceID.Encode()
	if err != nil {
		return nil, err
	}

	encodedUserID, err := m.UserID.Encode()
	if err != nil {
		return nil, err
	}

	key := make([]byte, len(encodedResourceID)+len(encodedUserID))
	copy(key, encodedResourceID)
	copy(key[len(encodedResourceID):], encodedUserID)

	return key, nil
}

func (c *Client) forEachUserResourceMapping(ctx context.Context, tx *bolt.Tx, fn func(*platform.UserResourceMapping) bool) error {
	cur := tx.Bucket(userResourceMappingBucket).Cursor()
	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		m := &platform.UserResourceMapping{}
		if err := json.Unmarshal(v, m); err != nil {
			return err
		}
		if !fn(m) {
			break
		}
	}

	return nil
}

func (c *Client) uniqueUserResourceMapping(ctx context.Context, tx *bolt.Tx, m *platform.UserResourceMapping) bool {
	key, err := userResourceKey(m)
	if err != nil {
		return false
	}

	v := tx.Bucket(userResourceMappingBucket).Get(key)
	return len(v) == 0
}

// DeleteUserResourceMapping deletes a user resource mapping.
func (c *Client) DeleteUserResourceMapping(ctx context.Context, resourceID platform.ID, userID platform.ID) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		m, err := c.findUserResourceMapping(ctx, tx, platform.UserResourceMappingFilter{
			ResourceID: resourceID,
			UserID:     userID,
		})
		if err != nil {
			return err
		}

		if err := c.deleteUserResourceMapping(ctx, tx, platform.UserResourceMappingFilter{
			ResourceID: resourceID,
			UserID:     userID,
		}); err != nil {
			return err
		}

		if m.ResourceType == platform.OrgsResourceType {
			return c.deleteOrgDependentMappings(ctx, tx, m)
		}

		return nil
	})
}

func (c *Client) deleteUserResourceMapping(ctx context.Context, tx *bolt.Tx, filter platform.UserResourceMappingFilter) error {
	ms, err := c.findUserResourceMappings(ctx, tx, filter)
	if err != nil {
		return err
	}
	if len(ms) == 0 {
		return fmt.Errorf("userResource mapping not found")
	}

	key, err := userResourceKey(ms[0])
	if err != nil {
		return err
	}

	return tx.Bucket(userResourceMappingBucket).Delete(key)
}

func (c *Client) deleteUserResourceMappings(ctx context.Context, tx *bolt.Tx, filter platform.UserResourceMappingFilter) error {
	ms, err := c.findUserResourceMappings(ctx, tx, filter)
	if err != nil {
		return err
	}
	for _, m := range ms {
		key, err := userResourceKey(m)
		if err != nil {
			return err
		}

		if err = tx.Bucket(userResourceMappingBucket).Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// This method deletes the user/resource mappings for resources that belong to an organization.
func (c *Client) deleteOrgDependentMappings(ctx context.Context, tx *bolt.Tx, m *platform.UserResourceMapping) error {
	bf := platform.BucketFilter{OrganizationID: &m.ResourceID}
	bs, err := c.findBuckets(ctx, tx, bf)
	if err != nil {
		return err
	}
	for _, b := range bs {
		if err := c.deleteUserResourceMapping(ctx, tx, platform.UserResourceMappingFilter{
			ResourceType: platform.BucketsResourceType,
			ResourceID:   b.ID,
			UserID:       m.UserID,
		}); err != nil {
			return err
		}
		// TODO(desa): add support for all other resource types.
	}

	return nil
}
