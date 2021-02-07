package azureBilling

// Global variables

// ###############################################################
// Lookups
// ###############################################################

var ConfigMap Config

var MeterLookup priceMeter
var VmSizeLookup vmLookup
var ManagedDiskLookup managedDiskLookup
var PlatformMapLookup platformMapLookup
var ReportingCategoryLookup reportingCategoryLookup
var SummaryCategoryLookup summaryCategoryLookup

// ###############################################################
// Aggregates popuplated during the read of the billing CSV
// ###############################################################
var AggregateResourceGroup aggregateResourceGroup
