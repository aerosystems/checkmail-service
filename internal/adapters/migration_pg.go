package adapters

import (
	"fmt"
	"gorm.io/gorm"
)

func AutoMigrateGORM(db *gorm.DB) error {
	if err := db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'domain_type') THEN
				CREATE TYPE domain_type AS ENUM ('unknown', 'blacklist', 'whitelist');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'match_type') THEN
				CREATE TYPE match_type AS ENUM ('prefix', 'suffix', 'contains', 'equals');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_type') THEN
				CREATE TYPE subscription_type AS ENUM ('unknown', 'trial', 'startup', 'business');
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create enum types: %v", err)
	}

	if err := db.AutoMigrate(&Domain{}, &Filter{}, &Review{}, &Access{}); err != nil {
		return fmt.Errorf("failed to AutoMigrateGORM: %v", err)
	}
	return nil
}
