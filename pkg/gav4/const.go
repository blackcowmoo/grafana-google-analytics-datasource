package gav4

import (
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
)

const (
	GaDefaultIdx           = 1
	GaAdminMaxResult       = 200
	GaReportMaxResult      = 100000
	GaRealTimeMinMinute    = 0 * time.Minute
	GaRealTimeMaxMinute    = 29 * time.Minute
	Ga360RealTimeMaxMinute = 59 * time.Minute
)

// Realtime metrics and dimensions not provided by Google Analytics
// 리얼타임 메트릭 디멘션은 구글 api에서 제공하지 않기때문에 상수화로 처리
var GaRealTimeDimensions = []model.MetadataItem{
	{
		ID: "appVersion",
		Attributes: model.MetadataItemAttribute{
			UIName:      "App version",
			Description: "The app's versionName (Android) or short bundle version (iOS).",
		},
	},
	{
		ID: "audienceId",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Audience ID",
			Description: "The numeric identifier of an Audience. Users are reported in the audiences to which they belonged during the report's date range. Current user behavior does not affect historical audience membership in reports.",
		},
	},
	{
		ID: "audienceName",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Audience name",
			Description: "The given name of an Audience. Users are reported in the audiences to which they belonged during the report's date range. Current user behavior does not affect historical audience membership in reports.",
		},
	},
	{
		ID: "city",
		Attributes: model.MetadataItemAttribute{
			UIName:      "City",
			Description: "The city from which the user activity originated.",
		},
	},
	{
		ID: "cityId",
		Attributes: model.MetadataItemAttribute{
			UIName:      "City ID",
			Description: "The geographic ID of the city from which the user activity originated, derived from their IP address.",
		},
	},
	{
		ID: "country",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Country",
			Description: "The country from which the user activity originated.",
		},
	},
	{
		ID: "countryId",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Country ID",
			Description: "The geographic ID of the country from which the user activity originated, derived from their IP address. Formatted according to ISO 3166-1 alpha-2 standard.",
		},
	},
	{
		ID: "deviceCategory",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Device category",
			Description: "The type of device: Desktop, Tablet, or Mobile.",
		},
	},
	{
		ID: "eventName",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Event name",
			Description: "The name of the event.",
		},
	},
	{
		ID: "minutesAgo",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Realtime minutes ago",
			Description: "The number of minutes ago that an event was collected. 00 is the current minute, and 01 means the previous minute.",
		},
	},
	{
		ID: "platform",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Platform",
			Description: "The platform on which your app or website ran; for example, web, iOS, or Android. To determine a stream's type in a report, use both platform and streamId.",
		},
	},
	{
		ID: "streamId",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Stream ID",
			Description: "The numeric data stream identifier for your app or website.",
		},
	},
	{
		ID: "streamName",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Stream name",
			Description: "The data stream name for your app or website.",
		},
	},
	{
		ID: "unifiedScreenName",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Page title and screen name",
			Description: "The page title (web) or screen name (app) on which the event was logged.",
		},
	},
}

var GaRealTimeMetrics = []model.MetadataItem{
	{
		ID: "activeUsers",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Active users",
			Description: "The number of distinct users who visited your site or app.",
		},
	},
	{
		ID: "conversions",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Conversions",
			Description: "The count of conversion events. Events are marked as conversions at collection time; changes to an event's conversion marking apply going forward. You can mark any event as a conversion in Google Analytics, and some events (i.e. first_open, purchase) are marked as conversions by default. To learn more, see About conversions.",
		},
	},
	{
		ID: "eventCount",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Event count",
			Description: "The count of events.",
		},
	},
	{
		ID: "screenPageViews",
		Attributes: model.MetadataItemAttribute{
			UIName:      "Views",
			Description: "The number of app screens or web pages your users viewed. Repeated views of a single page or screen are counted. (screen_view + page_view events).",
		},
	},
}
