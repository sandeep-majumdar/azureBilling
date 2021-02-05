package azureBilling

// Global variables

// ###############################################################
// Lookups
// ###############################################################

// MeterLookup.items=map[meterId]priceMeterItem
var MeterLookup priceMeter

// VmSizeLookup.items=map[lowercase VM]vmLookupItem
var VmSizeLookup vmLookup

// ManagedDiskLookup.items=map[lowercase MeterName]managedDiskLookupItem
var ManagedDiskLookup managedDiskLookup

// PlatformMapLookup.items=map[subscriptionid + "/" + rgName]PlatformMapLookupItem
var PlatformMapLookup platformMapLookup

// ReportingCategoryLookup.items=map[metercategory]ReportingCategoryLookupItem
var ReportingCategoryLookup reportingCategoryLookup

// ###############################################################
// Aggregates popuplated during the read of the billing CSV
// ###############################################################

// stringKey: reportingCategory + "/" + reportingSubCategory
var AggregateTotal aggregateTotal

// stringKey: reportingCategory + "/" + reportingSubCategory + "/" + SubscriptionId + "/" + ResourceGroup
var AggregateResourceGroup aggregateResourceGroup

// stringKey: reportingCategory + "/" + reportingSubCategory + "/" + SubscriptionId + "/" + ResourceGroup
var AggregatePlatform aggregatePlatform
