package campaign

const (
	QueryGetAll = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			goal_amount,
			current_amount,
			is_exclusive,
			donor_count,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
	`

	QueryGetCampaignByID = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			goal_amount,
			current_amount,
			is_exclusive,
			donor_count,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetAllImage = `
		SELECT
			id,
			campaign_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_images
		WHERE
			deleted_at IS NULL
	`

	QueryGetCampaignImageByID = `
		SELECT
			id,
			campaign_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_images
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetCampaignImages = `
		SELECT
			id,
			campaign_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_images
		WHERE
			deleted_at IS NULL
		AND
			campaign_id = ?
		ORDER BY
			is_primary DESC
	`

	QueryGetAllCategory = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_categories
		WHERE
			deleted_at IS NULL
	`

	QueryGetCampaignCategoryByID = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_categories
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetAllExclusiveCampaign = `
		SELECT
			id,
			campaign_id,
			winner_user_id,
			is_reward_money,
			reward,
			is_paid_off,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			exclusive_campaigns
		WHERE
			deleted_at IS NULL
	`

	QueryGetCampaignExclusiveByID = `
		SELECT
			id,
			campaign_id,
			winner_user_id,
			is_reward_money,
			reward,
			is_paid_off,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			exclusive_campaigns
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetCampaignExclusiveByCampaignID = `
		SELECT
			id,
			campaign_id,
			winner_user_id,
			is_reward_money,
			reward,
			is_paid_off,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			exclusive_campaigns
		WHERE
			deleted_at IS NULL
		AND
			campaign_id = ?
		ORDER BY
			id ASC
		LIMIT
			1
	`

	QueryGetCampaignExclusiveByWinnerUserID = `
		SELECT
			id,
			campaign_id,
			winner_user_id,
			is_reward_money,
			reward,
			is_paid_off,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			exclusive_campaigns
		WHERE
			deleted_at IS NULL
		AND
			winner_user_id = ?
		ORDER BY
			id ASC
	`

	QueryAdminDataTablesCampaigns = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			goal_amount,
			current_amount,
			is_exclusive,
			donor_count,
			(SELECT COUNT(id) FROM campaign_images WHERE deleted_at IS NULL AND campaign_id = campaigns.id) AS total_image,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesCampaigns = `
		SELECT
			COUNT(id) AS count_id
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
	`

	QueryAdminDataTablesCategories = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaign_categories
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesCategories = `
		SELECT
			COUNT(id) AS count_id
		FROM
			campaign_categories
		WHERE
			deleted_at IS NULL
	`

	QueryUserDataTablesCampaigns = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			goal_amount,
			current_amount,
			is_exclusive,
			donor_count,
			(SELECT COUNT(id) FROM campaign_images WHERE deleted_at IS NULL AND campaign_id = campaigns.id) AS total_image,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryCountAllUserDataTablesCampaigns = `
		SELECT
			COUNT(id) AS count_id
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryGetTotalDonation = QueryCountAllAdminDataTablesCampaigns

	QueryGetDonationCompleted = `
		SELECT
			COUNT(id) AS count_id
		FROM
			campaigns
		WHERE
			deleted_at IS NULL
		AND
			status = 'finished'
	`

	QueryAdminDataTablesWinnersExclusiveCampaigns = `
		SELECT
			id,
			campaign_id,
			COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') AS campaign_title,
			winner_user_id,
			COALESCE((SELECT name FROM users WHERE id = winner_user_id), '') AS winner_user_name,
			is_reward_money,
			reward,
			is_paid_off
		FROM
			exclusive_campaigns
		WHERE
			winner_user_id != 0
		AND
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesWinnersExclusiveCampaigns = `
		SELECT
			COUNT(id) AS count_id
		FROM
			exclusive_campaigns
		WHERE
			winner_user_id != 0
		AND
			deleted_at IS NULL
	`

	QueryGetOneRandomUserIDForWinnerExclusiveCampaign = `
		SELECT
			user_id
		FROM
			transactions
		WHERE
			user_id <> 0
		AND
			user_id IS NOT NULL
		AND
			user_id IN (
				SELECT
					id
				FROM
					users
				WHERE
					deleted_at IS NULL
			)
		AND
			status = 'paid'
		AND
			campaign_id = ?
		ORDER BY
			RAND()
		LIMIT
			1
	`
)
