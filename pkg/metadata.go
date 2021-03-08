package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	AttributeTypeDimension AttributeType = "DIMENSION"
	AttributeTypeMetric    AttributeType = "METRIC"
)

const metadataURL string = "https://www.googleapis.com/analytics/v3/metadata/ga/columns?pp=1"

type Metadata struct {
	Kind           string         `json:"kind"`
	Etag           string         `json:"etag"`
	TotalResults   int64          `json:"totalResults"`
	AttributeNames []string       `json:"attributeNames"`
	Items          []MetadataItem `json:"items"`
}

type MetadataItem struct {
	ID         string                `json:"id"`
	Kind       string                `json:"kind"`
	Attributes MetadataItemAttribute `json:"attributes"`
}

type MetadataItemAttribute struct {
	Type              AttributeType `json:"type,omitempty"`
	DataType          string        `json:"dataType,omitempty"`
	Group             string        `json:"group,omitempty"`
	Status            string        `json:"status,omitempty"`
	UIName            string        `json:"uiName,omitempty"`
	Description       string        `json:"description,omitempty"`
	AllowedInSegments string        `json:"allowedInSegments,omitempty"`
	AddedInAPIVersion string        `json:"addedInApiVersion,omitempty"`
}

type AttributeType string

func (ga *GoogleAnalytics) getMetadata() (*Metadata, error) {
	res, err := http.Get(metadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata api %w", err)
	}
	defer res.Body.Close()

	metadata := Metadata{}

	err = json.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("fail to parsing metadata to json %w", err)
	}
	return &metadata, nil
}

func (ga *GoogleAnalytics) getFilteredMetadata() ([]MetadataItem, []MetadataItem, error) {
	dimensionCacheKey := "ga:metadata:dimension"
	metricCacheKey := "ga:metadata:metric"
	if dimension, _, found := ga.Cache.GetWithExpiration(dimensionCacheKey); found {
		if metric, _, found := ga.Cache.GetWithExpiration(metricCacheKey); found {
			return dimension.([]MetadataItem), metric.([]MetadataItem), nil
		}
	}
	metadata, err := ga.getMetadata()
	if err != nil {
		return nil, nil, err
	}

	// length := int(metadata.TotalResults)
	var dimensionItems = make([]MetadataItem, 0)
	var metricItems = make([]MetadataItem, 0)
	for _, item := range metadata.Items {
		if item.Attributes.Status == "DEPRECATED" {
			continue
		}
		if item.Attributes.Type == AttributeTypeDimension {
			dimensionItems = append(dimensionItems, item)
		} else if item.Attributes.Type == AttributeTypeMetric {
			metricItems = append(metricItems, item)
		}
	}
	ga.Cache.Set(dimensionCacheKey, dimensionItems, time.Hour)
	ga.Cache.Set(metricCacheKey, metricItems, time.Hour)

	return dimensionItems, metadata.Items, nil
}

func (ga *GoogleAnalytics) GetDimensions() ([]MetadataItem, error) {
	dimension, _, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	return dimension, nil
}

func (ga *GoogleAnalytics) GetMetrics() ([]MetadataItem, error) {
	_, metric, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	return metric, nil
}
