// Package entities provides compatibility aliases for the entity subpackages.
//
// New code should import the focused package it needs:
//   - entities/action for action request payloads
//   - entities/configuration for simulation settings
//   - entities/event for scheduled events
//   - entities/topology for topology models
package entities

import (
	"chain_simulation/entities/action"
	"chain_simulation/entities/configuration"
	"chain_simulation/entities/event"
	"chain_simulation/entities/topology"
)

// Action request payloads.
type (
	Attack                            = action.Attack
	ScheduledMaliciousParams          = action.ScheduledMaliciousParams
	InsertSessionTableEntriesInstance = action.InsertSessionTableEntriesInstance
	ModifyBloomFilter                 = action.ModifyBloomFilter
	StartClient                       = action.StartClient
	InitOsmd                          = action.InitOsmd
	StartServer                       = action.StartServer
	SyncInstance                      = action.SyncInstance
)

var (
	NewAttackInstance           = action.NewAttackInstance
	NewScheduledMaliciousParams = action.NewScheduledMaliciousParams
	NewModifyBloomFilter        = action.NewModifyBloomFilter
	NewStartClient              = action.NewStartClient
	NewInitOsmd                 = action.NewInitOsmd
	NewStartServer              = action.NewStartServer
)

// Experiment scheduling models.
type (
	ConfigurationSetting = configuration.Setting
	Event                = event.Event
)

// Topology models.
type (
	RatioDistribution   = topology.RatioDistribution
	MaliciousParams     = topology.MaliciousParams
	SpecialParams       = topology.SpecialParams
	Node                = topology.Node
	Link                = topology.Link
	BlockChainParams    = topology.BlockChainParams
	SecPathMabParams    = topology.SecPathMabParams
	TopologyParams      = topology.Params
	TopologyStartParams = topology.StartParams
	DynamicParameters   = topology.DynamicParameters
)

var NewTopologyStartParams = topology.NewStartParams
