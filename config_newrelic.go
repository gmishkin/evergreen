package evergreen

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewRelicConfig struct {
	AccountID     string `bson:"account_id" json:"account_id" yaml:"account_id"`
	TrustKey      string `bson:"trust_key" json:"trust_key" yaml:"trust_key"`
	AgentID       string `bson:"agent_id" json:"agent_id" yaml:"agent_id"`
	LicenseKey    string `bson:"license_key" json:"license_key" yaml:"license_key"`
	ApplicationID string `bson:"application_id" json:"application_id" yaml:"application_id"`
}

func (c *NewRelicConfig) SectionId() string { return "newrelic" }

func (c *NewRelicConfig) Get(ctx context.Context) error {
	res := GetEnvironment().DB().Collection(ConfigCollection).FindOne(ctx, byId(c.SectionId()))
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			*c = NewRelicConfig{}
			return nil
		}
		return errors.Wrapf(err, "getting config section '%s'", c.SectionId())
	}

	if err := res.Decode(&c); err != nil {
		return errors.Wrapf(err, "decoding config section '%s'", c.SectionId())
	}

	return nil
}

func (c *NewRelicConfig) Set(ctx context.Context) error {
	_, err := GetEnvironment().DB().Collection(ConfigCollection).UpdateOne(ctx, byId(c.SectionId()), bson.M{
		"$set": bson.M{
			"account_id":     c.AccountID,
			"trust_key":      c.TrustKey,
			"agent_id":       c.AgentID,
			"license_key":    c.LicenseKey,
			"application_id": c.ApplicationID,
		},
	}, options.Update().SetUpsert(true))

	return errors.Wrapf(err, "updating config section '%s'", c.SectionId())
}

func (c *NewRelicConfig) ValidateAndDefault() error {
	allFieldsAreEmpty := c.AccountID == "" && c.TrustKey == "" && c.AgentID == "" && c.LicenseKey == "" && c.ApplicationID == ""
	allFieldsAreFilledOut := len(c.AccountID) > 0 && len(c.TrustKey) > 0 && len(c.AgentID) > 0 && len(c.LicenseKey) > 0 && len(c.ApplicationID) > 0

	if !allFieldsAreEmpty && !allFieldsAreFilledOut {
		return errors.New("must provide all fields or no fields for New Relic settings")
	}
	return nil
}
