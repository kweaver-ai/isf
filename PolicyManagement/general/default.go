package general

import (
	"encoding/json"
	"os"

	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/kweaver-ai/GoUtils/utilities"
)

var defaultPolicies map[string]PolicyValue

func getDefaultPolicy() map[string]PolicyValue {
	if defaultPolicies == nil {
		defaultPolicies = make(map[string]PolicyValue)
		addPolicy := func(policy PolicyValue) {
			defaultPolicies[policy.Name()] = policy
		}

		addPolicy(&PasswordStrengthMeter{
			Enable: false,
			Length: 8,
		})

		addPolicy(&MultiFactorAuth{
			Enable:             false,
			ImageVcode:         false,
			PasswordErrorCount: 0,
			SMSVcode:           false,
			OTP:                false,
		})

		addPolicy(&ClientRestriction{
			PCWEB:     false,
			MobileWEB: false,
			Windows:   false,
			Mac:       false,
			Android:   false,
			IOS:       false,
			Linux:     false,
		})

		addPolicy(&UserDocumentSharing{
			Anyshare: false,
			HTTP:     true,
		})

		addPolicy(&UserDocument{
			Create: true,
			Size:   5,
		})

		addPolicy(&NetworkResitriction{
			IsEnabled: false,
		})

		addPolicy(&NoNetworkPolicyAccessor{
			IsEnabled: false,
		})

		addPolicy(&SystemProtectionLevels{
			Level: 0,
		})
	}
	return defaultPolicies
}

// CreateDefaultPolicy add policy to db if not exists
func CreateDefaultPolicy() error {
	db, err := api.ConnectDB()
	if err != nil {
		return err
	}

	var dbValue []models.Policy[[]byte]
	err = db.Find(&dbValue).Error
	if err != nil {
		return err
	}
	var existsNames []string
	for _, policy := range dbValue {
		existsNames = append(existsNames, policy.Name)
	}

	addValue := make(map[string]PolicyValue)
	for name, value := range getDefaultPolicy() {
		if !utilities.InStrSlice(name, existsNames) {
			addValue[name] = value
		}
	}

	for name, value := range addValue {
		dbType := os.Getenv("DB_TYPE")
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		var policy interface{}
		if dbType == "DM8" {
			policy = models.Policy[string]{
				Name:    name,
				Default: string(data),
				Value:   string(data),
				Locked:  false,
			}
		} else {
			policy = models.Policy[[]byte]{
				Name:    name,
				Default: data,
				Value:   data,
				Locked:  false,
			}
		}

		err = db.Create(policy).Error
		if err != nil {
			return err
		}
	}
	return nil
}
