package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// MemoryStore is a placeholder interface for operator/context memory.
// A future mem0-backed implementation can satisfy this contract.
type MemoryStore interface {
	GetOperatorContext(ctx context.Context, key string) (*OperatorContextMemory, error)
	UpsertOperatorContext(ctx context.Context, mem OperatorContextMemory) error
	RecordForecastSummary(ctx context.Context, key string, summary ForecastMemorySummary) error
}

type OperatorContextMemory struct {
	Key                  string   `json:"key"`
	Location             string   `json:"location,omitempty"`
	Region               string   `json:"region,omitempty"`
	FacilityType         string   `json:"facilityType,omitempty"`
	PreferredAlertStyle  string   `json:"preferredAlertStyle,omitempty"`
	KnownVulnerabilities []string `json:"knownVulnerabilities,omitempty"`
	LastForecastSummary  string   `json:"lastForecastSummary,omitempty"`
}

type ForecastMemorySummary struct {
	Location string `json:"location"`
	Region   string `json:"region,omitempty"`
	Date     string `json:"date"`
	Summary  string `json:"summary"`
}

// InMemoryStore is a lightweight placeholder implementation for local development.
type InMemoryStore struct {
	mu       sync.RWMutex
	profiles map[string]OperatorContextMemory
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{profiles: make(map[string]OperatorContextMemory)}
}

func (s *InMemoryStore) GetOperatorContext(_ context.Context, key string) (*OperatorContextMemory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mem, ok := s.profiles[key]
	if !ok {
		return nil, nil
	}
	copy := mem
	return &copy, nil
}

func (s *InMemoryStore) UpsertOperatorContext(_ context.Context, mem OperatorContextMemory) error {
	if strings.TrimSpace(mem.Key) == "" {
		return fmt.Errorf("memory key is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.profiles[mem.Key]
	if ok {
		if mem.Location == "" {
			mem.Location = existing.Location
		}
		if mem.Region == "" {
			mem.Region = existing.Region
		}
		if mem.FacilityType == "" {
			mem.FacilityType = existing.FacilityType
		}
		if mem.PreferredAlertStyle == "" {
			mem.PreferredAlertStyle = existing.PreferredAlertStyle
		}
		if len(mem.KnownVulnerabilities) == 0 {
			mem.KnownVulnerabilities = existing.KnownVulnerabilities
		}
		if mem.LastForecastSummary == "" {
			mem.LastForecastSummary = existing.LastForecastSummary
		}
	}

	s.profiles[mem.Key] = mem
	return nil
}

func (s *InMemoryStore) RecordForecastSummary(ctx context.Context, key string, summary ForecastMemorySummary) error {
	mem, err := s.GetOperatorContext(ctx, key)
	if err != nil {
		return err
	}

	updated := OperatorContextMemory{Key: key}
	if mem != nil {
		updated = *mem
	}
	updated.Key = key
	updated.Location = summary.Location
	if summary.Region != "" {
		updated.Region = summary.Region
	}
	updated.LastForecastSummary = summary.Summary

	return s.UpsertOperatorContext(ctx, updated)
}

func memoryKey(location, region string) string {
	base := strings.ToLower(strings.TrimSpace(location))
	if region != "" {
		base += "::" + strings.ToLower(strings.TrimSpace(region))
	}
	return base
}
