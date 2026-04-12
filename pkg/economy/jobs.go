package economy

import (
	"fmt"
	"time"
)

// JobType represents the type of job
type JobType int

const (
	JobMiner JobType = iota
	JobBuilder
	JobHunter
	JobFarmer
	JobExplorer
	JobFisher
	JobCraftsman
	JobMerchant
	JobWarrior
	JobEnchanter
)

// String returns job name
func (j JobType) String() string {
	switch j {
	case JobMiner:
		return "Miner"
	case JobBuilder:
		return "Builder"
	case JobHunter:
		return "Hunter"
	case JobFarmer:
		return "Farmer"
	case JobExplorer:
		return "Explorer"
	case JobFisher:
		return "Fisher"
	case JobCraftsman:
		return "Craftsman"
	case JobMerchant:
		return "Merchant"
	case JobWarrior:
		return "Warrior"
	case JobEnchanter:
		return "Enchanter"
	}
	return "Unknown"
}

// JobDescription returns job description
func (j JobType) Description() string {
	descriptions := map[JobType]string{
		JobMiner:      "Mine ores and stone to earn money",
		JobBuilder:    "Place blocks and build structures",
		JobHunter:     "Kill mobs and creatures",
		JobFarmer:     "Grow crops and harvest resources",
		JobExplorer:   "Discover new areas and biomes",
		JobFisher:     "Catch fish and treasures",
		JobCraftsman:  "Craft items and equipment",
		JobMerchant:   "Trade with players and NPCs",
		JobWarrior:    "Win PvP battles and duels",
		JobEnchanter:  "Enchant items and use magic",
	}
	
	if desc, exists := descriptions[j]; exists {
		return desc
	}
	return "Unknown job"
}

// JobStats tracks job-specific statistics
type JobStats struct {
	BlocksMined      int64   `json:"blocks_mined"`
	BlocksPlaced     int64   `json:"blocks_placed"`
	MobsKilled       int64   `json:"mobs_killed"`
	DistanceTraveled float64 `json:"distance_traveled"`
	ItemsCrafted     int64   `json:"items_crafted"`
	FishCaught       int64   `json:"fish_caught"`
	TradesCompleted  int64   `json:"trades_completed"`
	PvPKills         int64   `json:"pvp_kills"`
	Enchantments     int64   `json:"enchantments"`
	CropsHarvested   int64   `json:"crops_harvested"`
}

// Job represents a player's job
type Job struct {
	Type        JobType   `json:"type"`
	PlayerID    string    `json:"player_id"`
	Level       int       `json:"level"`       // 1-100
	XP          int       `json:"xp"`
	XPToNext    int       `json:"xp_to_next"`
	
	// Income
	HourlyPay   float64   `json:"hourly_pay"`
	TotalEarned float64   `json:"total_earned"`
	
	// Stats
	Stats       JobStats  `json:"stats"`
	
	// Milestones
	Milestones  []JobMilestone `json:"milestones,omitempty"`
	
	// Timing
	JoinedAt    time.Time `json:"joined_at"`
	LastWorked  time.Time `json:"last_worked"`
}

// JobMilestone represents a completed milestone
type JobMilestone struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AchievedAt  time.Time `json:"achieved_at"`
	Reward      float64   `json:"reward"`
}

// NewJob creates a new job
func NewJob(jobType JobType, playerID string) *Job {
	return &Job{
		Type:        jobType,
		PlayerID:    playerID,
		Level:       1,
		XP:          0,
		XPToNext:    calculateXPForLevel(2),
		HourlyPay:   calculateHourlyPay(1),
		Stats:       JobStats{},
		Milestones:  make([]JobMilestone, 0),
		JoinedAt:    time.Now(),
		LastWorked:  time.Now(),
	}
}

// calculateXPForLevel calculates XP needed for a level
func calculateXPForLevel(level int) int {
	// XP curve: level^2 * 100
	return level * level * 100
}

// calculateHourlyPay calculates pay based on level
func calculateHourlyPay(level int) float64 {
	// Base: $10/hour at level 1
	// Growth: +$0.50 per level
	return 10.0 + (float64(level-1) * 0.5)
}

// AddXP adds XP and checks for level up
func (j *Job) AddXP(amount int) (leveledUp bool, newLevel int) {
	j.XP += amount
	j.LastWorked = time.Now()
	
	// Check for level up
	for j.XP >= j.XPToNext && j.Level < 100 {
		j.XP -= j.XPToNext
		j.Level++
		j.XPToNext = calculateXPForLevel(j.Level + 1)
		j.HourlyPay = calculateHourlyPay(j.Level)
		leveledUp = true
		newLevel = j.Level
		
		// Check for milestones
		j.checkMilestones()
	}
	
	return leveledUp, j.Level
}

// checkMilestones checks and awards milestones
func (j *Job) checkMilestones() {
	milestones := j.getMilestonesForLevel(j.Level)
	for _, m := range milestones {
		// Check if already achieved
		alreadyHas := false
		for _, existing := range j.Milestones {
			if existing.Name == m.Name {
				alreadyHas = true
				break
			}
		}
		
		if !alreadyHas {
			m.AchievedAt = time.Now()
			j.Milestones = append(j.Milestones, m)
			j.TotalEarned += m.Reward
		}
	}
}

// getMilestonesForLevel returns milestones for a given level
func (j *Job) getMilestonesForLevel(level int) []JobMilestone {
	milestones := make([]JobMilestone, 0)
	
	// Level 10: Novice
	if level == 10 {
		milestones = append(milestones, JobMilestone{
			Name:        "Novice " + j.Type.String(),
			Description: "Reached level 10",
			Reward:      50.0,
		})
	}
	
	// Level 25: Apprentice
	if level == 25 {
		milestones = append(milestones, JobMilestone{
			Name:        "Apprentice " + j.Type.String(),
			Description: "Reached level 25",
			Reward:      100.0,
		})
	}
	
	// Level 50: Journeyman
	if level == 50 {
		milestones = append(milestones, JobMilestone{
			Name:        "Journeyman " + j.Type.String(),
			Description: "Reached level 50",
			Reward:      250.0,
		})
	}
	
	// Level 75: Expert
	if level == 75 {
		milestones = append(milestones, JobMilestone{
			Name:        "Expert " + j.Type.String(),
			Description: "Reached level 75",
			Reward:      500.0,
		})
	}
	
	// Level 100: Master
	if level == 100 {
		milestones = append(milestones, JobMilestone{
			Name:        "Master " + j.Type.String(),
			Description: "Reached level 100 (Max)",
			Reward:      1000.0,
		})
	}
	
	return milestones
}

// GetTitle returns job title based on level
func (j *Job) GetTitle() string {
	switch {
	case j.Level >= 100:
		return "Master " + j.Type.String()
	case j.Level >= 75:
		return "Expert " + j.Type.String()
	case j.Level >= 50:
		return "Journeyman " + j.Type.String()
	case j.Level >= 25:
		return "Apprentice " + j.Type.String()
	case j.Level >= 10:
		return "Novice " + j.Type.String()
	default:
		return j.Type.String() + " Trainee"
	}
}

// RecordBlockMined records a block mined (for Miner job)
func (j *Job) RecordBlockMined(blockType string, xp int) {
	if j.Type == JobMiner {
		j.Stats.BlocksMined++
		j.AddXP(xp)
	}
}

// RecordBlockPlaced records a block placed (for Builder job)
func (j *Job) RecordBlockPlaced(xp int) {
	if j.Type == JobBuilder {
		j.Stats.BlocksPlaced++
		j.AddXP(xp)
	}
}

// RecordMobKilled records a mob kill (for Hunter job)
func (j *Job) RecordMobKilled(mobType string, xp int) {
	if j.Type == JobHunter {
		j.Stats.MobsKilled++
		j.AddXP(xp)
	}
}

// RecordDistance records distance traveled (for Explorer job)
func (j *Job) RecordDistance(distance float64) {
	if j.Type == JobExplorer {
		j.Stats.DistanceTraveled += distance
		// XP every 100 blocks
		xp := int(distance / 100)
		if xp > 0 {
			j.AddXP(xp)
		}
	}
}

// RecordItemCrafted records an item crafted (for Craftsman job)
func (j *Job) RecordItemCrafted(xp int) {
	if j.Type == JobCraftsman {
		j.Stats.ItemsCrafted++
		j.AddXP(xp)
	}
}

// RecordTrade records a trade (for Merchant job)
func (j *Job) RecordTrade(xp int) {
	if j.Type == JobMerchant {
		j.Stats.TradesCompleted++
		j.AddXP(xp)
	}
}

// RecordPvPKill records a PvP kill (for Warrior job)
func (j *Job) RecordPvPKill(xp int) {
	if j.Type == JobWarrior {
		j.Stats.PvPKills++
		j.AddXP(xp)
	}
}

// JobManager manages player jobs
type JobManager struct {
	jobs       map[string]*Job // key: "playerID_jobType"
	byPlayer   map[string][]JobType
	
	walletMgr  *WalletManager
	
	// Global settings
	XPModifiers map[JobType]float64
}

// NewJobManager creates a new job manager
func NewJobManager(walletMgr *WalletManager) *JobManager {
	return &JobManager{
		jobs:        make(map[string]*Job),
		byPlayer:    make(map[string][]JobType),
		walletMgr:   walletMgr,
		XPModifiers: make(map[JobType]float64),
	}
}

// generateJobKey generates a unique key for a job
func generateJobKey(playerID string, jobType JobType) string {
	return fmt.Sprintf("%s_%d", playerID, jobType)
}

// JoinJob allows a player to join a job
func (jm *JobManager) JoinJob(playerID string, jobType JobType) (*Job, error) {
	key := generateJobKey(playerID, jobType)
	
	if _, exists := jm.jobs[key]; exists {
		return nil, fmt.Errorf("already have this job")
	}
	
	job := NewJob(jobType, playerID)
	jm.jobs[key] = job
	jm.byPlayer[playerID] = append(jm.byPlayer[playerID], jobType)
	
	return job, nil
}

// GetJob gets a player's job
func (jm *JobManager) GetJob(playerID string, jobType JobType) (*Job, bool) {
	key := generateJobKey(playerID, jobType)
	job, exists := jm.jobs[key]
	return job, exists
}

// GetPlayerJobs gets all jobs for a player
func (jm *JobManager) GetPlayerJobs(playerID string) []*Job {
	jobTypes := jm.byPlayer[playerID]
	jobs := make([]*Job, 0, len(jobTypes))
	
	for _, jobType := range jobTypes {
		if job, exists := jm.GetJob(playerID, jobType); exists {
			jobs = append(jobs, job)
		}
	}
	
	return jobs
}

// GetPrimaryJob gets the highest level job for a player
func (jm *JobManager) GetPrimaryJob(playerID string) (*Job, bool) {
	jobs := jm.GetPlayerJobs(playerID)
	if len(jobs) == 0 {
		return nil, false
	}
	
	highest := jobs[0]
	for _, job := range jobs {
		if job.Level > highest.Level {
			highest = job
		}
	}
	
	return highest, true
}

// ProcessHourlyPay processes hourly pay for all jobs
func (jm *JobManager) ProcessHourlyPay() {
	for _, job := range jm.jobs {
		// Check if worked in last hour
		if time.Since(job.LastWorked) < 2*time.Hour {
			// Pay for active work
			pay := job.HourlyPay
			
			// Apply modifier
			if modifier, exists := jm.XPModifiers[job.Type]; exists {
				pay *= modifier
			}
			
			wallet := jm.walletMgr.GetOrCreateWallet(job.PlayerID)
			wallet.Add(pay, TransactionJob, "EMPLOYER", fmt.Sprintf("Hourly pay for %s job", job.Type.String()))
			job.TotalEarned += pay
		}
	}
}

// RecordActivity records activity for a player's job
func (jm *JobManager) RecordActivity(playerID string, jobType JobType, activity string, amount int) {
	job, exists := jm.GetJob(playerID, jobType)
	if !exists {
		return
	}
	
	switch activity {
	case "block_mined":
		job.RecordBlockMined("", amount)
	case "block_placed":
		job.RecordBlockPlaced(amount)
	case "mob_killed":
		job.RecordMobKilled("", amount)
	case "distance":
		job.RecordDistance(float64(amount))
	case "item_crafted":
		job.RecordItemCrafted(amount)
	case "trade":
		job.RecordTrade(amount)
	case "pvp_kill":
		job.RecordPvPKill(amount)
	}
}

// GetTopWorkers returns top workers by level
func (jm *JobManager) GetTopWorkers(jobType JobType, count int) []*Job {
	workers := make([]*Job, 0)
	
	for _, job := range jm.jobs {
		if job.Type == jobType {
			workers = append(workers, job)
		}
	}
	
	// Sort by level (bubble sort)
	for i := 0; i < len(workers); i++ {
		for j := i + 1; j < len(workers); j++ {
			if workers[i].Level < workers[j].Level {
				workers[i], workers[j] = workers[j], workers[i]
			}
		}
	}
	
	if count > len(workers) {
		count = len(workers)
	}
	
	return workers[:count]
}

// GetAvailableJobs returns all available job types
func (jm *JobManager) GetAvailableJobs() []JobType {
	return []JobType{
		JobMiner,
		JobBuilder,
		JobHunter,
		JobFarmer,
		JobExplorer,
		JobFisher,
		JobCraftsman,
		JobMerchant,
		JobWarrior,
		JobEnchanter,
	}
}

// SetXPModifier sets an XP modifier for a job type
func (jm *JobManager) SetXPModifier(jobType JobType, modifier float64) {
	jm.XPModifiers[jobType] = modifier
}
